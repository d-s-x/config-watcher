ARG BASEIMAGE=busybox
FROM $BASEIMAGE

USER 65534

ARG BINARY=config-watcher
COPY out/$BINARY /config-watcher

ENTRYPOINT ["/config-watcher"]
