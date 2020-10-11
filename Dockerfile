###################
# Build stage     #
###################
FROM golang:1.14-alpine as builder

# Get the GitHub workflow run ID and commit hash passed by GitHub actions
ARG GITHUB_SHA
ARG GITHUB_RUN_ID

# go-sqlite3 requires cgo to work. That's unfortunate because it means we can't use a scratch image which would have
# made the image as small as possible while also making it more secure by not even having a shell.
# We'll accept this for this demo since it wouldn't be a problem when speaking to a real database.
ENV CGO_ENABLED=1

RUN apk --update add ca-certificates make build-base sqlite

WORKDIR /app
COPY . .

RUN make build

###################
# Final image #
###################
#FROM scratch
#
#COPY --from=builder /etc/ssl/certs/ca-certificates.crt /etc/ssl/certs/ca-certificates.crt
#COPY --from=builder /app/build/service /

EXPOSE 8080

ENTRYPOINT ["build/service"]
