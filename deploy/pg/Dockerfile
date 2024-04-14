FROM postgres:latest

COPY ./deploy/pg/init-db.sh /docker-entrypoint-initdb.d/

EXPOSE 5432

CMD ["postgres"]