# dapr-starter

This Dapr starter project accelerates the development of new Dapr services in `go`. It includes `make` commands for: 

* `test`  - Tests the entire project
* `run`   - Runs the uncompiled code
* `build` - Builds local release binary
* `exec`  - Builds binary and runs it in Dapr locally
* `event` - Publishes sample message to Dapr pubsub API
* `image` - Builds and publish docker image to Docker Hub
* `lint`  - Lints the entire project
* `tag`   - Creates release tag
* `clean` - Cleans all runtime generated directory (bin, vendor)
* `help`  - Display available commands

This project also includes GitHub actions in [.github/workflows](.github/workflows) that test on each `push` and build images and mark release on each `tag`. 

[![Test](https://github.com/mchmarny/dapr-starter/workflows/Test/badge.svg)](https://github.com/mchmarny/dapr-starter/actions?query=workflow%3ATest) ![Release](https://github.com/mchmarny/dapr-starter/workflows/Release/badge.svg?query=workflow%3ARelease) [![Go Report Card](https://goreportcard.com/badge/github.com/mchmarny/dapr-starter)](https://goreportcard.com/report/github.com/mchmarny/dapr-starter) ![GitHub go.mod Go version](https://img.shields.io/github/go-mod/go-version/mchmarny/dapr-starter)

## how to use

1. Start by clicking on the Use this template button above the file list
2. Select the account you want to own the repository in the owner drop-down menu
3. Name your repository and optionally add description
4. Choose a repository visibility (Public, Internal, Public)
5. Finally, click Create repository from template

## make it your own 

After setting up the template, there is a few things you may want to change in your new project

### Makefile

* Change the `SERVICE_NAME` and `RELEASE_VERSION` variables
* If you have not already defined the `DOCKER_USER` environment variable set it directly here

### go.mod

* Update line 1 in [go.mod](go.mod) file by changing the github username org name and the project name to your own (`module github.com/mchmarny/dapr-starter`)
* Run `go mod tidy` and `go mod vendor`
* Run `make build` to make sure everything works (check [bin](bin) folder for results)

### deployment files

If deploying to Kubernates you may also want to consider updating the components and deployment files in the [deploy](deploy) directory. 

## Disclaimer

This is my personal project and it does not represent my employer. I take no responsibility for issues caused by this code. I do my best to ensure that everything works, but if something goes wrong, my apologies is all you will get.

## License

This software is released under the [MIT](./LICENSE)
