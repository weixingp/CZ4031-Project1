FROM golang:latest

RUN mkdir /App
ADD . /app
WORKDIR /app
RUN go build -o ./main ./exp.go

CMD ["./main"]
#docker build -t gobptree .
#docker run gobptree