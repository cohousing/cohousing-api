FROM golang:1.8-alpine
MAINTAINER Frank Bille-Stauner <frank@cohousing.nu>

WORKDIR /go/src/github.com/cohousing/cohousing-api
COPY . .

RUN go-wrapper download
RUN go-wrapper install

EXPOSE 8080

ENV GIN_MODE release

ENTRYPOINT ["go-wrapper", "run"]