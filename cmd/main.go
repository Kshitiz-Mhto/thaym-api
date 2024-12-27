package main

import (
	"fmt"
	"log"

	"ecom-api/internal/adapters/framework/right/db"
	"ecom-api/internal/application/api"
	"ecom-api/pkg/configs"

	"github.com/go-sql-driver/mysql"
)

func main() {
	cfg := mysql.Config{
		User:                 configs.Envs.DBUser,
		Passwd:               configs.Envs.DBPassword,
		Addr:                 configs.Envs.DBAddress,
		DBName:               configs.Envs.DBName,
		Net:                  "tcp",
		AllowNativePasswords: true,
		ParseTime:            true,
	}
	dbAdapter, err := db.NewDbAdapter(cfg)
	if err != nil {
		log.Fatalf("%v", err)
	}
	dbInstance := dbAdapter.GetDBInstance()
	server := api.NewAPIServer(fmt.Sprintf(":%s", configs.Envs.Port), dbInstance)
	if err := server.Run(); err != nil {
		log.Fatal(err)
	}
}
