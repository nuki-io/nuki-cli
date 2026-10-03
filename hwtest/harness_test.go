//go:build hwtest

// Package hwtest runs the BLE flows against a real, already paired Nuki device.
// See the hwtest target in the Makefile for the environment variables it reads.
package hwtest

import (
	"context"
	"encoding/binary"
	"encoding/hex"
	"encoding/json"
	"fmt"
	"log/slog"
	"os"
	"path/filepath"
	"runtime"
	"slices"
	"strings"
	"sync"
	"testing"
	"time"

	"github.com/nuki-io/nuki-cli/internal/authstore"
	"github.com/nuki-io/nuki-cli/pkg/blecommands"
	"github.com/nuki-io/nuki-cli/pkg/bleflows"
	"github.com/nuki-io/nuki-cli/pkg/nukible"
	"github.com/spf13/viper"
	"github.com/stretchr/testify/require"
)

const (
	levelRead = iota
	levelWrite
	levelActuate
)

var levelNames = []string{"read", "write", "actuate"}

const testTimeout = 2 * time.Minute

var (
	ble       *nukible.NukiBle
	deviceID  string
	store     *memStore
	ownAuthID uint32
	deviceCfg *blecommands.Config
	level     int
	allowed   map[string]bool
	recordDir string
)

func TestMain(m *testing.M) {
	if err := setup(); err != nil {
		fmt.Fprintln(os.Stderr, "hwtest setup failed:", err)
		os.Exit(1)
	}
	os.Exit(m.Run())
}

// TestHardware runs the phases in order and stops at the first failing one, so a
// broken parser is caught by the read phase before its output is written back.
func TestHardware(t *testing.T) {
	phases := []struct {
		name  string
		level int
		run   func(t *testing.T)
	}{
		{"read", levelRead, testRead},
		{"write", levelWrite, testWrite},
		{"actuate", levelActuate, testActuate},
		{"disruptive", levelRead, testDisruptive},
	}
	for _, p := range phases {
		ok := t.Run(p.name, func(t *testing.T) {
			requireLevel(t, p.level)
			p.run(t)
		})
		if !ok {
			t.Logf("%s phase failed, skipping the remaining phases", p.name)
			return
		}
	}
}

func setup() error {
	cfgPath := os.Getenv("NUKI_TEST_CONFIG")
	if cfgPath == "" {
		home, err := os.UserHomeDir()
		if err != nil {
			return err
		}
		cfgPath = filepath.Join(home, ".nukictl")
	}
	v := viper.New()
	v.SetConfigFile(cfgPath)
	v.SetConfigType("yaml")
	if err := v.ReadInConfig(); err != nil {
		return fmt.Errorf("reading %s: %w", cfgPath, err)
	}

	deviceID = os.Getenv("NUKI_TEST_DEVICE")
	if deviceID == "" {
		deviceID = v.GetString("activecontext")
	}
	if deviceID == "" {
		return fmt.Errorf("set NUKI_TEST_DEVICE or an active context with nukictl ble set-context")
	}
	auth, err := authstore.New(v).Load(deviceID)
	if err != nil {
		return err
	}
	ownAuthID = binary.LittleEndian.Uint32(auth.AuthId)
	store = &memStore{auth: map[string]bleflows.AuthorizeContext{deviceID: *auth}}

	lvl := os.Getenv("NUKI_TEST_LEVEL")
	if lvl == "" {
		lvl = "read"
	}
	level = slices.Index(levelNames, lvl)
	if level < 0 {
		return fmt.Errorf("NUKI_TEST_LEVEL must be one of %v, got %q", levelNames, lvl)
	}
	allowed = map[string]bool{}
	for _, a := range strings.Split(os.Getenv("NUKI_TEST_ALLOW"), ",") {
		allowed[strings.TrimSpace(a)] = true
	}
	if os.Getenv("NUKI_TEST_DEBUG") != "" {
		slog.SetLogLoggerLevel(slog.LevelDebug)
	}
	recordDir = os.Getenv("NUKI_TEST_RECORD")
	if recordDir != "" {
		if err := os.MkdirAll(recordDir, 0o755); err != nil {
			return err
		}
	}

	ble, err = nukible.NewNukiBle()
	if err != nil {
		return fmt.Errorf("enabling bluetooth: %w", err)
	}
	if runtime.GOOS == "linux" {
		if err := ble.ScanForDevice(deviceID, 10*time.Second); err != nil {
			return fmt.Errorf("scanning: %w", err)
		}
	}
	return probe()
}

// probe checks that the device is reachable and caches its config for capability checks.
func probe() error {
	flow, err := connect()
	if err != nil {
		return err
	}
	defer flow.DisconnectDevice()
	ctx, cancel := context.WithTimeout(context.Background(), 30*time.Second)
	defer cancel()
	deviceCfg, err = flow.GetConfig(ctx)
	if err != nil {
		return fmt.Errorf("reading config: %w", err)
	}
	fmt.Printf("hwtest: %s (firmware %s, hardware %s), level %s\n",
		deviceCfg.Name, deviceCfg.FirmwareVersion, deviceCfg.HardwareRevision, levelNames[level])
	return nil
}

