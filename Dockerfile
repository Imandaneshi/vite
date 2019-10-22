# 1. using golang alpine image to build server side codes
FROM golang:alpine as goBuilder

ADD . /code

WORKDIR /code

RUN go mod download

WORKDIR /code/cmd/vite

RUN go install


# 2. using node alpine image to build our front end project
FROM node:alpine as nodeBuilder

ADD ./client /code

WORKDIR /code

RUN yarn install

RUN yarn build


# 3. final image with no dependency, just binaries and static files
FROM alpine

# https://stackoverflow.com/questions/34729748/installed-go-binary-not-found-in-path-on-alpine-linux-docker/35613430#35613430
RUN mkdir /lib64 && ln -s /lib/libc.musl-x86_64.so.1 /lib64/ld-linux-x86-64.so.2

WORKDIR /app

COPY --from=goBuilder /go/bin/vite /usr/local/bin/vite

COPY --from=nodeBuilder /code/build /var/www/vite

EXPOSE 8062

CMD ["vite", "server", "--static-path", "/var/www/vite"]