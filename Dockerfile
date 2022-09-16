FROM golang:1.18
ENV GOPROXY=https://goproxy.cn
WORKDIR /build
COPY ./src .
RUN go mod tidy
RUN GO_EVN=${GO_EVN} go build

FROM  alpine:latest
RUN mkdir -p /cmd
WORKDIR  /cmd
COPY  --from=builder /build/main  .
EXPOSE 8000
CMD ["./main"]