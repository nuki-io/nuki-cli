package bleflows

import (
	"context"
	"fmt"
	"log/slog"
	"runtime"
	"time"

	"github.com/nuki-io/nuki-cli/pkg/blecommands"
	"github.com/nuki-io/nuki-cli/pkg/nukible"
)

type Flow struct {
	ble     *nukible.NukiBle
	handler *blecommands.BleHandler
	device  *nukible.Device
	authCtx *AuthorizeContext
	store   AuthStore
	id      string
}

const scanTimeout = 10 * time.Second

// ConnectAuthenticated enables the BLE adapter and returns an authenticated Flow for a paired device.
func ConnectAuthenticated(id string, store AuthStore) (*Flow, error) {
	ble, err := nukible.NewNukiBle()
	if err != nil {
		return nil, fmt.Errorf("failed to enable bluetooth: %w", err)
	}
	// Linux can only connect to devices found by a scan, macOS connects by ID directly.
	if runtime.GOOS == "linux" {
		if err = ble.ScanForDevice(id, scanTimeout); err != nil {
			return nil, fmt.Errorf("failed to scan for device: %w", err)
		}
	}
	flow, err := NewAuthenticatedFlow(ble, id, store)
	if err != nil {
		return nil, fmt.Errorf("failed to create BLE flow: %w", err)
	}
	return flow, nil
}

// ConnectUnauthenticated enables the BLE adapter, scans for the device and returns a Flow for pairing.
func ConnectUnauthenticated(id string, store AuthStore) (*Flow, error) {
	ble, err := nukible.NewNukiBle()
	if err != nil {
		return nil, fmt.Errorf("failed to enable bluetooth: %w", err)
	}
	if err = ble.ScanForDevice(id, scanTimeout); err != nil {
		return nil, fmt.Errorf("failed to scan for device: %w", err)
	}
	flow, err := NewUnauthenticatedFlow(ble, id, store)
	if err != nil {
		return nil, fmt.Errorf("failed to create BLE flow: %w", err)
	}
	return flow, nil
}

// NewAuthenticatedFlow creates a new Flow instance for a Nuki device that was already paired.
func NewAuthenticatedFlow(ble *nukible.NukiBle, id string, store AuthStore) (*Flow, error) {
	f := &Flow{
		ble:   ble,
		id:    id,
		store: store,
	}
	err := f.loadAuthContext(id)
	if err != nil {
		return nil, err
	}
	err = f.connect(id)
	if err != nil {
		return nil, err
	}
	err = f.device.DiscoverKeyturnerUsdio()
	if err != nil {
		return nil, err
	}
	f.initializeHandlerWithCrypto()

	return f, nil
}

// NewUnauthenticatedFlow creates a new Flow instance for a Nuki device that has not been paired yet.
func NewUnauthenticatedFlow(ble *nukible.NukiBle, id string, store AuthStore) (*Flow, error) {
	f := &Flow{
		ble:   ble,
		id:    id,
		store: store,
	}
	err := f.connect(id)
	if err != nil {
		return nil, err
	}
	err = f.device.DiscoverPairing()
	if err != nil {
		return nil, err
	}
	f.initializeHandler()

	return f, nil
}

func (f *Flow) connect(id string) error {
	addr, ok := f.ble.GetDeviceAddress(id)
	if !ok {
		return fmt.Errorf("requested device with MAC %s was not discovered", id)
	}

	device, err := f.ble.Connect(*addr)
	if err != nil {
		return fmt.Errorf("cannot connect to device %s. %s", id, err.Error())
	}
	f.device = device
	return nil
}

func (f *Flow) loadAuthContext(id string) error {
	ctx, err := f.store.Load(id)
	if err != nil {
		return fmt.Errorf("device is not paired yet. %s", err.Error())
	}
	f.authCtx = ctx
	return nil
}

func (f *Flow) initializeHandler() {
	f.handler = blecommands.NewBleHandler(nil, nil)
}

