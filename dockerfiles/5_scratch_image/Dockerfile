FROM golang:1.21 as builder

WORKDIR /usr/src/app
RUN useradd -u 10001 goapp

# pre-copy/cache go.mod for pre-downloading dependencies and only redownloading them in subsequent builds if they change
COPY ./go.mod ./go.sum ./
RUN go mod download && go mod verify

COPY . .
RUN CGO_ENABLED=0 GOOS=linux go build -v -o /usr/local/bin/app ./...

FROM scratch
WORKDIR /bin/
COPY --from=builder /etc/passwd /etc/passwd
COPY --from=builder /usr/share/zoneinfo /usr/share/zoneinfo
COPY --from=builder /etc/ssl/certs /etc/ssl/certs
COPY --from=builder /usr/local/bin/app .

USER goapp

EXPOSE 8080

HEALTHCHECK CMD curl --fail http://localhost:8080/health || exit 1
CMD ["./app"]  
