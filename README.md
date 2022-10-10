# FIFA World Cup Table (GO Example)

The responsibility of this project is to store the teams, games and results of the FIFA World Cup. Its rules were taken from Wikipedia on the page about the [FIFA Word Cup 2022](https://en.wikipedia.org/wiki/2022_FIFA_World_Cup).

This is an project example used as a reference for GO training. ***So don`t use this project structure as a silver bullet for any kind of problems.***

**Implemented features and requirements:**
- [x] Register teams;
- [ ] Generate group stage games;
- [ ] Register game results;
- [ ] Calculate team scores and ranking in the group stage;
- [ ] Generate Knockout stage games.

**Used Libraries**
- [httprouter](github.com/julienschmidt/httprouter) as HTTP handler;
- [mongo-driver](go.mongodb.org/mongo-driver) as MongoDB driver.

**Development Tools**
- [testify](github.com/stretchr/testify) for assertion on tests;
- [golangci-lint](https://github.com/golangci/golangci-lint) for run linters;
- [nancy](https://github.com/sonatype-nexus-community/nancy) to check for vulnerabilities.

**Structural Inspiration**
- **Concepts** was inspired by the *Domain Driven Design*;
- **Layout** was inspired by the [golang-standards](https://github.com/golang-standards/project-layout);
- **Layers Structure** was inspired by the *Hexagonal Architecture* (Ports & Adapters).

## Quick start

Requirements:
- [Golang](http://golang.org/)(>1.19)
- [GNU Make](https://www.gnu.org/software/make/)
- [Docker](http://docker.com)
- [Docker Compose](http://docker.com)

All main commands that you need to execute the project were added to Makefile as a single target. Here is the `make` or `make help` infos to help you to get started.

```console
~$ make help
usage: make [target]

build:
  build                           Build image.
  image                           Create release docker image.
  tag                             Add docker tag in release image.

check:
  lint                            Run lint on docker.
  audit                           Run audit on docker.
  test-unit                       Run unit tests on docker.
  test-integration                Run integration tests on docker.

default:
  help                            Show this help.

dependencies:
  clean                           Remove locally generated files.
  install                         Download dependencies.

run:
  run                             Run application on docker compose.
  stop                            Stop application running on docker compose.
```

That's it. Now you are good to go.

### Environment variables

| Name                  | Required  | Default Value           |
|-----------------------|-----------|-------------------------|
| MONGO_URL             | yes       | -                       |
| MONGO_DATABASE_NAME   | no        | `fifa-wct-go-example`   |
| HTTP_SERVER_PORT      | no        | `8080`                  |
| HTTP_SERVER_HOST      | no        | `localhost`             |

## Contributing

Pull requests are welcome. For major changes, please open an issue first to discuss what you would like to change.

Please make sure to update tests as appropriate.

## License

This project is available under the [MIT](https://choosealicense.com/licenses/mit/) license. See the LICENSE file for more info.
