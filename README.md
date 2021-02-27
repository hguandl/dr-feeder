# Doctor Feeder

A watcher for news by Arknights.

## Installation

### Download from GitHub Release
<https://github.com/hguandl/dr-feeder/releases/latest>

### Retrieve via Go command

```bash
$ export GOPROXY=https://goproxy.io,direct
$ go get -u github.com/hguandl/dr-feeder/v2
```

### Build from source

```bash
$ git clone https://github.com/hguandl/dr-feeder.git
$ cd dr-feeder
$ go build
$ go install
```

## Usage

```bash
$ dr-feeder -h
Usage of dr-feeder:
  -V    Print current version
  -c string
        Configuration file (default "config.yaml")
  -d    Debug with fake server
```

```bash
$ dr-feeder -c ./config.yaml
2021/02/18 01:01:04 Waiting for post "明日方舟客户端公告"...
2021/02/18 01:01:04 Waiting for post "微博 明日方舟Arknights"...
```

## Docker Support

### Retrieve Docker image
```bash
$ docker pull hguandl/dr-feeder:latest
```

### Startup container

The configuration file must be named by `config.yaml`. Suppose its **absolute** path is `/full/path/config.yaml`.

```bash
$ docker run -d -v /full/path/:/go/etc/ hguandl/dr-feeder
```

## Example Configuation File

<https://github.com/hguandl/dr-feeder/blob/master/config.yaml>
