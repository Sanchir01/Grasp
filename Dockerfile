FROM golang:1.22.0 as builder
COPY . .
RUN go mod download
RUN go build -o ./.bin/main ./cmd/main/main.go


CMD ["go", "-b"]