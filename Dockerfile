FROM --platform=$BUILDPLATFORM golang:1.23.4-alpine AS builder

WORKDIR /app
COPY main.go .
COPY go.mod .

ARG TARGETOS
ARG TARGETARCH

RUN <<EOF
go mod tidy 
GOOS=${TARGETOS} GOARCH=${TARGETARCH} go build
EOF

FROM scratch

WORKDIR /app

COPY --from=builder /app/amphipod .
COPY functions ./functions 
COPY tools ./tools
#CMD ["./amphipod"]
