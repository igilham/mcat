# mcat - MPEG TS Concatenation

mcat is currently a simple proof of concept using [gots][gots] to read [MPEG transport stream][mpeg_ts] packets, and write them to an output.

This is only an early prototype for playing with the API surface at present and is not intended for any kind of production use.

## Build instructions

You need to have [Go][golang] installed and set up.

```shell
go get github.com/igilham/mcat
cd ${GOPATH}/src/github.com/igilham/mcat
go build
```

## Testing

I have tested a simple use case with [Big Buck Bunny][bbb], using [VLC][vlc] to convert to H264/MP3 TS. I ran the resulting transport stream packets through mcat and compared the output file to the original.

### Stdout redirected to a file

```shell
$ ./mcat ~/Downloads/bbb_sunflower_1080p_60fps_normal.ts > ~/Downloads/test.ts
$ cmp ~/Downloads/bbb_sunflower_1080p_60fps_normal.ts ~/Downloads/test.ts
$ echo $?
0
```

### Output directly to a file

```shell
$ ./mcat -o ~/Downloads/test.ts ~/Downloads/bbb_sunflower_1080p_60fps_normal.ts
$ cmp ~/Downloads/bbb_sunflower_1080p_60fps_normal.ts ~/Downloads/test.ts
$ echo $?
0
```

## To do

I want to look into some of the following things to learn more about MPEG, the [gots][gots] API, and [Go][golang] in general.

- Explore options to correct the [programme clock reference (PCR)][pcr] while concatenating transport stream packets
- Explore options to drop [null packets][null_packet]

[bbb]: http://bbb3d.renderfarming.net/index.html
[golang]: https://golang.org
[gots]: https://github.com/Comcast/gots
[mpeg_ts]:https://en.wikipedia.org/wiki/MPEG_transport_stream
[pcr]: https://en.wikipedia.org/wiki/MPEG_transport_stream#PCR
[null_packet]: https://en.wikipedia.org/wiki/MPEG_transport_stream#Null_packets
[vlc]: https://videolan.org/vlc/
