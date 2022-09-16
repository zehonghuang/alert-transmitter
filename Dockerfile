FROM golang:1.18.2
RUN mkdir /app
ADD . /app/
WORKDIR /app
RUN go env -w GOPROXY=https://goproxy.cn
RUN go env -w GOSUMDB=off
RUN go mod download
RUN go build -o main .
RUN cd /app
EXPOSE 8000
CMD ["GO_ENV=${GO_ENV} ./main"]