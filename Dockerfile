FROM golang:1.19 as builder

WORKDIR /usr/src/app

# pre-copy/cache go.mod for pre-downloading dependencies and only redownloading them in subsequent builds if they change
COPY go.mod go.sum ./
RUN go mod download && go mod verify

COPY . .
RUN CGO_ENABLED=0 GOOS=linux go build -v -o /usr/local/bin/app ./...
# RUN CGO_ENABLED=0 GOOS=linux go build -a -installsuffix cgo -o /usr/local/bin/app .

FROM scratch
WORKDIR /bin/
COPY --from=builder /usr/local/bin/app .
COPY --from=builder /etc/ssl/certs/ca-certificates.crt /etc/ssl/certs/ca-certificates.crt
COPY --from=builder /etc/passwd /etc/passwd

EXPOSE 8080

USER app
HEALTHCHECK CMD curl --fail http://localhost:8080 || exit 1
CMD ["./app"]  

# FROM alpine:latest  
# RUN apk --no-cache add ca-certificates libc6-compat
# WORKDIR /app/
# COPY --from=builder /usr/local/bin/app .
# EXPOSE 8080
# CMD ["./app"]  
