#FROM golang:1.20-alpine AS builder
#
## Move to working directory (/build).
#WORKDIR /build
#
## Copy and download dependency using go mod.
#COPY go.mod go.sum ./
#RUN go mod download
#
## Copy the code into the container.
#COPY . .
#
## Set necessary environment variables needed for our image and build the API server.
#ENV CGO_ENABLED=0 GOOS=linux GOARCH=amd64
#RUN go build -ldflags="-s -w" -o apiserver .
#
#FROM scratch
#
## Copy binary and config files from /build to root folder of scratch container.
#COPY --from=builder ["/build/apiserver", "/build/.env", "/"]
#
## Command to run when starting the container.
#ENTRYPOINT ["/apiserver"]


FROM golang:latest

ENV PROJECT_DIR=/go/src/mapper \
    GO111MODULE=on \
    CGO_ENABLED=0

WORKDIR /go/src/mapper
COPY . .

RUN go mod download -x
RUN go get github.com/githubnemo/CompileDaemon
RUN go install github.com/githubnemo/CompileDaemon

EXPOSE 3000

ENTRYPOINT CompileDaemon -build="go build -o ./build/mapper main.go " -command="./build/mapper"
