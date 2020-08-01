FROM golang:1.14.5-alpine
ARG WORK_DIR
ENV WORK_DIR=${WORK_DIR:-/go/src/github.com/ulm0/deliver}
COPY . ${WORK_DIR}
WORKDIR $WORK_DIR/cmd/deliver
RUN apk add --no-cache upx && \
    CGO_ENABLED=0 go build -a -ldflags="-s -w -extldflags -static" -installsuffix cgo && \
    upx -q9 deliver

FROM busybox
COPY --from=0 /go/src/github.com/ulm0/deliver/cmd/deliver/deliver /bin/deliver
ADD https://raw.githubusercontent.com/containous/traefik/master/script/ca-certificates.crt /etc/ssl/certs/ca-certificates.crt
