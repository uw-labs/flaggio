######################################
# STEP 1 build go executable binary
######################################
FROM golang:1.14-alpine AS go_builder

RUN apk update && apk add --no-cache git make ca-certificates tzdata && update-ca-certificates
WORKDIR /flaggio

# Create appuser.
RUN adduser -D -g '' flaggio

COPY . .

# Fetch dependencies and build the binary
RUN make install && \
    GOOS=linux GOARCH=amd64 make build

######################################
# STEP 2 build frontend app
######################################
FROM node:12-alpine AS node_builder

ENV NODE_ENV production

WORKDIR /flaggio

COPY --from=go_builder /flaggio/web /flaggio

RUN npm install && npm run build

######################################
# STEP 3 build image
######################################
FROM scratch

COPY --from=go_builder /usr/share/zoneinfo /usr/share/zoneinfo
COPY --from=go_builder /etc/ssl/certs/ca-certificates.crt /etc/ssl/certs/
COPY --from=go_builder /etc/passwd /etc/passwd

COPY --from=go_builder /flaggio/bin/flaggio /flaggio
COPY --from=node_builder /flaggio/build /

USER flaggio

EXPOSE 8080
EXPOSE 8081

ENTRYPOINT ["/flaggio"]
