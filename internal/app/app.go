package app

import (
	"github.com/go-redis/redis/v8"

	"testovoe/internal/models"
	"testovoe/internal/services"
)

const (
	factsPerMinute = 1000
)

func Run() {
	rdb := redis.NewClient(&redis.Options{
		Addr:     "192.168.31.6:6379",
		Password: "",
		DB:       0,
	})
	fact := make(chan models.Fact)
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
