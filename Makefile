# For each target, we go inside the target dir
# and do `go get` for each of the dependencies
# and then call `go build`. Nothing fancy.

# This is a quick hack before I figure out the
# right way to do multiple binaries in a golang
# project.
# There are better makefiles available for a real
# project: https://github.com/cloudflare/hellogopher

.PHONY: epoll all epoll_clean ns ns_clean 3pc 3pc_clean gogob gogob_clean

all: epoll ns

clean: epoll_clean ns_clean 3pc_clean

epoll:
	cd epoll && mkdir -p bin && go build -o bin/epoll

epoll_clean:
	cd epoll && rm -rf bin

ns:
	cd ns && mkdir -p bin && \
		go get github.com/opencontainers/runtime-spec/specs-go && \
		go build -o bin/ns

ns_clean:
	cd ns && rm -rf bin


3pc:
	cd 3pc && mkdir -p bin && go build -o bin/3pc

3pc_clean:
	cd 3pc && rm -rf bin

gogob:
	cd gogob && mkdir -p bin && go build -o bin/gogob


gogob_clean:
	cd gogob && rm -rf bin
