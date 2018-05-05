epoll
====

This demonstrates basic usage of epoll api. 

epoll is a datastructure in kernel. It can be created/modified/waited on using the following three syscalls:

`epoll_create` can be used to create a new insance of epoll data structure.
```
int epoll_create(int size);
```

`epoll_ctl` can be used to register for events.
```
int epoll_ctl(int epfd, int op, int fd, struct epoll_event *event);
```

`epoll_wait` is a blocking call that will return if the fds receive the event
that we expressed our interest in.
```
int epoll_wait(int epfd, struct epoll_event *evlist, int maxevents, int timeout);
```
