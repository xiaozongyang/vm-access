FROM --platform=$BUILDPLATFORM golang:1.24.4 as builder

WORKDIR /workspace
ENV GOPROXY=https://goproxy.cn,direct

COPY go.mod go.sum ./
RUN --mount=type=ssh go mod download

COPY . .

ARG TARGETARCH
RUN CGO_ENABLED=0 GOOS=linux GOARCH=$TARGETARCH go build -a -o vm-access-api cmd/vm-access-api/main.go
RUN CGO_ENABLED=0 GOOS=linux GOARCH=$TARGETARCH go build -a -o vm-access-proxy cmd/vm-access-proxy/main.go

FROM --platform=$BUILDPLATFORM scratch
COPY --from=builder /workspace/vm-access-api /vm-access-api
COPY --from=builder /workspace/vm-access-proxy /vm-access-proxy
COPY --from=builder /workspace/static /static
COPY --from=builder /workspace/templates /templates

CMD ["/vm-access-proxy"]
