FROM golang:1.23-bookworm AS builder
WORKDIR /app
COPY . .
RUN make build

FROM gcr.io/distroless/static:nonroot
COPY --from=builder /app/bin/refresher /refresher
ENTRYPOINT ["/refresher"]