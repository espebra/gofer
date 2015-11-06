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

The following will print the message *foo* in the channel *#bar*.

```
$ curl -i -X POST "http://localhost:8080/" -d "message=foo" -d "target=#bar"
```

The following will print the message *zoo* as a private query to the user *qux*.

```
$ curl -i -X POST "http://localhost:8080/" -d "message=zoo" -d "target=qux"
```
