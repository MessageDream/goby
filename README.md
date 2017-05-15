# goby

 As we know microsoft CodePush cloud is slow in China, we can use goby to build our's. We can use [qiniu](http://www.qiniu.com/) or [OSS](https://www.aliyun.com/product/oss) to store the files, because it's simple and quick!  Or you can use local storage, just modify conf/app.ini file, it's simple configure.

## INSTALL FROM SOURCE CODE

```shell
$ git clone https://github.com/MessageDream/goby.git
$ cd goby
$ go build goby.go
$ ./goby server #启动服务 浏览器中打开 http://127.0.0.1:3000
```

## License
MIT License [read](https://github.com/MessageDream/LICENSE)