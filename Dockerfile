###############################################################################
# roboll/cfssl-sidecar
###############################################################################
FROM alpine:3.2

ADD cfssl-sidecar-linux-amd64 /cfssl-sidecar
ENTRYPOINT ["/cfssl-sidecar"]
