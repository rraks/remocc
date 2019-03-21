FROM golang

ENV GO111MODULE=on

WORKDIR /go/src/github.com/rraks/remocc

COPY go.mod .
COPY go.sum .

RUN go mod download

COPY . .

RUN CGO_ENABLED=0 GOOS=linux GOARCH=amd64 go build ./

CMD ["/go/src/github.com/rraks/remocc/remocc"]
EXPOSE 3000
