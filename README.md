# Doctor Feeder

A watcher for news by Arknights.

## Installation

### Retrieved via Go command

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

### Example Configuation File

```yaml
version: 1.0
notifiers:
- type: custom
  api_url: "http://localhost:8080/"
- type: tgbot
  api_key: "123456:foobar"
  chats:
    - "114514"
    - "-1919810"
- type: workwx
  corpid: "wx20190501"
  agentid: 1000002
  corpsecret: "rhodesisland"
  touser: "@all"
- type: bark
  tokens:
    - "qwertyuiop"
    - "asdfghjkl"
- type: ifttt
  webhooks:
    - event: "call"
      api_key: "amiya"
    - event: "sleep"
      api_key: "doctor"
```
