FROM golang:1.16-alpine AS dev

WORKDIR /src
COPY boot.go ./
COPY go.mod ./
RUN go build -o /boot

EXPOSE  8000

CMD [ "/boot" ]