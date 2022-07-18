FROM golang:1.18-alpine AS builder

WORKDIR /app
COPY . /app
RUN apk add make bash
RUN make build


FROM golang:1.18-alpine
COPY --from=builder /app/simpledb /app/simpledb

ENV SIMPLEDB_PATH=/database \
    SIMPLEDB_HOST=0.0.0.0 \
    SIMPLEDB_PORT=5432
EXPOSE $SIMPLEDB_PORT

ENTRYPOINT /app/simpledb
