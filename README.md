# gofer

Simple IRC bot written in Go, based on [thoj/go-ircevent](https://github.com/thoj/go-ircevent). It will listen for incoming HTTP requests and will relay the content onto IRC as specified.

## Installation

Install ``golang``:

```
$ sudo yum/apt-get/brew install golang
```

Create the Go workspace and set the GOPATH environment variable:

```
$ mkdir ~/go
$ cd ~/go
$ mkdir src bin pkg
$ export GOPATH="~/go"
```

Download and install ``gofer``. The binary will be created as ``~/go/bin/gofer``.

```
$ go get github.com/espebra/gofer
$ cd src/github.com/espebra/gofer
$ go install
```

Copy ``~/go/src/github.com/espebra/gofer/config.json.example`` to ``~/.gofer.json`` and modify it to fit your needs.

## Usage

The built in help text will show the various command line arguments available:

```
~/go/bin/gofer --help
```

Some arguments commonly used to start ``gofer`` are:

```
~/go/bin/gofer --config ~/gofer.json
```

## HTTP interface

The following will print the PRIVMSG *foo* in the channel *#bar*.

```
$ curl -d "message=foo" http://localhost:8080/channel/bar/privmsg
```

The following will print the ACTION *foo* in the channel *#bar*.

```
$ curl -d "message=foo" http://localhost:8080/channel/bar/action
```

The following will print the PRIVMSG *zoo* as a private message to the user *qux*.

```
$ curl -d "message=zoo" http://localhost:8080/user/qux/privmsg
```

The following will print the ACTION *zoo* as a private message to the user *qux*.

```
$ curl -d "message=zoo" http://localhost:8080/user/qux/action
```

## Command execution

Potentially dangerous. IRC messages on the following syntax are turned into commands:

```
!command arg1 arg2
````

``command`` needs to exist as an executable file in a directory ``CommandDirectory`` which is specified in the configuration file. The arguments ``arg1, arg2, ...`` are sent as arguments to the ``command``. The stdout is printed back on IRC.


