# Go REST API with sqlc

This repository demonstrates how to set up a Go REST API using the `sqlc` tool for generating type-safe SQL queries from PostgreSQL database schemas.

## Features

- RESTful API structure
- Type-safe SQL queries using `sqlc`
- PostgreSQL database integration
- CRUD operations

## Requirements

- Installing recent versions of sqlc requires Go 1.21+.
- PostgreSQL
- sqlc (installation instructions below)

## Getting Started

### Installation

1. **Clone the repository:**

    ```
    git clone git clone https://github.com/muthukumar89uk/go-restapi-with-sqlc.git
    ```
   Click here to directly [download it](https://github.com/muthukumar89uk/go-restapi-with-sqlc/zipball/master).

2. **Install sqlc:**

    ```sh
    go install github.com/sqlc-dev/sqlc/cmd/sqlc@latest
    ```

3. **Install dependencies:**

    ```sh
    go mod tidy
    ```
### Generate Type-safe SQL Queries

1. **Define SQL Queries:**

    Create a `sqlc.yaml` file in your project root with the following content:

    ```yaml
    version: "1"
    packages:
      - name: "db"
        path: "internal/db"
        queries: "./sql/"
        schema: "./schema/"
        engine: "postgresql"
    ```

2. **Create SQL Files:**

    Create a directory structure for your SQL files:

    ```
    .
    ├── sql
    │   ├── queries.sql
    |   └── schema.sql
    ```

    - **sql/queries.sql:**

        ```sql
        -- name: CreateCareer :one
       INSERT INTO career (Company,Position,Jobtype,Description,StartDate,EndDate)
       VALUES ($1, $2,$3,$4,$5,$6)
       RETURNING *;
       
       -- name: GetCareerByJobId :one
       SELECT * FROM career
       WHERE jobid = $1 LIMIT 1;
       
       -- name: GetAllCareerDetails :many
       SELECT * FROM career;
       
       -- name: UpdateCareerByJobId :one
       UPDATE career
       SET company=$1,position=$2,jobtype=$3,description=$4
       WHERE jobid = $5
       RETURNING *;
       
       -- name: DeleteCareerByJobId :one
       DELETE
       FROM career
       WHERE jobid = $1
       RETURNING *;
        ```

    - **sql/schema.sql:**

        ```sql
       CREATE TABLE IF NOT EXISTS Career (
          JobID  BIGSERIAL PRIMARY KEY,
          Company VARCHAR(255) NOT NULL,
          Position VARCHAR(255) NOT NULL,
          Jobtype  VARCHAR(255) NOT NULL,
          Description VARCHAR(255) NOT NULL,
          StartDate DATE  NOT NULL ,
          EndDate DATE NOT NULL
        );
        ```

3. **Generate Code:**

    Run the `sqlc` command to generate Go code from the SQL queries:

    ```sh
    sqlc generate
    ```

### Run the Application

1. **Run the application:**

    ```sh
    go run .
    ```

2. **Access the API:**

    The API will be available at `http://localhost:8080`.

    **Example Endpoints:**

    - `GET /get-all-career-details` - Retrieve all career-details 
    - `GET /getcareerdetail/:id` - Retrieve an career details by ID
    - `POST /create-career` - Create a new career
    - `PUT /updatecareer/:id` - Update an existing career
    - `DELETE /deletecareer/:id` - Delete a career
