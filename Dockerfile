FROM golang:1.25-alpine AS build

WORKDIR /build

COPY go.mod go.sum ./
RUN go mod download

COPY api ./api
COPY *.go ./

RUN CGO_ENABLED=0 go build -ldflags "-w -s" -o /todotxt

FROM gcr.io/distroless/static-debian13:nonroot AS deploy

# hadolint ignore=DL3045
COPY frontend/dist/spa /frontend/dist/spa/
COPY --from=build /todotxt /

EXPOSE 3000
CMD ["/todotxt"]
