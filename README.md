# goby 

As we know microsoft CodePush cloud is slow in China, we can use goby to build our's. We can use [qiniu](http://www.qiniu.com/) or [OSS](https://www.aliyun.com/product/oss) to store the files, because it's simple and quick!  Or you can use local storage, just modify conf/app.ini file, it's simple configure.

[![Build Status](https://travis-ci.org/MessageDream/goby.svg?branch=master)](https://travis-ci.org/MessageDream/goby)
[![Build status](https://ci.appveyor.com/api/projects/status/7f1h1vkrs1f6n9qi/branch/master?svg=true&passingText=windows%20build%20passing&failingText=windows%20build%20failing)](https://ci.appveyor.com/project/MessageDream/goby)
[![Go Report Card](https://goreportcard.com/badge/github.com/MessageDream/goby)](https://goreportcard.com/report/github.com/MessageDream/goby)
## INSTALL FROM BINARY
#### Linux
* [386](https://github.com/MessageDream/goby/releases/download/v0.1.1/goby-v0.1.1-linux-386.zip)
* [amd64](https://github.com/MessageDream/goby/releases/download/v0.1.1/goby-v0.1.1-linux-amd64.zip)

#### Windows
* [386](https://github.com/MessageDream/goby/releases/download/v0.1.1/goby-v0.1.1-windows-386.zip)
* [amd64](https://github.com/MessageDream/goby/releases/download/v0.1.1/goby-v0.1.1-windows-amd64.zip)

#### Darwin
* [386](https://github.com/MessageDream/goby/releases/download/v0.1.1/goby-v0.1.1-darwin-386.zip)
* [amd64](https://github.com/MessageDream/goby/releases/download/v0.1.1/goby-v0.1.1-darwin-amd64.zip)


```shell
$ unzip goby-*.zip
$ ./goby server #open http://127.0.0.1:3000 in browser and configure it.
$ code-push login http://127.0.0.1:3000 
```

## INSTALL FROM SOURCE CODE

### Dependencies:

* [go](https://github.com/golang/go)
* [glide](https://github.com/Masterminds/glide)
* [code-push-cli](https://github.com/Microsoft/code-push/tree/master/cli)  (optional)
* Use [react-native-code-push](https://github.com/Microsoft/react-native-code-push) or [react-native-goby](https://github.com/MessageDream/react-native-goby) in client project.

```shell
$ git clone https://github.com/MessageDream/goby.git
$ cd goby
$ glide install
$ go build goby.go
$ ./goby server #open http://127.0.0.1:3000 in browser and configure it.
$ code-push login http://127.0.0.1:3000 
```

## DIFF UPDATE
If you want to client just download the diff code， according to the following steps:

* Edit `conf/app.ini` > `[package]` > `ENABLE_GOOGLE_DIFF = true`
* Use [react-native-goby](https://github.com/MessageDream/react-native-goby) replace [react-native-code-push](https://github.com/Microsoft/react-native-code-push) in client project.

## SUPPORTED
### Database
* mysql
* postgres
* sqlite3
* tidb

### Cache
* memory
* file
* memcache
* redis

### Storage
* [qiniu](http://www.qiniu.com/)
* [OSS](https://www.aliyun.com/product/oss)
* local

### Log
* console
* file
* database
* connection(tcp、unix or udp)
* email(smtp)

## License
MIT License [read](https://github.com/MessageDream/goby/blob/master/LICENSE)