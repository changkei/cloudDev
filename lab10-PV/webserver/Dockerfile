FROM golang:1.15.7-alpine AS build

RUN apk add git
RUN go get -u go.mongodb.org/mongo-driver/mongo

WORKDIR /src/
COPY webserver.go /src/
RUN CGO_ENABLED=0 go build -o /bin/webserver

FROM scratch
COPY --from=build /bin/webserver /bin/webserver
ENTRYPOINT ["/bin/webserver"]
