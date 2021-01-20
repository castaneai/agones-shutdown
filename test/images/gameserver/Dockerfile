FROM golang:alpine3.12 AS builder
ARG VERSION
RUN apk add --no-cache git
WORKDIR /go/src/gameserver
COPY . /go/src/gameserver
RUN CGO_ENABLED=0 go build -ldflags "-X main.version=$VERSION" -o /bin/gameserver main.go

FROM scratch
COPY --from=builder /bin/gameserver /bin/gameserver
ENTRYPOINT ["/bin/gameserver"]
