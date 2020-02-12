FROM alpine:latest
MAINTAINER "NEO Dev <everestmx@gmail.com>"

# Install dependencies
RUN apk add --update git --no-cache make musl-dev go bash curl

COPY . /app

WORKDIR /app

# Configure Go
ENV GOROOT /usr/lib/go
ENV GOPATH /go
ENV PATH /go/bin:$PATH

RUN mkdir -p ${GOPATH}/src ${GOPATH}/bin

RUN go get -u github.com/Masterminds/glide/

WORKDIR $GOPATH

CMD ["make"]

ADD docker-entrypoint.sh /docker-entrypoint.sh

# Fix executable
RUN chmod +x /docker-entrypoint.sh

EXPOSE 8080

ENTRYPOINT ["/docker-entrypoint.sh"]
