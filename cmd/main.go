package main

import (
	"arceus/pkg/carbon/pkg/config"

	_ "github.com/go-sql-driver/mysql"

	"arceus/internal/server"
)

func main() {
	flags := config.ParseFlags()
	server.Run(flags)
}
