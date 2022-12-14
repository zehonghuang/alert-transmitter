# Base image, golang 1.18
FROM golang:1.18.3
ENV GOPROXY=https://goproxy.cn
WORKDIR /workspace
# Copy all files into the image
COPY . .
COPY ./resource/* ./resource/
COPY ./template/* ./template/
# Run go mod
RUN go mod download
# Expose ports
EXPOSE 8000
# Run Go program, just like locally
ENTRYPOINT ["go","run","."]