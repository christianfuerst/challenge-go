# REST-API backed by sqlite3 db written in Go

Just a small app created for a coding challenge from this link:
https://www.maibornwolff.de/coding-challenge-node-js

A detailed specification for this app can be found there.

## Why Go and not NodeJS?

The coding challenge clearly states that the project should be written in NodeJS, so what the heck are you doing?

I'am in the process of learning Go, so I thought the coding challenge would be a nice opportunity to get more practice in Go and see if this fairly new programming language is capable to fullfill those task, that NodeJS can.

I guess it can - coding the project was really straight forward and painless (no callback hell, you know).

## Getting Started

These instructions will get you a copy of the project up and running on your local machine for development and testing purposes. See deployment for notes on how to deploy the project on a live system.

### Prerequisites

- Get an API key from https://openweathermap.org
- Add your API key to config.json
- In order to compile the project you need
    - Go: https://golang.org/doc/install
    - GCC: https://gcc.gnu.org/install/binaries.html

### Build

Clone this repository, change directory, create default folder and run the application.

```
git clone https://github.com/christianfuerst/challenge-go
cd challenge-go
md logs
md db
go run server.go db.go api.go helper.go
```

Now the application should be running and you can access the api with your browser.

```
Sample http-requests:
http://localhost:3000/weather?city=Berlin&day=2018-12-03
http://localhost:3000/weather?city=Berlin&month=2018-12
```


