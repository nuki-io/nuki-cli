package cmd

import (
	"os"
	"path/filepath"
	"testing"

	"github.com/spf13/viper"
	"github.com/stretchr/testify/require"
)

func TestConfigFlagFileType(t *testing.T) {
	tests := []struct {
		name    string
		file    string
		content string
		want    string
	}{
		{name: "no extension falls back to yaml", file: "nukiconf", content: "activecontext: old\n", want: "activecontext: new\n"},
		{name: "extension decides the type", file: "nukiconf.json", content: `{"activecontext": "old"}`, want: "{\n  \"activecontext\": \"new\"\n}"},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			t.Cleanup(viper.Reset)
			path := filepath.Join(t.TempDir(), tt.file)
			require.NoError(t, os.WriteFile(path, []byte(tt.content), 0o600))
			cfgFile = path
			t.Cleanup(func() { cfgFile = "" })

			initConfig()
			require.Equal(t, "old", viper.GetString("activecontext"))
			viper.Set("activecontext", "new")
			writeConfig()

			got, err := os.ReadFile(path)
			require.NoError(t, err)
			require.Equal(t, tt.want, string(got))
		})
	}
}

func TestWriteConfigRestrictsPermissions(t *testing.T) {
	tests := []struct {
		name   string
		create bool
	}{
		{name: "existing world-readable file", create: true},
		{name: "new file", create: false},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			t.Cleanup(viper.Reset)
			home := t.TempDir()
			t.Setenv("HOME", home)
			path := filepath.Join(home, ".nukictl")
			if tt.create {
				require.NoError(t, os.WriteFile(path, []byte("activecontext: x\n"), 0o644))
			}

			initConfig()
			viper.Set("activecontext", "y")
			writeConfig()

			info, err := os.Stat(path)
			require.NoError(t, err)
			require.Equal(t, configPerm, info.Mode().Perm())
		})
	}
}
