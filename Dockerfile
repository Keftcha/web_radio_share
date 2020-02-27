FROM golang:1.14 as builder

WORKDIR /go/src/app
COPY . .

RUN go build -o /bin/app


FROM scratch

WORKDIR /bin/app
COPY --from=builder /bin/app .

CMD ["./app"]
