FROM golang:1.20-alpine AS dev


WORKDIR /src
COPY . /src
RUN go build -o /boot github.com/portal-co/server/void

EXPOSE  8000

CMD [ "/boot" ]