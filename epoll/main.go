package main

import (
	"fmt"
	"net"
	"os"
	"strconv"
	"syscall"
	"time"
)

/*
epoll is a data structure in kernel. It has following syscalls that
can be used to create/modify/wait on it.

int epoll_create(int size);
int epoll_ctl(int epfd, int op, int fd, struct epoll_event *event);
int epoll_wait(int epfd, struct epoll_event *evlist, int maxevents, int timeout);
*/

/*
GO provides a wrapper for each of these epoll syscalls.
*/

func sayHello(cfd int) {
	fmt.Println("accepted the connection.")
	syscall.Write(cfd, []byte("Hola!\nBye!\n"))
	// Keep the socket open for 1 seconds.
	time.Sleep(1 * time.Second)
	fmt.Println("closing the connection socket.")
	defer syscall.Close(cfd)
}

const MaxConnections = 50

func main() {

	// Create a socket in a non blocking mode.
	fd, err := syscall.Socket(syscall.AF_INET, syscall.O_NONBLOCK|syscall.SOCK_STREAM, 0)
	if err != nil {
		fmt.Printf("socket error: %s\n", err)
		os.Exit(1)
	}
	defer syscall.Close(fd)

	fmt.Println("Created a new socket.")
	if err := syscall.SetNonblock(fd, true); err != nil {
		fmt.Printf("setnonblock error: %s\n", err)
		os.Exit(1)
	}

	port, _ := strconv.Atoi(os.Args[1])
	addr := syscall.SockaddrInet4{Port: port}
	copy(addr.Addr[:], net.ParseIP("0.0.0.0.0").To4())

	err = syscall.Bind(fd, &addr)
	if err != nil {
		fmt.Printf("failed to bind: %s\n", err)
		os.Exit(1)

	}
	fmt.Println("Bind the socket to address.")

	err = syscall.Listen(fd, 10)
	if err != nil {
		fmt.Printf("cannot listen on the socket. error: %s\n", err)
		os.Exit(1)
	}
	fmt.Println("Start listening for connections.")

	// Create an even poll data structure.
	// We use the epoll_create1 variant for simplicity.
	// int epoll_create(int size);
	epfd, err := syscall.EpollCreate1(0)
	if err != nil {
		fmt.Println("epoll_create1: ", err)
		os.Exit(1)
	}
	defer syscall.Close(epfd)
	fmt.Println("Create the epoll datastructure.")

	// Now we prepare to call epoll_ctl to tell kernel
	// about the events that we are interested in for
	// the given fds.
	var event syscall.EpollEvent
	event.Events = syscall.EPOLLIN
	event.Fd = int32(fd)

	// int epoll_ctl(int epfd, int op, int fd, struct epoll_event *event);
	if err := syscall.EpollCtl(epfd, syscall.EPOLL_CTL_ADD, fd, &event); err != nil {
		fmt.Println("epollctl error: ", err)
		os.Exit(1)
	}
	fmt.Println("Told kernel about the connection fd we are interested in.")

	// Create an array of events for kernel to respond us back
	var events [MaxConnections]syscall.EpollEvent
	for {		
		fmt.Println("Waiting for events.")
		//	int epoll_wait(int epfd, struct epoll_event *evlist, int maxevents, int timeout);
		nvenets, err := syscall.EpollWait(epfd, events[:], -1)
		if err != nil {
			fmt.Println("epoll wait: ", err)
			continue
		}
		fmt.Println("Received some events.")

		for i := 0; i < nvenets; i++ {
			// Something is received on the listening socket.
			if int(events[i].Fd) == fd {
				// Accept the connection, send back the response
				// and close the connection.
				cfd, _, err := syscall.Accept(fd)
				if err != nil {
					fmt.Println("error in accept.")
					continue
				}
				go sayHello(cfd)
			} else {
				fmt.Println("received something on the socket that we are not interested in?")
			}
		}
	}
}
