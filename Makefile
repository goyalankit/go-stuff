.PHONY: epoll all epoll_clean ns ns_clean


all: epoll ns

clean: epoll_clean ns_clean

epoll:
	cd epoll && mkdir -p bin && go build -o bin/epoll

epoll_clean:
	cd epoll && rm -rf bin

ns:
	cd ns && mkdir -p bin && go build -o bin/ns

ns_clean:
	cd ns && rm -rf bin
