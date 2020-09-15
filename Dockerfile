FROM golang:1.15-alpine AS build_base

WORKDIR /tmp/app

COPY app/ .

RUN go mod download

RUN go build -o ./out/kube-client .

RUN ls -alt ./out/kube-client

FROM alpine:3.9

COPY --from=build_base /tmp/app/out/kube-client /app/kube-client

CMD ["/app/kube-client"]