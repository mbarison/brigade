FROM alpine:3.7

RUN apk update && apk add --no-cache \
    ca-certificates \
    && update-ca-certificates

COPY rootfs/brigade-ai-gateway /usr/bin/brigade-ai-gateway

#ENV

CMD /usr/bin/brigade-ai-gateway
