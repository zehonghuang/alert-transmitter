FROM golang:1.8
RUN mkdir /app
ADD . /app/
WORKDIR /app
RUN go build -o main .
RUN adduser -S -D -H -h /app appuser
USER appuser
CMD ["GO_ENV=${GO_ENV} ./main"]