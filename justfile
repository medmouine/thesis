#!/usr/bin/env just --unstable --justfile
set shell := ["zsh", "-cu"]
set dotenv-load := true
export WORKSPACE_ROOT := / justfile_directory()
export TARGET := env_var_or_default('JUST_TARGET', 'mapper')

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

test *T=TARGET: (clean T)
    cd {{T}} && gotestsum -- -coverprofile=cover.out ./...

build-dev *args: (lintall "mapper")
    cd mapper && docker build -f ./dev.Dockerfile {{args}} .

build *args: (lintall "mapper")
    cd mapper && docker build -f ./Dockerfile {{args}} .

push img *args:
    docker push {{img}} {{args}}

dpush img:
    just build "-t {{img}}"
    just push {{img}}

compose-up:
    docker compose -f deploy/docker-compose.yml -p local up -d

compose-down:
    docker compose -f deploy/docker-compose.yml -p local down

kind-up:
    kind create cluster --config deploy/local-cluster.yaml --name $CLUSTER_NAME

kind-down:
    kind delete clusters $CLUSTER_NAME

cloud-up:
    deploy/cloud-infra/doctl_up

cloud-down:
    deploy/cloud-infra/doctl_down
