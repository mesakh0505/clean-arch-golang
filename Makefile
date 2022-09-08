ifndef $(GOPATH)
    GOPATH=$(shell go env GOPATH)
    export GOPATH
endif

# Database
POSTGRES_USER ?= user
POSTGRES_PASSWORD ?= password
POSTGRES_ADDRESS ?= localhost:5432
POSTGRES_DATABASE ?= todo_app

.PHONY: engine
engine:
	go build -o engine app/main.go 

# Install all the tools for development
.PHONY: init
init: lint-prepare mockery-prepare migrate-prepare

# Install the mockery. This command will install the mockery in the GOPATH/bin folder
mockery-prepare: 
	 @go get github.com/vektra/mockery
	 @go get github.com/vektra/mockery/.../

# Use the mockery to generate mock interface
mockery-gen:
	# $(GOPATH)/bin/mockery --name ITodoUsecase
	# $(GOPATH)/bin/mockery --name ITodoRepository
	mockery --dir=domain --name ITodoUsecase 
	mockery --dir=domain --name ITodoRepository


.PHONY: lint-prepare
lint-prepare:
	@echo "Preparing Linter"
	curl -sfL https://raw.githubusercontent.com/golangci/golangci-lint/master/install.sh| sh -s latest

.PHONY: lint
lint:
	@echo "Applying linter"
	./bin/golangci-lint run ./...

.PHONY: run
run:
	@docker-compose up -d --build

.PHONY: stop
stop:
	@docker-compose down

# Short test is used for testing the whole unit-test
.PHONY: short-test
short-test:
	@go test -v -cover --short -race ./...

# Full test is used for testing the whole application including the database query directly to a live database
# This may takes time.
.PHONY: full-test
full-test:
	@echo "Running the full test..."
	@go test -v -cover -race ./...

.PHONY: full-test-local
full-test-local:
	@docker-compose -f docker-compose.test.yaml up -d postgres-test
	@make full-test
	@docker-compose -f docker-compose.test.yaml down --volumes

.PHONY: docker-test
docker-test:
	@docker-compose -f docker-compose.test.yaml up --build --abort-on-container-exit
	@docker-compose -f docker-compose.test.yaml down --volumes

.PHONY: migrate-prepare
migrate-prepare:
	@go get -tags 'postgres' -u github.com/golang-migrate/migrate/v4/cmd/migrate
	@go build -a -o ./bin/migrate -tags 'postgres' github.com/golang-migrate/migrate/v4/cmd/migrate

.PHONY: migrate-up
migrate-up:
	@bin/migrate -database "postgres://$(POSTGRES_USER):$(POSTGRES_PASSWORD)@$(POSTGRES_ADDRESS)/$(POSTGRES_DATABASE)?sslmode=disable" \
	-path=internal/postgres/migrations up

.PHONY: migrate-down
migrate-down:
	@bin/migrate -database "postgres://$(POSTGRES_USER):$(POSTGRES_PASSWORD)@$(POSTGRES_ADDRESS)/$(POSTGRES_DATABASE)?sslmode=disable" \
	-path=migrations/ down

.PHONY: migrate-drop
migrate-drop:
	@bin/migrate -database "postgres://$(POSTGRES_USER):$(POSTGRES_PASSWORD)@$(POSTGRES_ADDRESS)/$(POSTGRES_DATABASE)?sslmode=disable" \
	-path=migrations/ drop

.PHONY: migrate-create
migrate-create:
	@bin/migrate create -ext sql -dir internal/postgres/migrations ${name}

.PHONY: clean
clean: 
	@make stop
	@docker-compose -f docker-compose.test.yaml down --volumes

.PHONY: proto
proto:
	@protoc --go_out=. --go_opt=paths=source_relative \
		--go-grpc_out=. --go-grpc_opt=paths=source_relative \
		domain/proto/*.proto