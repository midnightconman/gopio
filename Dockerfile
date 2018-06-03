FROM golang:1.10.2 as build

ARG VERSION="0.1.0"

ENV APP="gopio"
ENV APP_DIR="/go/src/github.com/midnightconman/${APP}"
ENV GOPIO_PORT="8443"

# Copy source directories individually for better cache
COPY ./vendor ${APP_DIR}/vendor
COPY ./pb ${APP_DIR}/pb
COPY ./schema ${APP_DIR}/schema
COPY ./server ${APP_DIR}/server
COPY ./Gopkg.lock ${APP_DIR}/Gopkg.lock
COPY ./Gopkg.toml ${APP_DIR}/Gopkg.toml

RUN set -ex \
      && cd ${APP_DIR}/server \
      && GOOS=linux GOARCH=arm CGO_ENABLED=0 go install -a \
           -ldflags "-s -w -extldflags \"-static\"" \
      && chmod 0755 /go/bin/linux_arm/server

FROM scratch

COPY --from=build /go/bin/linux_arm/server /

EXPOSE 8080
ENTRYPOINT ["/server"]
CMD [""]
