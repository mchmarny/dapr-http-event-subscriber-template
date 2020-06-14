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