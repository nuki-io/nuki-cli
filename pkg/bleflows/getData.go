package bleflows

import (
	"fmt"
	"log/slog"

	"github.com/nuki-io/nuki-cli/pkg/blecommands"
)

func (f *Flow) GetConfig() (*blecommands.Config, error) {
	nonce, err := f.getChallenge()
	if err != nil {
		return nil, fmt.Errorf("failed to get challenge from device: %w", err)
	}

	cfg := &blecommands.RequestConfig{Nonce: nonce}
	msg := f.handler.ToEncryptedMessage(cfg, GetNonce24())
	raw, err := f.device.WriteUsdio(msg)
	if err != nil {
		return nil, fmt.Errorf("failed to get config from device: %w", err)
	}
	res, err := f.handler.FromEncryptedDeviceResponse(raw)
	if err != nil {
		return nil, fmt.Errorf("failed to get config from device: %w", err)
	}

	return res.(*blecommands.Config), nil
}

func (f *Flow) RequestData(cmd blecommands.CommandCode) (*blecommands.Response, error) {
	cfg := &blecommands.RequestData{CommandIdentifier: cmd}
	msg := f.handler.ToEncryptedMessage(cfg, GetNonce24())
	raw, err := f.device.WriteUsdio(msg)
	if err != nil {
		return nil, fmt.Errorf("failed to request data from device: %w", err)
	}
	res, err := f.handler.FromEncryptedDeviceResponse(raw)
	if err != nil {
		return nil, fmt.Errorf("failed to request data from device: %w", err)
	}

	return &res, nil
}

func (f *Flow) GetStatus() (*blecommands.KeyturnerStates, error) {
	res, err := f.RequestData(blecommands.CommandKeyturnerStates)
	if err != nil {
		return nil, fmt.Errorf("failed to request keyturner states: %w", err)
	}
	state, ok := (*res).(*blecommands.KeyturnerStates)
	if !ok {
		return nil, fmt.Errorf("failed to cast response to KeyturnerStates: %w", err)
	}
	return state, nil
}

func (f *Flow) GetLogs(start int, count int) ([]blecommands.LogEntry, error) {
	nonce, err := f.getChallenge()
	if err != nil {
		return nil, fmt.Errorf("failed to get challenge from device: %w", err)
	}

	cfg := &blecommands.RequestLogEntries{
		StartIndex:  uint32(start),
		Count:       uint16(count),
		Nonce:       nonce,
		SortOrder:   blecommands.LogSortOrderDescending,
		TotalCount:  0x00,
		SecurityPin: blecommands.NewPin(f.authCtx.Pin),
	}
	msg := f.handler.ToEncryptedMessage(cfg, GetNonce24())
	entries := []blecommands.LogEntry{}
	_, err = f.device.WriteUsdioWithCallback(msg,
		func(buf []byte, sem chan error) []byte {
			f.onRequestLogResponse(buf, sem, &entries)
			return buf
		})
	return entries, err
}

func (f *Flow) onRequestLogResponse(buf []byte, sem chan error, entries *[]blecommands.LogEntry) {
	slog.Debug("Received response", "buf", fmt.Sprintf("%x", buf))
	res, err := f.handler.FromEncryptedDeviceResponse(buf)
	if err != nil {
		slog.Error("Failed to decrypt response", "err", err)
		sem <- err
		return
	}
	slog.Debug("Received log entry response", "cmd", res.GetCommandCode(), "payload", res)
	if res.GetCommandCode() == blecommands.CommandLogEntry {
		entry, ok := res.(*blecommands.LogEntry)
		if ok {
			*entries = append(*entries, *entry)
		}
	}
	if s, ok := res.(*blecommands.Status); ok && s.Status == blecommands.StatusComplete {
		sem <- nil
	}
}
