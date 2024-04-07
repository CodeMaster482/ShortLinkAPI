#Builder
FROM golang:1.21.0-alpine AS builder

COPY . /github.com/ShortLinkAPI/
WORKDIR /github.com/ShortLinkAPI/

RUN go mod download
RUN go clean --modcache
RUN CGO_ENABLED=0 GOOS=linux go build -mod=readonly -o server ./cmd/ShortLinkAPI/main.go

#Runner
FROM scratch AS run

WORKDIR /docker-filmLibrary/

COPY --from=builder /github.com/ShortLinkAPI/server .

EXPOSE 8080

ENTRYPOINT ["./server"]