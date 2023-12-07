package main

import (
	"Service/database"
	"Service/handler"
	"flag"
	"log"
)

func main() {
	// Флаг для выбора между Redis и PostgreSQL
	var useRedis bool
	flag.BoolVar(&useRedis, "use-redis", false, "Use Redis for URL storage")
	flag.Parse()

	var storage database.URLStorage
	//var err error

	if useRedis {
		// Использование Redis в качестве хранилища
		redisClient, err := database.ConnectToRedis()
		if err != nil {
			log.Fatalf("Failed to connect to Redis: %s", err)
			return
		}
		storage = &database.RedisStorage{Client: redisClient}
	} else {
		// Использование PostgreSQL в качестве хранилища
		db, err := database.ConnectToDB()
		if err != nil {
			log.Fatalf("Failed to connect to PostgreSQL: %s", err)
			return
		}
		defer db.Close()
		storage = &database.PostgresStorage{DB: db}
	}

	handler.StartServer(storage)
}
