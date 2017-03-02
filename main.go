package main

import (
	"flag"
	"os"
	"fmt"
	_ "github.com/go-sql-driver/mysql"
	"database/sql"
	"strconv"
)

func defaultString(varName, defaultValue string) string {
	if "" != os.Getenv(varName) {
		return os.Getenv(varName)
	}

	return defaultValue
}

func defaultInt(varName string, defaultValue int) (int, error) {
	if "" != os.Getenv(varName) {
		value, err := strconv.Atoi(os.Getenv(varName))

		if nil != err {
			return -1, err
		}

		return value, nil
	}

	return defaultValue, nil
}

func main() {
	host := "localhost"
	port := 3306

	newVal, err := defaultInt("MYSQL_PORT", 3306)

	if nil == err {
		port = newVal
	}

	user := "root"
	password := "password"
	database := "information_schema"
	expectedCount := 1

	flag.StringVar(&host, "host", defaultString("MYSQL_HOST", host), "MySQL Host")
	flag.IntVar(&port, "port", port, "MySQL Port")
	flag.StringVar(&user, "user", defaultString("MYSQL_USER", user), "User")
	flag.StringVar(&password, "password", defaultString("MYSQL_PASSWORD", password), "Password")
	flag.StringVar(&database, "database", defaultString("MYSQL_DATABASE", database), "Database Name")

	flag.IntVar(&expectedCount, "expected", expectedCount, "Expected Count")

	flag.Parse()

	dataSourceName := fmt.Sprintf("%s:%s@/tcp(%s:%d)/%s?charset=utf8", user, password, host, port, database)

	db, err := sql.Open("mysql", dataSourceName)

	if nil != err {
		panic(err)
	}

	defer func() {
		db.Close()
	}()

	remaining := os.Args[flag.NArg():]

	for _, query := range remaining {
		query, err := db.Query(query)

		if nil != err {
			panic(err)
		}

		count := 0

		for query.Next() {
			count++
		}

		if count != expectedCount {
			panic(fmt.Errorf("Unexpected row count for query: '%s'", query))
		}
	}
}
