FROM golang:1.11.5-alpine

WORKDIR /go/src/github.com/ravernkoh/jabba

COPY . .

RUN go build ./cmd/jabba

FROM alpine:3.9

ENV HOSTNAME jabba.ravernkoh.me
ENV PORT 80

WORKDIR /app

COPY --from=0 /go/src/github.com/ravernkoh/jabba/jabba .

CMD [ "./jabba" ]
