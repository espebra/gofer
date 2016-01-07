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

Potentially dangerous.

The bot will listen for private messages and messages on any channel it has joined. For each message, the bot will iterate over files in a directory specific for that channel and execute them with the sender nick name and the message as arguments.

The channel name is used as part of the path in the file system to the command that will be executed. This makes it possible to have different commands enabled in different channels.

Example:

On the channel *#foo*, the user *bar* says *Something is weird with #1337!*. The bot will then look into the directory ``scripts/#foo/`` (``scripts`` set by the ``ScriptDirectory`` configuration options) and iterate over the files there. It will execute each of the files as:

```
scripts/#foo/filename "bar" "Something is weird with #1337!"
```

The ``ScriptDirectory`` configuration option needs to point to a directory. A given directory structure within this directory is required. Consider the following example configuration:

```
{
    [...]
    "ScriptDirectory": "/usr/local/gofer/scripts/"
    [...]
}
```

The above configuration will require the directory structure which is created with the following commands:

```
$ sudo mkdir /usr/local/gofer/scripts
$ sudo mkdir /usr/local/gofer/scripts/#foo
$ sudo mkdir /usr/local/gofer/scripts/#anotherchannel
```

The commands that should be enabled for the IRC channel ``#foo`` are copied or symlinked to ``/usr/local/gofer/scripts/#foo/``. Likewise for channel ``#anotherchannel``.

The commands must exit with 0 for the output to be printed to IRC. Multi line output is supported.
