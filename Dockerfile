FROM alpine:3.18.6

COPY tfx /usr/bin/tfx
ENTRYPOINT ["/usr/bin/tfx"]
