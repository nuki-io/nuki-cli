git_hash := $(shell git describe --always --tags)
current_time = $(shell date +"%Y-%m-%dT%H:%M:%S")

# Add linker flags
linker_flags = '-s -X github.com/nuki-io/nuki-cli/cmd.BuildTime=${current_time} -X github.com/nuki-io/nuki-cli/cmd.Version=${git_hash}'
bin_name=nukictl$(BIN_SUFFIX)

.PHONY:
build:
	go generate ./...
	go build -ldflags=${linker_flags} -o ${bin_name} .

test:
	go test ./...

# Runs against the paired device from ~/.nukictl. Environment:
#   NUKI_TEST_LEVEL     read (default), write or actuate
#   NUKI_TEST_ALLOW     comma list of opt-ins: unlatch, reboot, set-pin
#   NUKI_TEST_TEMP_PIN  temporary PIN for set-pin
#   NUKI_TEST_DEVICE    device ID, defaults to the active context
#   NUKI_TEST_CONFIG    config file, defaults to ~/.nukictl
#   NUKI_TEST_RECORD    where passing tests save responses for replay, empty disables
#   NUKI_TEST_DEBUG     non-empty enables debug logs including raw payloads
#   RUN                 test filter, e.g. RUN=Hardware/read/Config
NUKI_TEST_RECORD ?= $(CURDIR)/pkg/blecommands/testdata/recorded
.PHONY: hwtest
hwtest:
	NUKI_TEST_RECORD=$(NUKI_TEST_RECORD) go test -tags hwtest -count=1 -p 1 -v -timeout 20m ./hwtest/ $(if $(RUN),-run '$(RUN)')
