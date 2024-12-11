FROM golang:alpine

WORKDIR /func

COPY . .

RUN go build -o binary .

# Command to run the executable
ENTRYPOINT ["/func/binary"]