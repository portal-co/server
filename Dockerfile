FROM golang:1.16-alpine AS dev

WORKDIR /src
COPY . /src
RUN go build -o /boot

CMD [ "/boot" ]