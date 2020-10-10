###################
# CA Certificates #
###################
FROM alpine:latest as certs
RUN apk --update add ca-certificates

###################
# Build stage     #
###################
FROM golang:1.14 as builder

# Get the GitHub workflow run ID and commit hash passed by GitHub actions
ARG GITHUB_SHA
ARG GITHUB_RUN_ID

# go-sqlite3 requires cgo to work
ENV CGO_ENABLED=1

WORKDIR /app
COPY . .

RUN make build

###################
# Final image #
###################
FROM scratch

COPY --from=certs /etc/ssl/certs/ca-certificates.crt /etc/ssl/certs/ca-certificates.crt
COPY --from=builder /app/build/service /

EXPOSE 8080

ENTRYPOINT ["./service"]
