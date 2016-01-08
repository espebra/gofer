# gofer

Simple IRC bot written in Go, based on [thoj/go-ircevent](https://github.com/thoj/go-ircevent). It will listen for incoming HTTP requests and relay the content onto IRC as specified. It will also execute commands based on messages in IRC channels and print back the result.

## Table of contents

* [Installation](#installation)
* [HTTP API interface](#http-api-interface)
* [Command execution](#command-execution)

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

Download and install ``gofer``. The statically linked binary will be created as ``~/go/bin/gofer``.

```
$ go get github.com/espebra/gofer
```

Copy ``~/go/src/github.com/espebra/gofer/config.json.example`` to ``~/.gofer.json`` and modify it to fit your needs.

## Usage

The built in help text will show the various command line arguments available:

```
~/go/bin/gofer --help
```

The common way to start gofer is:

```
~/go/bin/gofer --config ~/gofer.json
```

An example [systemd service script](systemd/gofer.service) is provided to make it easy to daemonize the bot and log to syslog.

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

On the channel *#foo*, the user *bar* says *Something is weird with #1337!*. The bot will then look into the directory ``scripts/#foo/`` (``scripts`` set by the ``ScriptDirectory`` configuration option) and iterate over the files there. The files will be executed as follows:

```
scripts/#foo/script1 "bar" "Something is weird with #1337!"
scripts/#foo/script2 "bar" "Something is weird with #1337!"
```

It's up to each script to parse the message and decide if they have anything relevant to say about it. If the exit code is 0 and some output is printed, the output is printed to IRC. Multi line output is supported.

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
