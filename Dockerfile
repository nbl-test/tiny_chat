## Build
FROM golang:1.20-buster AS build
WORKDIR /app
ENV CGO_ENABLED=0
ENV GOPROXY=https://goproxy.cn,direct
COPY . /app
RUN go mod tidy
RUN go build --ldflags '-w -s -extldflags "-static"' -asmflags -trimpath -o /tiny_chat 

## Deploy
FROM scratch
EXPOSE 8080
COPY --from=build /tiny_chat /tiny_chat
ENTRYPOINT ["/tiny_chat"]