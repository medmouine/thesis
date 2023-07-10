FROM alpine:3.14

RUN mkdir -p kubeedge

COPY ./bin/main kubeedge/bin/mapper

VOLUME kubeedge/res/
VOLUME /opt/kubeedge/deviceProfile.json

WORKDIR kubeedge/bin

ENTRYPOINT ["./mapper"]

