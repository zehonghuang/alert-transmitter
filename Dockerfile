# Base image, golang 1.18
FROM golang:1.18.3
ENV GOPROXY=https://goproxy.cn
WORKDIR /workspace
# Copy all files into the image
COPY . .
# Run go mod
RUN go mod download
# Expose ports
EXPOSE 8000
# Run Go program, just like locally
ENTRYPOINT ["GO_EVN=${GO_EVN}","go","run","main.go"]