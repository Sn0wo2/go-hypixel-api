FROM --platform=$BUILDPLATFORM alpine AS certs
RUN apk add --no-cache ca-certificates

FROM scratch
COPY --from=certs /etc/ssl/certs/ca-certificates.crt /etc/ssl/certs/ca-certificates.crt

ARG TARGETPLATFORM
COPY $TARGETPLATFORM/go-hypixel-api /app/go-hypixel-api
WORKDIR /app
ENTRYPOINT ["/app/go-hypixel-api"]
