FROM golang:1.19.3
RUN go install github.com/volatiletech/sqlboiler/v4@latest
RUN go install github.com/volatiletech/sqlboiler/v4/drivers/sqlboiler-psql@latest
WORKDIR /src/sqlboiler
CMD ["sqlboiler", "psql"]
