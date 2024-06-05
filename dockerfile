FROM golang:latest

WORKDIR /app

COPY . .
RUN go mod download && go build -o social-network-api cmd/main.go

CMD [ "./social-network-api", "-port", "5000" ]