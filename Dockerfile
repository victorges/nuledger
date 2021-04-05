FROM golang:1.15 AS build

WORKDIR /app
COPY . ./

RUN make test
RUN make build

FROM alpine AS run

WORKDIR /app
COPY --from=build /app/build/authorizer /app

ENTRYPOINT ./authorizer