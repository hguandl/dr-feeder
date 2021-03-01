# Doctor Feeder

A watcher for news by Arknights.

## Usage

### Download pre-compiled binaries from GitHub Release

<https://github.com/hguandl/dr-feeder/releases/latest>

### Command line arguments

```bash
$ dr-feeder -h
Usage of dr-feeder:
  -V    Print current version
  -c string
        Configuration file (default "config.yaml")
  -d    Debug with fake server
```

### Execution example

```bash
$ dr-feeder -c ./config.yaml
2021/02/18 01:01:04 Waiting for post "明日方舟客户端公告"...
2021/02/18 01:01:04 Waiting for post "微博 明日方舟Arknights"...
```

### Example Configuation File

<https://github.com/hguandl/dr-feeder/blob/master/config.yaml>

## Docker Support

As an alternative, you can also use a Docker container.

### Retrieve Docker image
```bash
$ docker pull hguandl/dr-feeder:latest
```

### Startup container

The configuration file must be named by `config.yaml`. Suppose its **absolute** path is `/full/path/config.yaml`.

```bash
$ docker run -d -v /full/path/:/go/etc/ hguandl/dr-feeder
```

## Development

```bash
$ git clone https://github.com/hguandl/dr-feeder.git
$ cd dr-feeder
```

### About fake server

The fake server returns testing data in `tests` directory. When `dr-feeder` is running on debug mode (`-d`), it will use `debug_url` in the configuration file as the API URL instead of the original (real) URL. It is useful for development and testing.

### Testing steps

* Setup `debug_url` in the configuration file:
```yaml
watchers:
- type: weibo
  debug_url: "http://localhost:8088/weibo"
  uid: 6279793937

### Arknights game announcements
- type: akanno
  debug_url: "http://localhost:8088/akanno"
  channel: IOS
```

* Run the fake server:
```bash
$ go run ./cmd/fake-server
2021/03/01 19:12:24 Listen at :8088

```

* Run the program in debug mode:
```bash
$ go run . -d -c <config_file>.yaml
Running on debug mode...

```
