# AWS Aurora DB Data Generator

Since inserting large amounts of data using the Insert statement takes too long, we will use the AWS Aurora native feature `LOAD DATA FROM S3` syntax to generate a CSV file to store in AWS S3.

```sql
LOAD DATA FROM S3 's3://<s3-bucket-name>/path//table_name.csv'
INTO TABLE customer 
FIELDS TERMINATED BY ',' 
LINES TERMINATED BY '\n';
```

## Requirements

- RDB (MySQL or PostgreSQL)

## How to use

```sh

go run ./cmd/test_data_generator/ --app <app_name> --count 10000
```
