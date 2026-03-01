FROM alpine:3.20.3

ARG TARGETPLATFORM
COPY ${TARGETPLATFORM}/tfx /usr/bin/tfx
ENTRYPOINT ["/usr/bin/tfx"]
