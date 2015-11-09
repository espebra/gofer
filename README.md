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

## HTTP API interface

### Say privmsg in specified channel

|                       | Value                 	|
| --------------------- | ----------------------------- |
| **Method**            | POST				|
| **URL**               | /channel/:channel/privmsg	|
| **URL parameters**    | *None*                	|
| **Success response**  | ``200``               	|
| **Error response**    |                       	|
| **Form data**         | *message=Some message text*	|

###### Example

The following will print the PRIVMSG *foo* in the channel *#bar*.

```
$ curl -d "message=foo" http://localhost:8080/channel/bar/privmsg
```

### Say action in specified channel

|                       | Value                 	|
| --------------------- | ----------------------------- |
| **Method**            | POST				|
| **URL**               | /channel/:channel/action	|
| **URL parameters**    | *None*                	|
| **Success response**  | ``200``               	|
| **Error response**    |                       	|
| **Form data**         | *message=Some action text*	|

###### Example

The following will print the ACTION *foo* in the channel *#bar*.

```
$ curl -d "message=foo" http://localhost:8080/channel/bar/action
```

### Say privmsg to specified user

|                       | Value                 	|
| --------------------- | ----------------------------- |
| **Method**            | POST				|
| **URL**               | /user/:nickname/privmsg	|
| **URL parameters**    | *None*                	|
| **Success response**  | ``200``               	|
| **Error response**    |                       	|
| **Form data**         | *message=Some message text*	|

###### Example

The following will print the PRIVMSG *zoo* as a private message to the user *qux*.

```
$ curl -d "message=zoo" http://localhost:8080/user/qux/privmsg
```

### Say action to specified user

|                       | Value                 	|
| --------------------- | ----------------------------- |
| **Method**            | POST				|
| **URL**               | /user/:nickname/action	|
| **URL parameters**    | *None*                	|
| **Success response**  | ``200``               	|
| **Error response**    |                       	|
| **Form data**         | *message=Some action text*	|

###### Example

The following will print the ACTION *zoo* as a private message to the user *qux*.

```
$ curl -d "message=zoo" http://localhost:8080/user/qux/action
```

## Command execution

Potentially dangerous. IRC messages on the following syntax are turned into commands:

```
!command arg1 arg2
````

``command`` needs to exist as an executable file in a given directory. The arguments ``arg1, arg2, ...`` are sent as arguments to the ``command``. Stdout is printed back on IRC.

The channel name is used as part of the path in the file system to the command that will be executed. This makes it possible to have different commands enabled in different channels.

The ``CommandDirectory`` configuration option needs to point ot a directory. A given directory structure within this directory is required. Consider the following example configuration:

```
{
    [...]
    "CommandDirectory": "/usr/local/gofer/scripts/"
    [...]
}
```

The above configuration will require the directory structure which is created with the following commands:

```
$ sudo mkdir /usr/local/gofer/scripts
$ sudo mkdir /usr/local/gofer/scripts/channel
$ sudo mkdir /usr/local/gofer/scripts/channel/bar
$ sudo mkdir /usr/local/gofer/scripts/channel/foo
```

The commands that should be enabled for the IRC channel ``#bar`` are copied or symlinked to ``/usr/local/gofer/scripts/channel/bar/``. Likewise for channel ``#foo``.

One example command, ``time``, is available in the scripts directory in this repository. It will allow IRC users to get the current time in a specified timezone:

```
22:18:43   user | !time europe/moscow
22:18:43  gofer | Sat Nov  7 00:18:43 MSK 2015
22:19:22   user | !time
22:19:22  gofer | Fri Nov  6 22:19:21 CET 2015
22:19:40   user | !time utc
22:19:40  gofer | Fri Nov  6 21:19:40 UTC 2015
```

