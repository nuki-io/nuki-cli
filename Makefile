git_hash := $(shell git describe --always --tags)
current_time = $(shell date +"%Y-%m-%dT%H:%M:%S")

# Add linker flags
linker_flags = '-s -X github.com/nuki-io/nuki-cli/cmd.BuildTime=${current_time} -X github.com/nuki-io/nuki-cli/cmd.Version=${git_hash}'
bin_name=nukictl

.PHONY:
build:
	go generate ./...
	go build -ldflags=${linker_flags} -o ${bin_name} .

test:
	go test ./...
