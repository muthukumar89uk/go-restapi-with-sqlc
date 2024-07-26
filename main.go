package main

import (
	"context"
	_ "embed"
	"fmt"
	"jobApps/drivers"
	router "jobApps/routers"
)

//go:embed sql/schema.sql
var schemaSql string

func main() {
	conn, err := drivers.DataBaseConnection()
	if err != nil {
		fmt.Println(err)
		return
	}

	_, err = conn.Exec(context.Background(), schemaSql)
	if err != nil {
		fmt.Println("table creation failed:", err)
		return
	}

	router.Router(conn)
}
