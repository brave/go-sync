FROM public.ecr.aws/docker/library/postgres:16

RUN apt update && apt install -y git make gcc postgresql-server-dev-16

RUN git clone https://github.com/pgpartman/pg_partman
RUN cd pg_partman && make NO_BGW=1 install

RUN git clone https://github.com/citusdata/pg_cron
RUN cd pg_cron && make && make install
