# mcat - MPEG TS Concatenation

mcat is currently a simple proof of concept using [gots][gots] to read MPEG transport stream packets, and write them to an output.

This is only an early prototype for playing with the API surface at present and is not intended for any kind of production use.

## Build instructions

```shell
go get github.com/igilham/mcat
cd ${GOPATH}/src/github.com/igilham/mcat
go build
```

[gots]: https://github.com/Comcast/gots
