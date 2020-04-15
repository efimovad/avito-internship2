FROM golang:1.13-stretch AS builder

WORKDIR /usr/src/app

COPY go.mod .
RUN go mod download

COPY . .
RUN make build

FROM ubuntu:19.04
ENV DEBIAN_FRONTEND=noninteractive
ENV PORT 8080
EXPOSE $PORT

RUN apt-get update && apt-get install -y build-essential

COPY --from=builder /usr/src/app/ .

CMD ./server