git_hash := $(shell git rev-parse --short HEAD || echo 'main')
current_time = $(shell date +"%Y-%m-%d:T%H:%M:%S")

# Add linker flags
linker_flags = '-s -X go.nuki.io/nuki/nukictl/cmd.BuildTime=${current_time} -X go.nuki.io/nuki/nukictl/cmd.Version=${git_hash}'
bin_name=nukictl

.PHONY:
build:
	go build -ldflags=${linker_flags} -o ${bin_name} .
