# Build context must be the repo root: `docker build -f api.Dockerfile .`
FROM node:20 AS web-builder
WORKDIR /web
COPY web/package.json web/package-lock.json ./
RUN npm ci
COPY web/ .
RUN npm run build

FROM golang:1.26 AS builder
WORKDIR /src
COPY go.mod go.sum ./
RUN go mod download
COPY . .
COPY --from=web-builder /web/dist internal/webassets/dist
RUN CGO_ENABLED=0 go build -o /out/api ./cmd/api

FROM gcr.io/distroless/static-debian12
COPY --from=builder /out/api /api
ENTRYPOINT ["/api"]
