## servefile
Quick way to serve any directory or file

### Installation
Go is a prerequisite, install it from [here](http://golang.org/doc/install.html) if you do not have it installed.

```shell
$ go get github.com/abiosoft/servefile
```

### Usage
If a file is specified, the file will be served irrespective of the URL request path.

```shell
$ servefile -h
Usage of servefile:
  -f=".": path to file or directory to serve, defaults to current directory
  -h=false: show this help
  -p=8080: port to listen on

$ servefile -f /var/www/html
2015/04/26 01:43:33 servefile serving on 8080
```
