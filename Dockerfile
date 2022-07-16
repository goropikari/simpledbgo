# FROM golang:1.18-buster AS builder

# WORKDIR /app
# COPY . /app
# RUN apt-get update && apt-get upgrade -y
# RUN make build


# FROM golang:1.18-buster
# COPY --from=builder /app/simpledb /app/simpledb

# ENv SIMPLEDB_PATH=/database \
# ENV DBMS_HOST=0.0.0.0 \
#     DBMS_PORT=5432
# EXPOSE $DBMS_PORT

# ENTRYPOINT /app/simpledb


FROM golang:1.18-alpine AS builder

WORKDIR /app
COPY . /app
RUN apk add make bash
RUN make build


FROM golang:1.18-alpine
COPY --from=builder /app/simpledb /app/simpledb

ENV SIMPLEDB_PATH=/database \
    DBMS_HOST=0.0.0.0 \
    DBMS_PORT=5432
EXPOSE $DBMS_PORT

ENTRYPOINT /app/simpledb
