FROM golang:1.22 AS build
WORKDIR /src
COPY . .
RUN CGO_ENABLED=0 go build -trimpath -ldflags='-s -w' -o /out/timeseriesd ./cmd/timeseriesd
FROM gcr.io/distroless/static-debian12
COPY --from=build /out/timeseriesd /timeseriesd
EXPOSE 8080
ENTRYPOINT ["/timeseriesd"]
