# build stage
FROM golang:1.16-alpine AS build-env
ADD . /src
RUN cd /src && go mod download
RUN cd /src && CGO_ENABLED=0 GOOS=linux GOARCH=amd64 go build -o app

# final stage
FROM alpine
COPY --from=build-env /etc/ssl/certs/ca-certificates.crt /etc/ssl/certs/
COPY --from=build-env /src/app /
COPY --from=build-env /src/config.toml /