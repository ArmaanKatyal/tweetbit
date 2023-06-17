FROM golang:1.20 as builder

WORKDIR /app

COPY go.mod go.sum ./

RUN go mod download

COPY . .

RUN go build -o main

RUN chmod +x main

CMD [ "/app/main" ]
