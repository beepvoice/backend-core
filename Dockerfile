FROM golang:1.13-rc-alpine as build

RUN apk add --no-cache git

WORKDIR /src
COPY go.mod go.sum .env *.go ./
RUN CGO_ENABLED=0 go build -ldflags "-s -w"

FROM scratch

COPY --from=build /src/core /core
COPY --from=build /src/.env /.env

ENTRYPOINT ["/core"]
