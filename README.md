# taskfs

## INSTALL

```
$ go get github.com/lufia/taskfs
```

## USAGE

```
$ mkdir mtpt
$ taskfs mtpt
$ echo add github $github_token >mtpt/ctl
$ ls mtpt/github
$ cat mtpt/github/repo@user#1/message
$ fusermount -u mtpt
```

## DEVELOPMENT

```
$ docker build -t taskfs:latest .
$ docker run -t -i --rm --cap-add SYS_ADMIN --device /dev/fuse taskfs
```
