ifeq ($(OS),Windows_NT)
    DOCKER_EXEC_FLAGS = -it
else
    ifeq ($(CI),true)
        DOCKER_EXEC_FLAGS = -i
    else
        DOCKER_EXEC_FLAGS = -it
    endif
endif

ifeq ($(OS),Windows_NT)
  RM_BIN = if exist bin rmdir /s /q bin
  MKDIR_BIN = mkdir bin
else
  RM_BIN = rm -rf bin
  MKDIR_BIN = mkdir bin
endif

.PHONY: build generate test docker-build docker-up docker-down setup init docker-mysql start-server start-client

# Build the application
build:
	$(RM_BIN)
	$(MKDIR_BIN)
	go build -o bin/server.exe ./cmd/server
	go build -o bin/client.exe ./cmd/client

# Generate SQLBoiler models
generate:
	buf generate
	go mod tidy
	sqlboiler mysql --wipe --config sqlboiler.toml --output gen/db --pkgname db

# Run tests
test:
	go test -v ./internal/repository

docker-build:
	docker compose build

docker-up:
	docker compose up -d

docker-down:
	docker compose down

docker-mysql:
	@echo Starting MySQL container...
	docker compose up -d mysql
	@echo Waiting for MySQL to be ready...
	@setlocal ENABLEDELAYEDEXPANSION && \
	for /L %%i in (1,1,30) do ( \
		docker exec todo01-mysql-1 mysqladmin ping -h127.0.0.1 -uroot -proot --silent >nul 2>&1 && ( \
			echo MySQL is ready! && exit /b 0 \
		) || ( \
			echo Waiting for MySQL to be ready... && ping -n 3 127.0.0.1 >nul \
		) \
	)

# Development setup
setup:
	go install github.com/volatiletech/sqlboiler/v4@latest
	go install github.com/volatiletech/sqlboiler/v4/drivers/sqlboiler-mysql@latest
	go install github.com/bufbuild/buf/cmd/buf@latest
	go install google.golang.org/protobuf/cmd/protoc-gen-go@latest
	go install github.com/bufbuild/connect-go/cmd/protoc-gen-connect-go@latest
	go mod tidy

test-db:
	echo "test"

start-server:
	./bin/server.exe

client-create:
	./bin/client.exe -action=create -title="Test Task" -description="This is a test" -due-date=2024-07-01

client-list:
	./bin/client.exe -action=list

client-update:
	./bin/client.exe -action=update -id=1 -title="Updated Task" -description="Updated description" -status=completed

client-delete:
	./bin/client.exe -action=delete -id=1



