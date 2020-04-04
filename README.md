# battleship

Battleship is simple application, that emulates this game.

## Toolset

* Golang 1.14

## Environment setup

Install Go: 
```bash
cd /tmp
wget https://dl.google.com/go/go1.14.linux-amd64.tar.gz
sudo rm -rf /usr/local/go
sudo tar -C /usr/local -xzf go1.14.linux-amd64.tar.gz
``` 
Setup environment variables, insert into your .profile file or just execute: 
```bash
export GOROOT=/usr/local/go
export GOPATH=$HOME/go
export PATH=$GOPATH/bin:$GOROOT/bin:$PATH
```

## Makefile

This project is integrated with make tool, so to check all the possibilities:
```bash
make help
```

To build project:
```bash
make build
``` 
To run project:
```bash
make run
``` 
To check project with static analyzers:
```bash
make static_check
``` 
To check project with static checks and unit tests:
```bash
make check
```

## Run tests 

Run tests:
```bash
make test
```

## API

perfectly working swagger API can be found at 
http://localhost:8080/swagger/index.html after the server is started.

swagger sources are also available in **docs** dir in the root of project 

