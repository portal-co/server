FROM golang:1.16-alpine AS dev

WORKDIR /src
COPY . /src
RUN go build -o /boot portal.pc/server/main_

EXPOSE  8000

CMD [ "/boot" ]