FROM golang:1.24

WORKDIR /usr/src/app

# # pre-copy/cache go.mod for pre-downloading dependencies and only redownloading them in subsequent builds if they change
# COPY ./backend/go.mod ./backend/go.sum ./backend
#
# WORKDIR /usr/src/app/backend
#
# RUN go mod download
#
# WORKDIR /usr/src/app
COPY . .
# RUN cd ./backend && go build -v -o /usr/local/bin/app ./...
#
# CMD ["app"]