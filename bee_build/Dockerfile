FROM golang:1.15.0
RUN go get -u github.com/astaxie/beego
RUN go get -u github.com/lib/pq
RUN export GO111MODULE=on && go get github.com/beego/bee
WORKDIR /go/src