FROM golang:1.17 as builder

WORKDIR /app

COPY . .

#RUN CGO_ENABLED=0 GOOS=linux go build -a -installsuffix cgo -o vstorage-csi .
RUN CGO_ENABLED=0 GOOS=linux go build -a -installsuffix cgo -o vstorage-csi main.go 

FROM alpine:latest  
RUN apk --no-cache add ca-certificates

WORKDIR /root/

COPY --from=builder /app/vstorage-csi .

CMD ["./vstorage-csi"]

