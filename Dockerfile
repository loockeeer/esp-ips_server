FROM golang AS build

COPY go.mod /src/espips_server/
WORKDIR /src/espips_server/

RUN go mod download

COPY . /src/espips_server/

RUN make build


FROM debian

WORKDIR /usr/src/app

COPY --from=build /src/espips_server/build/main.out .

CMD [ "./main.out" ]