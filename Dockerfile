#build stage
FROM golang:1.16.0-buster as build-env
ARG GH_TOKEN
RUN git config --global url."https://${GH_TOKEN}:x-oauth-basic@github.com/ProjectAthenaa".insteadOf "https://github.com/ProjectAthenaa"
RUN git config --global url."https://${GH_TOKEN}:x-oauth-basic@github.com/A-Solutionss".insteadOf "https://github.com/A-Solutionss"

RUN mkdir /app
ADD . /app
WORKDIR /app
RUN --mount=type=cache,target=/root/.cache/go-build
RUN go mod download
RUN go build -o shape_gen


# final stage
FROM debian:buster-slim
WORKDIR /app
COPY --from=build-env /app/shape_gen /app/

ENV REDIS_URL="rediss://default:n6luoc78ac44pgs0@test-redis-do-user-9223163-0.b.db.ondigitalocean.com:25061"
ENV ENVIRONMENT="Development"

EXPOSE 3000 3000

ENTRYPOINT ./shape_gen