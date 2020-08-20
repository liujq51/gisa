FROM golang
MAINTAINER  liujq51
WORKDIR /go/src
#COPY . .
EXPOSE 80
#RUN go build ./main/master.go -o master
#CMD ["/go/src/main/master", "-config", "/go/src/main/master.json"]
ENTRYPOINT ["/go/src/backend"]
