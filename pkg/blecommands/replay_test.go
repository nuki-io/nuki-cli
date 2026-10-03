package blecommands_test

import (
	"encoding/hex"
	"encoding/json"
	"os"
	"path/filepath"
	"strings"
	"testing"

	"github.com/nuki-io/nuki-cli/pkg/blecommands"
	"github.com/stretchr/testify/require"
)

// recordedResponse mirrors the fixture format written by the hardware tests in hwtest/.
type recordedResponse struct {
	Command string          `json:"command"`
	Code    uint16          `json:"code"`
	Payload string          `json:"payload"`
	Parsed  json.RawMessage `json:"parsed,omitempty"`
	Error   string          `json:"error,omitempty"`
}

// TestReplayRecordedResponses re-parses payloads captured from a real device and
// compares the result with what the parser produced at recording time.
func TestReplayRecordedResponses(t *testing.T) {
	files, err := filepath.Glob(filepath.Join("testdata", "recorded", "*.json"))
	require.NoError(t, err)
	if len(files) == 0 {
		t.Skip("no recorded responses, run make hwtest to capture some")
	}

	for _, file := range files {
		raw, err := os.ReadFile(file)
		require.NoError(t, err)
		var recs []recordedResponse
		require.NoError(t, json.Unmarshal(raw, &recs), file)

		name := strings.TrimSuffix(filepath.Base(file), ".json")
		for i, rec := range recs {
			t.Run(name+"/"+rec.Command, func(t *testing.T) {
				payload, err := hex.DecodeString(rec.Payload)
				require.NoError(t, err)

				res, err := blecommands.ParseResponse(blecommands.CommandCode(rec.Code), payload)
				if rec.Error != "" {
					require.EqualError(t, err, rec.Error, "entry %d", i)
				} else {
					require.NoError(t, err, "entry %d", i)
				}
				if rec.Parsed != nil {
					got, err := json.Marshal(res)
					require.NoError(t, err)
					require.JSONEq(t, string(rec.Parsed), string(got), "entry %d", i)
				}
			})
		}
	}
}
