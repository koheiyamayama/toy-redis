FROM golang:1.21.4-bullseye as build
WORKDIR /app
COPY . /app
RUN go build -o /app/cmd/toy-redis /app

FROM debian:bullseye-slim as app
COPY --from=build /app/cmd/toy-redis ./
CMD [ "./toy-redis" ]
