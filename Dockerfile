FROM golang:1.16-alpine

WORKDIR /app

# Download Go modules
COPY go.mod .
COPY go.sum .

RUN go mod download

COPY *.go ./

# Build
RUN go build -o /diagram

EXPOSE 8080

CMD [ "/diagram" ]
