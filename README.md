# goby

As we know microsoft CodePush cloud is slow in China, we can use goby to build our's. We can use [qiniu](http://www.qiniu.com/) or [OSS](https://www.aliyun.com/product/oss) to store the files, because it's simple and quick!  Or you can use local storage, just modify conf/app.ini file, it's simple configure.

## INSTALL FROM SOURCE CODE

### Dependencies:

* [glide](https://github.com/Masterminds/glide)
* [code-push-cli](https://github.com/Microsoft/code-push/tree/master/cli)
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

## License
MIT License [read](https://github.com/MessageDream/goby/blob/master/LICENSE)