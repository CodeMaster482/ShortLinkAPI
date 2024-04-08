#Builder
FROM golang:1.21.0-alpine AS builder

COPY . /github.com/ShortLinkAPI/
WORKDIR /github.com/ShortLinkAPI/

RUN go mod download
RUN go clean --modcache
RUN CGO_ENABLED=0 GOOS=linux go build -mod=readonly -o server ./cmd/ShortLinkAPI/main.go

RUN chmod +x ./server

#Runner
#scratch
FROM scratch AS run 

WORKDIR /docker-ShortLinkAPI/

COPY --from=builder /github.com/ShortLinkAPI/server .

EXPOSE 8080 9000

ARG memo
ENV memo ${MEMO}

ENTRYPOINT ["./server", "--in-memo", "false"]