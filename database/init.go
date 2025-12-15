package database

import (
	"fmt"

	"github.com/dtekltd/common/system"
)

func init() {
	fmt.Println("> init database")

	// redis
	if rc := system.Env("REDIS_CONN"); rc != "" {
		var err error
		RedisClient, err = NewRedisClient(rc)
		if err != nil {
			system.Logger.Panicf("Failed to init redis connection \"%s\" - %v", rc, err)
		}
	}

	// gorm db
	switch system.Env("DB_TYPE") {
	case "mysql":
		InitMySqlDB(system.Env("DB_RUNTIME_TYPE") == "mysql")
	case "postgres":
		InitPostgresDB(system.Env("DB_RUNTIME_TYPE") == "postgres")
	}
}
