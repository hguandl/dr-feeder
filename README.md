# Doctor Feeder

A watcher for news by Arknights.

## Docker

### Retrieve Docker image
```bash
$ docker pull hguandl/dr-feeder:latest
```

### Startup container

The configuration file must be named by `config.yaml`. Suppose it is placed in `/some/path/config.yaml`.

```bash
$ docker run -d -v /some/path/:/go/etc/ dr-feeder
```

## Installation

### Retrieve via Go command

```bash
$ go get -u github.com/hguandl/dr-feeder
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
  -c string
    	Configuration file (default "config.yaml")
```

```bash
$ dr-feeder -c ./config.yaml
2021/02/18 01:01:04 Waiting for post "明日方舟客户端公告"...
2021/02/18 01:01:04 Waiting for post "微博 明日方舟Arknights"...
```

## Example Configuation File

```yaml
### Currently the version is fixed to 1.0
version: 1.0

notifiers:
### Custom API server
# Send an HTTP POST form to <api_url>,
# which contains "title", "body" and "url".
##
- type: custom
  api_url: "http://localhost:8080/"

### Telegram bot
# Authencate a telegram bot with <api_key>
# to send messages to <chats>.
##
- type: tgbot
  api_key: "123456:foobar"
  chats:
    - "114514"
    - "-1919810"

### Work Wechat app
# See https://hguandl.com/posts/weibo-watcher-deploy/
##
- type: workwx
  corpid: "wx20190501"
  agentid: 1000002
  corpsecret: "rhodesisland"
  touser: "@all"

### iOS Bark
# Official API server: https://api.day.app
##
- type: bark
  tokens:
    - "qwertyuiop"
    - "asdfghjkl"

### IFTTT Webhook
# See https://maker.ifttt.com/use/demo-page
# (use your API key to replace "demo-page" above)
##
- type: ifttt
  webhooks:
    - event: "call"
      api_key: "amiya"
    - event: "sleep"
      api_key: "doctor"
```
