FROM golang:1.19-alpine as builder_base
WORKDIR /build
ENV GO111MODULE=on \
    GOPROXY=https://goproxy.cn,direct
COPY go.mod .
COPY go.sum .
RUN go mod download
COPY . .
# -ldflags="-w -s" means omit DWARF symbol table and the symbol table and debug information
RUN CGO_ENABLED=0 GOOS=linux GOARCH=amd64 go build \
    -o /main \
    ./server/main.go

FROM alpine:latest as runner
WORKDIR /root
COPY --from=builder_base /main .
COPY --from=builder_base /build/data /root/data
CMD ["./main"]