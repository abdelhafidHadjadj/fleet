FROM golang AS build-stage


WORKDIR /app

COPY go.mod go.sum ./

RUN go mod download

COPY . .

RUN CGO_ENABLED=0 GOOS=linux go build -o /fleet

FROM build-stage AS run-test-stage

RUN go test -v ./...

FROM gcr.io/distroless/cc-debian12	AS build-release-stage

WORKDIR /

COPY --from=build-stage /fleet /fleet

COPY --from=build-stage /app/.env .

EXPOSE 8080

ENTRYPOINT [ "/fleet" ]