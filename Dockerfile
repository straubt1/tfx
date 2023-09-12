FROM alpine:3.18.3

COPY tfx /usr/bin/tfx
ENTRYPOINT ["/usr/bin/tfx"]
