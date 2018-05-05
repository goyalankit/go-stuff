package main

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

func main() {

}
