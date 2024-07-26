package drivers

import (
	"context"
	_ "embed"
	"fmt"
	"jobApps/helper"
	"os"

	"github.com/jackc/pgx/v4"
)

func DataBaseConnection() (*pgx.Conn, error) {
	err := helper.Configure(".env")
	if err != nil {
		fmt.Println("error is loading env file")
		return nil, err
	}

	host := os.Getenv("DB_HOST")
	port := os.Getenv("DB_PORT")
	password := os.Getenv("DB_PASSWORD")
	dbname := os.Getenv("DB_NAME")
	user := os.Getenv("DB_USER")

	connectionURI := fmt.Sprintf("postgresql://%s:%s@%s:%s/%s", user, password, host, port, dbname)

	conn, err := pgx.Connect(context.Background(), connectionURI)
	if err != nil {
		return nil, err
	}

	fmt.Println("Database Connected Successfully!!!...")

	return conn, nil
}

func TableCreation(conn *pgx.Conn) (err error) {
	createTablesQuery := `
	--users table
	CREATE TABLE IF NOT EXISTS users (
		UserID BIGSERIAL PRIMARY KEY,
		Username VARCHAR(255) NOT NULL,
		Email VARCHAR(255) NOT NULL,
		PhoneNumber VARCHAR(20),
		Password VARCHAR(255) NOT NULL,
		Role VARCHAR(255) NOT NULL,
		CreatedAt TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
		UpdatedAt TIMESTAMP DEFAULT CURRENT_TIMESTAMP
	);

	--profile table
	CREATE TABLE IF NOT EXISTS Profile (
		UserID  BIGSERIAL PRIMARY KEY,
		FullName VARCHAR(255) NOT NULL,
		Age Integer ,
		Gender VARCHAR(10)  NOT NULL,
		Address VARCHAR(255) NOT NULL
	);

	--career table
	CREATE TABLE IF NOT EXISTS Career (
		JobID  BIGSERIAL PRIMARY KEY,
		Company VARCHAR(255) NOT NULL,
		Position VARCHAR(255) NOT NULL,
		Jobtype  VARCHAR(255) NOT NULL,
		Description VARCHAR(255) NOT NULL,
		StartDate DATE,
		EndDate DATE
	);`

	_, err = conn.Exec(context.Background(), createTablesQuery)
	if err != nil {
		fmt.Println("Database creation Failed:", err)
		return
	}

	return nil
}
