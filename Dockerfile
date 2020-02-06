FROM golang:1.13-alpine AS builder
WORKDIR /cleaner
ENV CGO_ENABLED=0
COPY *.go go.mod go.sum ./
RUN go build -ldflags="-s -w"

ENTRYPOINT ["/bin/sh", "-c", "env && /cleaner/artifact-cleaner"]
#FROM gcr.io/distroless/static
#ENTRYPOINT ["/usr/bin/cleaner"]
#COPY --from=builder /cleaner/artifact-cleaner /usr/bin/cleaner
