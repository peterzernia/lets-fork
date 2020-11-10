# [Let's Fork](https://letsfork.app)
The server for a Let's Fork built with Go, Redis and Docker.

## Demo
[![Let's Fork](https://letsfork.app/thumbnail.png)](https://vimeo.com/477511968 "Let's Fork")

## Development
Requirements: [Docker Compose](https://docs.docker.com/compose/)

Copy `.env.dist` to `.env`.

Get an api key from [Yelp](https://www.yelp.com/developers) and set `API_TOKEN=`

Run `make build` builds Docker image.

Run `make up` to starts the server. 