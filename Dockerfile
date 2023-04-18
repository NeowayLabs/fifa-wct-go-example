# --- Base ----
FROM golang:1.19 AS base
WORKDIR $GOPATH/src/github.com/NeowayLabs/fifa-wct-go-example

# ---- Dependencies ----
FROM base AS dependencies
ENV GO111MODULE=on
COPY go.mod .
COPY go.sum .
RUN go mod download

# ---- Test Integration ----
FROM dependencies AS test-integration
COPY . .
ARG MONGO_URL
RUN MONGO_URL=$MONGO_URL go test -race -tags integration -coverprofile=coverage.out -covermode=atomic -timeout 300s -v ./...
RUN grep -v "_mock" coverage.out >> filtered_coverage.out
RUN go tool cover -func filtered_coverage.out

# ---- Test ----
FROM dependencies AS test-unit
COPY . .
RUN go test -race -tags unit -coverprofile=coverage.out -covermode=atomic -timeout 300s -v ./...
RUN grep -v "_mock" coverage.out >> filtered_coverage.out
RUN go tool cover -func filtered_coverage.out

# ---- Lint ----
FROM dependencies AS lint
COPY . .
RUN curl -sSfL https://raw.githubusercontent.com/golangci/golangci-lint/master/install.sh | sh -s -- -b $(go env GOPATH)/bin v1.49.0
RUN golangci-lint run

# ---- Audit ----
FROM dependencies AS audit
COPY . .
RUN curl -sSfL https://github.com/sonatype-nexus-community/nancy/releases/download/v0.1.17/nancy-linux.amd64-v0.1.17 -o $(go env GOPATH)/bin/nancy && \
    chmod +x $(go env GOPATH)/bin/nancy
RUN nancy go.sum

# ---- Build ----
FROM dependencies AS build
COPY . .
ARG VERSION
ARG BUILD
ARG DATE
RUN CGO_ENABLED=0 GOOS=linux GOARCH=amd64 go build -v -a -installsuffix cgo -o /go/bin/fifa-wct-go-example -ldflags "-X main.version=${VERSION} -X main.build=${BUILD} -X main.date=${DATE}" ./cmd/fifa-wct-go-example

# --- Release ----
FROM debian:stable-slim AS image
COPY --from=build /go/bin/fifa-wct-go-example /fifa-wct-go-example
USER nobody
ENTRYPOINT ["/fifa-wct-go-example"]