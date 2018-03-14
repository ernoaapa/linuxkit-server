
FROM linuxkit/alpine:07f7d136e427dc68154cd5edbb2b9576f9ac5213 as mirror

RUN mkdir -p /out/etc/apk && cp -r /etc/apk/* /out/etc/apk/
RUN apk add --no-cache --initdb -p /out \
  ca-certificates \
  dosfstools \
  multipath-tools
RUN apk add --no-cache --initdb -p /out --repository http://dl-3.alpinelinux.org/alpine/edge/main \
  parted
RUN apk add --no-cache --initdb -p /out --repository http://dl-3.alpinelinux.org/alpine/edge/community \
  docker
RUN rm -rf /out/etc/apk /out/lib/apk /out/var/cache

FROM scratch
COPY --from=mirror /out/ /

LABEL maintainer="Erno Aapa <erno.aapa@gmail.com>"

COPY linuxkit-server /

ENTRYPOINT ["/linuxkit-server"]