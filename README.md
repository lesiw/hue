# hue

`hue` colors standard error (and optionally standard output) for the command you
provide.

You can customize `hue`'s output colors with the environment variables `HUEOUT`
and `HUEERR`. For example, `HUEERR=bold` will set standard error to bold,
`HUEOUT=blue` will set standard output to blue, and `HUEERR=redbg` will give
standard error a red background. Multiple attributes can be provided in a comma
separated list; for example, `HUEOUT=bold,green` will make standard output both
bold and green.

Power users can provide `HUEOUT` and `HUEERR` specific ANSI escape codes if they
so desire.

## Demo

![lesiw.io/hue animated demo](./demo.gif)

## Installation

### curl

```sh
curl lesiw.io/hue | sh
```

### go install

```sh
go install lesiw.io/hue@latest
```

## Details

By default, `hue` uses [go-iomux](https://github.com/Netflix/go-iomux), which is
a Go port of [io-mux](https://github.com/joshtriplett/io-mux), which is a Rust
implementation of the technique in
[rederr](https://github.com/poettering/rederr). This has the benefit of
preserving order between stdout and stderr at the small cost of using sockets
for stdout and stderr.

In the rare instance where sockets do not perform as expected, `HUEASYNC=1` will
fall back to an approach without sockets, which will break message ordering.
