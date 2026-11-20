FROM golang:1.27.1-alpine AS build

WORKDIR /src
COPY go.mod ./
COPY . .
RUN CGO_ENABLED=0 go build -trimpath -ldflags="-s -w" -o /server .

FROM alpine:3.22

WORKDIR /app
RUN apk add --no-cache bash
COPY --from=build /server ./server
COPY --from=build /src/templates ./templates
COPY --from=build /src/static ./static
COPY --from=build /src/standard.txt /src/shadow.txt /src/thinkertoy.txt /src/Binary.flf ./

USER 65532:65532
EXPOSE 8080
LABEL org.opencontainers.image.title="ASCII Art Web" \
      org.opencontainers.image.description="A Go web application that renders text as ASCII art" \
      org.opencontainers.image.authors="Ghazal and Mohammad"
ENTRYPOINT ["./server"]