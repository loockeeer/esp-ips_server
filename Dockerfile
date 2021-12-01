FROM golang AS build

WORKDIR /usr/src/app

COPY . .

RUN make build

FROM debian

WORKDIR /usr/src/app

COPY --from=build /usr/src/app/build/main.out .

CMD [ "./main.out" ]