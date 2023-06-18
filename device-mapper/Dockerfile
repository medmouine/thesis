FROM alpine:3.14
RUN mkdir -p kubeedge
COPY ./bin/tempsim kubeedge/bin/
COPY ./res kubeedge/res/
WORKDIR kubeedge/bin
ENTRYPOINT ["./tempsim"]