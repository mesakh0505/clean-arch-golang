# clean-arch-golang
Sample repository of implementing [Clean Architecture](https://blog.cleancoder.com/uncle-bob/2012/08/13/the-clean-architecture.html) in Golang

## Documentations and Endpoint
The Open API 3 Documentations can be seen in [docs/openapi.yaml](/docs/openapi.yaml)

## Development

### Prerequisite
- Docker & Docker-compose
- Unix-Family Workspace (Mac, Linux)
- Go 1.13+ (support go module)
- Protobuf
- protoc-gen-go

### Checkout
First, clone the application to any directory as long as not in GOPATH (since this project use Go Module, so clonning in GOPATH is not recommended)

```
$ git clone git@github.com:LieAlbertTriAdrian/clean-arch-golang.git your_target_directory

$ cd your_target_directory
```

### Install The Development Tools
To install the development tools like linter, mock generator, migrator etc you can just run this command

```
$ make init
$ make mockery-prepare
```

#### Linter
To run the linter

```
$ make lint
```

The linter will running and check your code directly.

#### Mock Generator
To generate the mock function based on interface you can do this command

```
$ mockery-gen
```

And if you want to add another interface you can just edit the Makefile command and add your Interface name under the `mockery-gen` command.

```
mockery-gen:
  $(GOPATH)/bin/mockery --name ITodoUsecase
  $(GOPATH)/bin/mockery --name ITodoRepository
  $(GOPATH)/bin/mockery --name IYourNewInterface
```

### Testing

This project has 2 kind of testing
 - unit-testing
 - integration testing to database

#### Unit Testing

To run the unit testing you can do it with this comand

```
$ make short-test
```
#### Integration Testing to Database
And to run the full testing (unit testing plus integration testing with database), you can run it with this command

```
# For local environment. Will spawn the dockerized database and run the test againts it in your local
$ make full-test-local

# For local test with docker. Instead run the test in your local, you can run the test in the docker with all the dependencies such as Database
$ make docker-test
```

## Running All the Application Services

Please follow this steps to run the application

### Dockerize
To build your local Docker image, you need to set a private key that connected to your LieAlbertTriAdrian Github account on `.ssh/id_rsa`

```
$ make docker
```

### Configuration File
To run the service in docker-compose, please create a `.env` in the root project. The configuration example can be looked in .env.example. Or you can just copy this to run locally.

_.env_
```
ENV=debug
POSTGRES_HOST=todo-service-postgress
POSTGRES_PORT=5432
POSTGRES_USER=user
POSTGRES_PASSWORD=password
POSTGRES_DATABASE=todo_app
CONTEXT_TIMEOUT=2s
SERVER_ADDRESS=:9090
DB_MAX_CONN_LIFE_TIME_S=300
DB_MAX_OPEN_CONNECTION=100
DB_MAX_IDLE_CONNECTION=10
```

### Run the Service

After making the environment file, you can just run this command to make it run and live in docker-compose
```
$ make run
```
After running the service, before to test it, you need to run the migration script. You just simply do this command (For Local development).

```
$ make migrate-up
```

```
$ make migrate-create name=haha 
```

Generate protobuf file
```
$ make proto
```

## Deployment
TODO(LieAlbertTriAdrian)

---
# Roadmap
## [Goal]
- Great development experience(Solid tooling, Intuitive structure, Easy to customize)
- Best practices
- Ease of starting a new project

## [Scope/Reqs]
- [x] Add postgres sample
- [x] Add linter
- [x] Add mockery
- [ ] Add mongo sample
- [x] Add rest sample
- [x] Add grpc sample
- [ ] Add event bus sample e.g. rabbitmq / sqs
- [ ] Add postgres testing
- [ ] Add mongo testing
- [ ] Add rest testing
- [ ] Add grpc testing
- [ ] Add event bus testing
- [ ] Add observability - tracing
- [ ] Add error logging - sentry
- [ ] Add e2e testing
- [ ] Add load testing
- [ ] Add grpc gateway
- [ ] Add watchmode - debugger vscode
- [ ] Add deployment
- [ ] Integrate openAPI validation
- [ ] Add graphql
- [ ] Fix and check for typos in documentation
- [ ] Add nix-shell

## References
- https://medium.com/yemeksepeti-teknoloji/mocking-an-interface-using-mockery-in-go-afbcb83cc773
- https://medium.com/@OmisNomis/creating-an-rpc-server-in-go-3a94797ab833
- https://medium.com/swlh/using-grpc-and-protobuf-in-golang-9c218d662db3
- https://grpc.io/docs/languages/go/quickstart/