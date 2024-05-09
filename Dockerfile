FROM golang:1.22.0 as builder
COPY . .
RUN go mod download


ENTRYPOINT ["top", "-b"]