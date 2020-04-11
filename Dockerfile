FROM golang:1.11.5-alpine

WORKDIR /go/src/github.com/ravern/jabba

COPY . .

RUN apk add --no-cache git

RUN go get -u github.com/gobuffalo/packr/packr

RUN go get -u -d github.com/magefile/mage && \
  cd /go/src/github.com/magefile/mage && \
  go run bootstrap.go

RUN mage prod

FROM alpine:3.9

ENV HOSTNAME jabba.ravern.co
ENV PORT 80

WORKDIR /app

COPY --from=0 /go/src/github.com/ravern/jabba/releases/jabba .

CMD [ "./jabba" ]
