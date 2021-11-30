FROM golang AS build

WORKDIR /usr/src/app

COPY . .

RUN go build .

FROM alpine

WORKDIR /usr/src/app

COPY ./config.yaml .

COPY --from=build /usr/src/app/build/main .

CMD [ "./main" ]