FROM alpine:latest
MAINTAINER "NEO Dev <everestmx@gmail.com>"

# Install dependencies
RUN apk add --update git --no-cache

COPY . /app

WORKDIR /app

#RUN go get -v -d && CGO_ENABLED=0 GOOS=linux go build -a -installsuffix cgo -o app-service

#
# This results in a single layer image
#
# Install dependencies
RUN apk add --update ca-certificates bash curl --no-cache

ADD docker-entrypoint.sh /docker-entrypoint.sh

# Copy squid-auth from `golang`
#COPY --from=build /app/app-service /usr/bin/app-service

# Fix executable
RUN chmod +x /docker-entrypoint.sh

EXPOSE 8080

ENTRYPOINT ["/docker-entrypoint.sh"]
