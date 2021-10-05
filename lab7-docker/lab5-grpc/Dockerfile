FROM golang:1.15.7-alpine AS build

RUN apk add git
RUN go get -u google.golang.org/grpc

COPY server.go /go/src
WORKDIR /go/src/movieapi
COPY movieapi /go/src/movieapi

RUN CGO_ENABLED=0 go build -o /bin/moviegrpc /go/src/server.go

FROM scratch
COPY --from=build /bin/moviegrpc /bin/moviegrpc
ENTRYPOINT ["/bin/moviegrpc"]