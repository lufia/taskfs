# taskfs

## INSTALL

```
$ go install github.com/lufia/taskfs@latest
```

## USAGE

```
$ mkdir mtpt
$ taskfs mtpt
$ echo add github $github_token >mtpt/ctl
$ ls mtpt/github.com
$ cat mtpt/githubcom/repo@user#1/message
$ fusermount -u mtpt
```

## DEVELOPMENT

```
$ GOOS=linux go build -o taskfs
$ docker build -t taskfs:latest .
$ docker run -t -i --rm --cap-add SYS_ADMIN --device /dev/fuse taskfs
```
