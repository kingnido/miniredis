FROM golang:1.11-alpine
RUN mkdir /app
ADD . /app
WORKDIR /app
RUN go build -o redis .
CMD ["/app/redis", "-port=8000"]
