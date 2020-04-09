FROM golang:1-alpine AS builder

# Get ca-certificates, zoneinfo and git
RUN apk --update add ca-certificates tzdata

# All these steps will be cached
RUN mkdir /app
WORKDIR /app
COPY go.mod . 
COPY go.sum .

# Get dependancies - will also be cached if we won't change mod/sum
RUN go mod download

# COPY the source code as the last step
COPY . .

# Compile
RUN CGO_ENABLED=0 GOOS=linux GOARCH=amd64 go build -a -installsuffix cgo -o /go/bin/app ./cmd/bot

# Image
FROM scratch

COPY --from=builder /usr/share/zoneinfo /usr/share/zoneinfo
COPY --from=builder /etc/ssl/certs/ca-certificates.crt /etc/ssl/certs/

COPY --from=builder /go/bin/app /app

ENTRYPOINT ["/app"]