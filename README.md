# taskfs

## INSTALL

```
$ go get github.com/lufia/taskfs
```

## USAGE

```
$ taskfs -t $github_token mtpt
$ ls mtpt/github
$ fusermount -u mtpt
```

## DEVELOPMENT

```
$ docker build -t taskfs:latest .
$ docker run -t -i --rm --cap-add SYS_ADMIN --device /dev/fuse taskfs
```
