FROM golang:1.18.2
RUN mkdir /app
ADD . /app/
WORKDIR /app
RUN go env -w GOPROXY=https://goproxy.cn
RUN go env -w GOSUMDB=on
RUN go mod download
RUN go build -o main .
RUN adduser -S -D -H -h /app appuser
USER appuser
CMD ["GO_ENV=${GO_ENV} ./main"]