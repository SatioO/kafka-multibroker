FROM golang:alpine as builder

RUN apk update && apk add --no-cache git

WORKDIR $GOPATH/src/consumer
COPY go.mod go.sum ./

RUN go mod download

COPY . .

RUN CGO_ENABLED=0 GOOS=linux GOARCH=amd64 go build -a -installsuffix cgo -ldflags="-w -s" -o /bin/consumer main.go

FROM scratch
COPY --from=builder /bin/consumer /bin/consumer
ENTRYPOINT ["/bin/consumer"]