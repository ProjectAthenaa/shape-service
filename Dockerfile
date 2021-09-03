#build stage
FROM golang:1.17.0-buster as build-env
ARG GH_TOKEN
RUN git config --global url."https://${GH_TOKEN}:x-oauth-basic@github.com/ProjectAthenaa".insteadOf "https://github.com/ProjectAthenaa"
RUN mkdir /app
ADD . /app
WORKDIR /app
RUN --mount=type=cache,target=/root/.cache/go-build
RUN go build -ldflags "-s -w" -o shape_gen


# final stage
FROM debian:buster-slim
WORKDIR /app
COPY --from=build-env /app/shape_gen /app/

RUN apt-get update \
 && apt-get install -y --no-install-recommends ca-certificates

RUN update-ca-certificates


EXPOSE 3000 3000

ENTRYPOINT ./shape_gen