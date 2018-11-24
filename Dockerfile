# Note: The newer busybox:glibc is missing libpthread.so.0.
FROM alpine:latest
MAINTAINER iYOCHU Nigeria Ltd

RUN apk add --no-cache ca-certificates

# Add the executable
COPY bittrexmicro /bittrex/bittrexmicrox




EXPOSE 5030

CMD ["/bittrex/bittrexmicrox"]