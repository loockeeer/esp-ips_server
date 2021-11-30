FROM golang AS build

WORKDIR /usr/src/app

COPY . .

RUN make build

FROM alpine

WORKDIR /usr/src/app

COPY ./config.yaml .

COPY --from=build /usr/src/app/build/main .

CMD [ "./main" ]