.PHONY: epoll all epoll_clean


all: epoll

clean: epoll_clean

epoll:
	cd epoll && mkdir -p bin && go build -o bin/epoll

epoll_clean:
	cd epoll && rm -rf bin
