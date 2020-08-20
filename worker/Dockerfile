FROM golang
MAINTAINER  liujq51
WORKDIR /go/src
#COPY . .
#RUN go build ./worker.go -o worker
ENTRYPOINT ["/go/src/main/worker","-config","/go/src/main/worker.json"]
#ENTRYPOINT ["ls", "-al","/go/src"]
