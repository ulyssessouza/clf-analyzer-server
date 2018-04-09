[![Build Status](https://travis-ci.org/ulyssessouza/clf-analyzer-server.svg?branch=master)](https://travis-ci.org/ulyssessouza/clf-analyzer-server)

# What?
The idea of this project is to implement a client-server CLF (https://en.wikipedia.org/wiki/Common_Log_Format) analyzer having it's connections based in websocket endpoints

# Requirements
To build the projects you will need a configured GO environment.
To have so, please refer to https://golang.org/doc/install and follow the sections
- Download the Go distribution
- Install the Go tools

After this you should be ready to build the projects

# How to build and run
You have 2 main ways to build the projects

## 1) Get directly from github.com
```sh
$ go get -u github.com/ulyssessouza/clf-analyser-server
$ go get -u github.com/ulyssessouza/clf-analyser-client
```

At this point you should have the code of the 2 projects under `$GOPATH/src/github.com/ulyssessouza`
and the binaries under `$GOPATH/bin`

## CLI arguments
Please take the time to explore the configurable CLI arguments by using the help section executing:
```sh
$ clf-analyzer-server -h
$ clf-analyzer-client -h
```
# How does it work?
The main optimisation is to have a light processing in the server side.
For this, the server uses a Sqlite for persistency, since having all the ingested logs in memory would be a problem VERY FAST :D


The amount of queries executed is also optimized, since we know that the refresh rate is 10 seconds.
Instead of having the client pulling data from the server, is the server that pushes data to the connected clients every 10 seconds through websockets.
With this the server executes just one SQL query every 10 seconds to update the cache then pushes it to all the connected clients.


The logs are stored in a file called ```sqlite_clf_analyzer.db``` placed in the same directory where the binary was executed

# Observations
### Log generator used for some stress tests:
https://github.com/kiritbasu/Fake-Apache-Log-Generator