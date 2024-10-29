# linux-overlay

A simple OBS overlay for Linux that displays the state of keyboard and mouse buttons.

## Build

```sh
$ ./build.sh
```

Alternately, you can build the binary yourself:

```sh
$ go build -o dist ./cmd/overlay
```

## Run

```sh
$ sudo ./dist/overlay
```

## Configuration

The configuration file is located at `dist/config.yaml`.