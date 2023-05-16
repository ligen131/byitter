FROM golang:1.20-alpine3.17 as builder
COPY . /src
WORKDIR /src
ENV GOPROXY "https://goproxy.cn"
RUN go build -o /build/byoj .

FROM alpine:3.17 as prod
COPY --from=builder /build/byoj /usr/bin/byoj
WORKDIR /app
ENTRYPOINT [ "byoj" ]