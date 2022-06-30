FROM golang:1.18-bullseye as builder
WORKDIR /app
COPY go.* ./
RUN go mod download
COPY . .
RUN CGO_ENABLED=0 GOOS=linux go build -o exporter

FROM gcr.io/distroless/static-debian11:nonroot
COPY --from=builder /app/exporter /app/exporter
USER nonroot
ENTRYPOINT ["/app/exporter"]
