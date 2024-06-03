package app

import (
	"flag"
	"fmt"
	"log"

	"github.com/go-redis/redis/v8"
	"golang.org/x/net/context"

	"testovoe/internal/models"
	"testovoe/internal/services"
)

const (
	factsPerMinute = 1000
)

// go run main.go -ip 192.168.31.6 -port 6379
// равносильно
// go run main.go
func Run() {
	var ip string
	var port string
	var testContext = context.Background()
	// Для запуска через консоль
	flag.StringVar(&ip, "ip", "192.168.31.6", "IP address of Redis server")
	flag.StringVar(&port, "port", "6379", "Port of Redis server")
	flag.Parse()

	rdb := redis.NewClient(&redis.Options{
		Addr:     fmt.Sprintf("%s:%s", ip, port),
		Password: "",
		DB:       0,
	})

	fact := make(chan models.Fact)
	if rdb.Ping(testContext).Err() != nil {
		log.Printf("Нет подключения к Redis")
	}
	//Запускаем нонстоп заполнение буфера фактами
	go services.RedisFiller(rdb, factsPerMinute)

	//Запускаем обработку буфера фактов в канал fact
	go services.FetchFromRedis(rdb, fact)
	var counter = 0
	//Запускаем отправку фактов в API
	go services.SendFactsToAPI(fact, &counter)

	//Ожидаем команду завершения программы
	select {}
}
