set dotenv-load
export WORKSPACE_ROOT := / justfile_directory()
export TARGET := env_var_or_default('JUST_TARGET', '')
export KIND_CLUSTER_NAME := env_var_or_default('CLUSTER_NAME', 'local')
export KO_DOCKER_REPO := env_var_or_default('KO_DOCKER_REPO', 'ko.local')

export TARGET_DISK_PATH := "./" + TARGET
export TARGET_MOD_PATH := ""

mapper:
    just ctx mapper

lintall *T=TARGET: lint sec

sec *T=TARGET:
    cd {{T}} && gosec -exclude=G401,G404,G501,G502,G505 ./...

#critic *T=TARGET:
#    cd {{T}} && gocritic check -enableAll ./...

lint *T=TARGET:
    cd {{T}} && golangci-lint run ./...

clean *T=TARGET:
    rm -rf ${TARGET}/build

build *T=TARGET: (clean T)
    cd {{T}} && CGO_ENABLED=0 go build -ldflags="-w -s" -o ./build/{{T}} ./main.go

test *T=TARGET: (clean T)
    cd {{T}} && gotestsum -- -coverprofile=cover.out ./...

t := ""
cbuild *T=TARGET: (lintall T)
    @ftag="{{ if t != "" { "-t " + t } else { "" } }}" && cd {{T}} && ko build $ftag .

compose-up:
    docker compose -f deploy/docker-compose.yml -p local up -d

compose-down:
    docker compose -f deploy/docker-compose.yml -p local down

kind-up:
    KIND_EXPERIMENTAL_PROVIDER=podman kind create cluster --config config/local-cluster.yaml --name ${KIND_CLUSTER_NAME}

kind-down:
    kind delete clusters ${KIND_CLUSTER_NAME}

ip:
    podman inspect ${TARGET} --format {{{{.NetworkSettings.IPAddress}}
