# hue

`hue` colors standard error (and optionally standard output) for the command you
provide. For example, `hue ls /etc/passwd /badfile` will print `/etc/passwd`
normally, but the "No such file or directory" error will be red.

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
curl -L lesiw.io/hue | sh
```

### go install

```sh
go install lesiw.io/hue@latest
```
