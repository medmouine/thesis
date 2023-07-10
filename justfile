set dotenv-load
export TARGET := env_var_or_default('JUST_TARGET', 'all')
export KIND_CLUSTER_NAME := env_var_or_default('CLUSTER_NAME', 'local')

mapper:
    just ctx target=mapper

@default:
    echo "Current Target: ${TARGET}"
    echo "___________________________"

@ctx target:
    echo "Set target to {{target}}"
    sed -i '' "s/JUST_TARGET=.*$/JUST_TARGET={{target}}/" "{{justfile_directory()}}/.env"

lintall: default lint sec critic

sec:
    cd ${TARGET} && gosec ./...

critic:
    cd ${TARGET} && gocritic check -enableAll ./...

lint:
    cd ${TARGET} && golangci-lint run ./...

clean:
    rm -rf ${TARGET}/build

build $APP_NAME='app': lintall clean
    export CGO_ENABLED=0
    go build -ldflags="-w -s" -o ${TARGET}/build/${APP_NAME} ${TARGET}/main.go

test: lintall clean
    cd ${TARGET} && go test -v ./...

cbuild *tag='latest':
    cd ${TARGET} && podman build -t ${TARGET}:{{tag}} .

compose:
    docker compose -f docker-compose.yml ${TARGET} up

local-up:
    KIND_EXPERIMENTAL_PROVIDER=podman kind create cluster --config config/local-cluster.yaml --name ${KIND_CLUSTER_NAME}

local-down:
    kind delete clusters ${KIND_CLUSTER_NAME}

ip:
    podman inspect ${TARGET} --format {{{{.NetworkSettings.IPAddress}}