// connect retries because the lock needs a moment to advertise again after a disconnect.
func connect() (*bleflows.Flow, error) {
	var err error
	for attempt := 1; attempt <= 3; attempt++ {
		var flow *bleflows.Flow
		flow, err = bleflows.NewAuthenticatedFlow(ble, deviceID, store)
		if err == nil {
			return flow, nil
		}
		time.Sleep(2 * time.Second)
		if runtime.GOOS == "linux" {
			_ = ble.ScanForDevice(deviceID, 10*time.Second)
		}
	}
	return nil, fmt.Errorf("connecting to %s: %w", deviceID, err)
}

// newFlow connects to the device for the duration of t. Cleanups registered by the
// test after this call run before the device is disconnected.
func newFlow(t *testing.T) (context.Context, *bleflows.Flow) {
	t.Helper()
	flow, err := connect()
	require.NoError(t, err)
	t.Cleanup(func() { flow.DisconnectDevice() })
	if recordDir != "" {
		rec := &recorder{name: t.Name(), counts: map[blecommands.CommandCode]int{}}
		flow.SetResponseRecorder(rec.add)
		t.Cleanup(func() { rec.write(t) })
	}
	ctx, cancel := context.WithTimeout(context.Background(), testTimeout)
	t.Cleanup(cancel)
	return ctx, flow
}

// cleanupCtx is for restoring device state, which must work even if the test ran out of time.
func cleanupCtx(t *testing.T) context.Context {
	ctx, cancel := context.WithTimeout(context.Background(), time.Minute)
	t.Cleanup(cancel)
	return ctx
}

func requireLevel(t *testing.T, l int) {
	t.Helper()
	if level < l {
		t.Skipf("requires NUKI_TEST_LEVEL=%s", levelNames[l])
	}
}

func requireAllowed(t *testing.T, name string) {
	t.Helper()
	if !allowed[name] {
		t.Skipf("requires NUKI_TEST_ALLOW=%s", name)
	}
}

func requireKeypad(t *testing.T) {
	t.Helper()
	if !deviceCfg.HasKeypad && !deviceCfg.HasKeypad2 {
		t.Skip("no keypad paired")
	}
}

// requireKnown fails if v is outside its enum. Stringer renders unknown values as "Type(n)".
func requireKnown(t *testing.T, v fmt.Stringer, field string) {
	t.Helper()
	s := v.String()
	require.False(t, s == "" || strings.Contains(s, "("), "%s has unknown value %q", field, s)
}

func requireRecent(t *testing.T, ts time.Time, field string) {
	t.Helper()
	d := time.Since(ts)
	require.Less(t, d.Abs(), 10*time.Minute,
		"%s is %s off from now (%s). Wrong parsing, or device clock needs nukictl ble update-time", field, d, ts)
}

func logJSON(t *testing.T, label string, v any) {
	t.Helper()
	b, err := json.MarshalIndent(v, "", "  ")
	require.NoError(t, err)
	t.Logf("%s: %s", label, b)
}

// memStore keeps credentials in memory so tests never modify the config file.
type memStore struct {
	mu   sync.Mutex
	auth map[string]bleflows.AuthorizeContext
}

func (s *memStore) Load(id string) (*bleflows.AuthorizeContext, error) {
	s.mu.Lock()
	defer s.mu.Unlock()
	a, ok := s.auth[id]
	if !ok {
		return nil, fmt.Errorf("no authorization for %s", id)
	}
	return &a, nil
}

func (s *memStore) Store(id string, ctx *bleflows.AuthorizeContext) error {
	s.mu.Lock()
	defer s.mu.Unlock()
	s.auth[id] = *ctx
	return nil
}

// recordedResponse is replayed by TestReplayRecordedResponses in pkg/blecommands.
type recordedResponse struct {
	Command string          `json:"command"`
	Code    uint16          `json:"code"`
	Payload string          `json:"payload"`
	Parsed  json.RawMessage `json:"parsed,omitempty"`
	Error   string          `json:"error,omitempty"`
}

// maxRecordedPerCommand keeps polling loops and log streams from bloating fixtures.
const maxRecordedPerCommand = 3

type recorder struct {
	name    string
	counts  map[blecommands.CommandCode]int
	entries []recordedResponse
}

func (r *recorder) add(cmd blecommands.CommandCode, payload []byte) {
	// Keypad code payloads contain the codes themselves.
	if cmd == blecommands.CommandChallenge || cmd == blecommands.CommandKeypadCode || r.counts[cmd] >= maxRecordedPerCommand {
		return
	}
	r.counts[cmd]++
	rec := recordedResponse{Command: cmd.String(), Code: uint16(cmd), Payload: hex.EncodeToString(payload)}
	res, err := blecommands.ParseResponse(cmd, payload)
	if err != nil {
		rec.Error = err.Error()
	}
	if res != nil {
		rec.Parsed, _ = json.Marshal(res)
	}
	r.entries = append(r.entries, rec)
}

func (r *recorder) write(t *testing.T) {
	if len(r.entries) == 0 || t.Failed() {
		return
	}
	b, err := json.MarshalIndent(r.entries, "", "  ")
	if err != nil {
		t.Errorf("encoding recorded responses: %v", err)
		return
	}
	file := filepath.Join(recordDir, strings.NewReplacer("/", "_", " ", "_").Replace(r.name)+".json")
	if err := os.WriteFile(file, b, 0o644); err != nil {
		t.Errorf("writing recorded responses: %v", err)
	}
}
