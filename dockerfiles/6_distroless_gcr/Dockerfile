FROM golang:1.21 as builder

WORKDIR /usr/src/app

# pre-copy/cache go.mod for pre-downloading dependencies and only redownloading them in subsequent builds if they change
COPY ./go.mod ./go.sum ./
RUN go mod download && go mod verify

COPY . .
RUN CGO_ENABLED=0 GOOS=linux go build -v -o /usr/local/bin/app ./...


FROM gcr.io/distroless/static:nonroot
WORKDIR /bin/
COPY --from=builder /usr/local/bin/app .

EXPOSE 8080

HEALTHCHECK CMD curl --fail http://localhost:8080/health || exit 1
CMD ["/bin/app"]  
