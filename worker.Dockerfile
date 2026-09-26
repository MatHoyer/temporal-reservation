# Build context must be the repo root: `docker build -f worker.Dockerfile .`
FROM golang:1.26 AS builder
WORKDIR /src
COPY go.mod go.sum ./
RUN go mod download
COPY . .
RUN CGO_ENABLED=0 go build -o /out/worker ./cmd/worker

FROM gcr.io/distroless/static-debian12
COPY --from=builder /out/worker /worker
ENTRYPOINT ["/worker"]
