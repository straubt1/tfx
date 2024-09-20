FROM alpine:3.20.3

COPY tfx /usr/bin/tfx
ENTRYPOINT ["/usr/bin/tfx"]
