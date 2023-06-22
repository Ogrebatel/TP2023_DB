FROM golang:1.19 AS build

COPY . /server/

WORKDIR /server/

RUN go build app/cmd/main.go

FROM ubuntu:20.04
COPY . .

RUN apt-get -y update && apt-get install -y tzdata
ENV TZ=Russia/Moscow
RUN ln -snf /usr/share/zoneinfo/$TZ /etc/localtime && echo $TZ > /etc/timezone

ENV PGVER 12
RUN apt-get -y update && apt-get install -y postgresql-$PGVER
USER postgres

RUN /etc/init.d/postgresql start &&\
    psql --command "CREATE USER db_pg WITH SUPERUSER PASSWORD 'db_postgres';" &&\
    createdb -O db_pg db_forum &&\
    psql -f db/db.sql -d db_forum &&\
    /etc/init.d/postgresql stop

RUN echo "host all  all    0.0.0.0/0  md5" >> /etc/postgresql/$PGVER/main/pg_hba.conf
RUN echo "listen_addresses='*'" >> /etc/postgresql/$PGVER/main/postgresql.conf

USER root
COPY --from=build /server/main .

EXPOSE 5000

CMD date && service postgresql start && ./main