func (f *Flow) initializeHandlerWithCrypto() {
	crypto := blecommands.NewCrypto(f.authCtx.SharedKey)
	f.handler = blecommands.NewBleHandler(crypto, f.authCtx.AuthId)
}

func (f *Flow) getChallenge(ctx context.Context) ([]byte, error) {
	msg := f.handler.ToEncryptedMessage(&blecommands.RequestData{CommandIdentifier: blecommands.CommandChallenge}, GetNonce24())
	raw, err := f.device.WriteUsdio(ctx, msg)
	if err != nil {
		return nil, fmt.Errorf("failed to get challenge from device: %w", err)
	}
	res, err := f.handler.FromEncryptedDeviceResponse(raw)
	if err != nil {
		return nil, fmt.Errorf("failed to get challenge from device: %w", err)
	}
	c, ok := res.(*blecommands.Challenge)
	if !ok {
		return nil, fmt.Errorf("expected challenge from device, got %s", res.GetCommandCode())
	}
	return c.Nonce, nil
}

func (f *Flow) UpdateAuthCtxFromConfig(cfg *blecommands.Config) {
	f.authCtx.Name = cfg.Name
	f.authCtx.NukiId = cfg.NukiID
	f.store.Store(f.id, f.authCtx)
}

// SetResponseRecorder passes every raw response payload received by this flow to fn.
func (f *Flow) SetResponseRecorder(fn func(cmd blecommands.CommandCode, payload []byte)) {
	f.handler.SetRecorder(fn)
}

func (f *Flow) DisconnectDevice() error {
	if f.device == nil {
		return fmt.Errorf("no device connected")
	}
	f.device.Disconnect()
	f.device = nil
	return nil
}

// performSimpleOp sends a command that requires a challenge+PIN and waits for StatusComplete.
// The caller provides an already-built request (with nonce and pin already set).
func (f *Flow) performSimpleOp(ctx context.Context, req blecommands.Request) error {
	return f.exchange(ctx, req, nil)
}

// responseHandler receives each non-final response. Returning true ends the exchange early,
// for commands where the device does not always follow its answer with StatusComplete.
type responseHandler func(blecommands.Response) (done bool)

// exchange sends an encrypted request and passes each response to handle (if non-nil)
// until the device reports StatusComplete or handle returns true.
func (f *Flow) exchange(ctx context.Context, req blecommands.Request, handle responseHandler) error {
	msg := f.handler.ToEncryptedMessage(req, GetNonce24())
	ch, stop := f.device.WriteUsdioStream(ctx, msg)
	defer stop()

	if err := receiveUntilComplete(ctx, ch, f.handler.FromEncryptedDeviceResponse, handle); err != nil {
		return fmt.Errorf("%s: %w", req.GetCommandCode(), err)
	}
	return nil
}

// collectResponses sends req and returns every *T response received before StatusComplete.
func collectResponses[T any, PT responsePtr[T]](ctx context.Context, f *Flow, req blecommands.Request) ([]T, error) {
	var items []T
	if err := f.exchange(ctx, req, collectInto[T, PT](&items)); err != nil {
		return nil, err
	}
	return items, nil
}

type responsePtr[T any] interface {
	*T
	blecommands.Response
}

func collectInto[T any, PT responsePtr[T]](items *[]T) responseHandler {
	return func(res blecommands.Response) bool {
		if r, ok := res.(PT); ok {
			*items = append(*items, *r)
		}
		return false
	}
}

func receiveUntilComplete(
	ctx context.Context,
	ch <-chan []byte,
	decode func([]byte) (blecommands.Response, error),
	handle responseHandler,
) error {
	for {
		select {
		case buf := <-ch:
			res, err := decode(buf)
			if err != nil {
				return err
			}
			slog.Debug("Received response", "cmd", res.GetCommandCode(), "payload", res)
			if s, ok := res.(*blecommands.Status); ok && s.Status == blecommands.StatusComplete {
				return nil
			}
			if handle != nil && handle(res) {
				return nil
			}
		case <-ctx.Done():
			return ctx.Err()
		}
	}
}
