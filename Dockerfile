FROM golang:1.18-alpine as build

WORKDIR /go/src/
COPY . .

RUN CGO_ENABLED=0 go build \
    -mod=vendor \
    -o /clank \
    clank/main.go

FROM alpine
RUN mkdir clank
COPY --from=build /clank clank/clank
ENV CLANK_LOG_LEVEL=info
ENTRYPOINT ["/clank/clank"]
