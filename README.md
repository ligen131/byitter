# byitter

A simple web application.

**This project is work in progress.**

## Overview

|  No. | Project                           | Completed |
| ---: | :-------------------------------- | :-------: |
|    1 | User register                     |     ✅     |
|    2 | User Login (Authorization by JWT) |     ✅     |
|    3 | Post (Add and get list)           |     ✅     |
|    4 | API Document                      |     ❌     |
|    5 | Docker                            |     ✅     |

## Usage

PostgreSQL is required.

### Docker

```shell
$ docker build -t byoj .
$ docker run -d -v ${PWD}:/app -p 3435:3435 byoj:0.2
```

### Build from Source Code

Modify `config.yml`:

```yaml
database:
    hostname: 127.0.0.1 # host.docker.internal
```

Then run

```shell
$ go mod tidy
$ go build -o byoj .
$ ./byoj
```

## Development

Using following command to commit:

```shell
$ git add .
$ npm run commit
```

## Project Logs

See <https://ligen.life/2022/byoj-project-logs>
