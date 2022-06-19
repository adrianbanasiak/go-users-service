FROM golang:1.18-alpine AS build
RUN apk add git

ENV CGO_ENABLED=0
WORKDIR /src

COPY . /src/

RUN go get .
RUN go build -o bin/server . && chmod +x bin/server

FROM scratch AS bin
COPY --from=build /src/bin/server /server

CMD ["/server"]