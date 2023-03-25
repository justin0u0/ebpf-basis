# FROM gcr.io/distroless/base-debian10 AS server
FROM ubuntu:22.04 AS base

RUN apt-get update && apt-get install -y iproute2 netcat

COPY bin/cmd /

FROM base AS server

EXPOSE 8080 8081

CMD ["/cmd", "server"]

FROM rabbitmq:3.11.4-management AS rabbit

RUN apt-get update && apt-get install -y iproute2

COPY bin/cmd /

FROM postgres:15.2 AS postgres-ebpf

RUN apt-get update && apt-get install -y iproute2

COPY bin/cmd /

FROM bitnami/pgbouncer:1.18.0 AS pgbouncer-ebpf

USER root
RUN apt-get update && apt-get install -y iproute2
USER 1001

COPY bin/cmd /
