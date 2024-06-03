package services

import (
	"bytes"
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"net/url"
	"os"
	"time"

	"github.com/go-redis/redis/v8"
	"golang.org/x/net/context"

	"testovoe/internal/models"
)

var (
	// Задаем контекст для бд
	ctx = context.Background()
	// Забираем из окружения переменные
	method      = os.Getenv("METHOD")
	urlStr      = os.Getenv("URL_STR")
	bearerToken = os.Getenv("BEARER_TOKEN")
)

// Инициализируем переменные стандартными значениями, если они не были переданы в окружение
func init() {
	if method == "" {
		method = "POST"
	}
	if urlStr == "" {
		urlStr = "https://development.kpi-drive.ru/_api/facts/save_fact"
	}
	if bearerToken == "" {
		bearerToken = "48ab34464a5573519725deb5865cc74c"
	}
}

// RedisFiller заполняет Redis данными фактов
func RedisFiller(rdb *redis.Client, factsPerMinute int) {
	for {
		for i := 0; i < factsPerMinute; i++ {
			// Задержка для постепенной отправки запросов
			time.Sleep(time.Duration(time.Minute.Milliseconds() / int64(factsPerMinute)))
			// Заполняем структуру факта любыми данными(просто для примера)
			fact := models.Fact{
				PeriodStart:         "2024-05-01",
				PeriodEnd:           "2024-05-31",
				PeriodKey:           "month",
				IndicatorToMoId:     227373,
				IndicatorToMoFactId: 0,
				Value:               1,
				FactTime:            "2024-05-31",
				IsPlan:              0,
				AuthUserId:          40,
				Comment:             "buffer Last_name",
			}

			// Конвертируем структуру в байты
			data, err := json.Marshal(fact)
			if err != nil {
				log.Printf("Ошибка при переводе факта в байты: %v", err)
			}

			// Отправляем запрос в Redis как очередь
			err = rdb.RPush(ctx, "facts_queue", data).Err()
			if err != nil {
				log.Printf("Ошибка при добавлении факта в Redis: %v", err)
			}
		}
	}
}

// FetchFromRedis извлекает данные из Redis и отправляет их в канал
func FetchFromRedis(rdb *redis.Client, facts chan models.Fact) {
	for {
		// Извлечение данных из Redis и тут же удаление записи из базы
		data, err := rdb.LPop(ctx, "facts_queue").Bytes()
		if err != nil {
			// Опциональное ожидание данных из Redis
			if err == redis.Nil {
				// Если данных нет, подождать и попробовать снова
				time.Sleep(1 * time.Second)
				continue
			}
			log.Printf("Ошибка при извлечении данных из Redis: %v", err)
		}

		// Демаршалинг данных в структуру Fact
		var fact models.Fact
		if err := json.Unmarshal(data, &fact); err != nil {
			log.Printf("Ошибка при демаршалинге факта: %v", err)
		}

		// Отправка факта в канал, блокирующаяся до тех пор, пока канал не освободится
		facts <- fact
	}
}

// SendFactsToAPI получает данные из канала и отправляет их на сервер
func SendFactsToAPI(facts chan models.Fact, counter *int) {
	for {
		// Ожидание данных из канала (блокируется, если канал пуст)
		fact := <-facts
		if *counter < 10 {
			*counter++
		} else {
			break
		}
		// Отправка данных на сервер
		if err := sendFact(fact); err != nil {
			log.Printf("Ошибка при отправке факта в API: %v", err)
			// Обработка ошибок, например, повторная отправка факта или логирование
		}
	}
}

// sendFact отправляет данные факта на API сервер
func sendFact(fact models.Fact) error {
	// Собираем тело запроса
	data := url.Values{}
	data.Set("period_start", fact.PeriodStart)
	data.Set("period_end", fact.PeriodEnd)
	data.Set("period_key", fact.PeriodKey)
	data.Set("indicator_to_mo_id", fmt.Sprintf("%d", fact.IndicatorToMoId))
	data.Set("indicator_to_mo_fact_id", fmt.Sprintf("%d", fact.IndicatorToMoFactId))
	data.Set("value", fmt.Sprintf("%d", fact.Value))
	data.Set("fact_time", fact.FactTime)
	data.Set("is_plan", fmt.Sprintf("%d", fact.IsPlan))
	data.Set("auth_user_id", fmt.Sprintf("%d", fact.AuthUserId))
	data.Set("comment", fact.Comment)
	// Собираем запрос
	req, err := http.NewRequest(method, urlStr, bytes.NewBufferString(data.Encode()))
	if err != nil {
		return fmt.Errorf("ошибка при создании запроса: %w", err)
	}
	// Добавляем заголовки
	req.Header.Add("Authorization", fmt.Sprintf("Bearer %s", bearerToken))
	req.Header.Add("Content-Type", "application/x-www-form-urlencoded")
	// Устанавливаем соединение с сервером ( по-хорошему его нужно держать пока не закончится буфер)
	client := &http.Client{}
	resp, err := client.Do(req)
	if err != nil {
		return fmt.Errorf("ошибка при выполнении запроса: %w", err)
	} else {
		fmt.Printf("Сервер вернул код ответа: %d", resp.StatusCode)
	}
	// После обработки ошибок указываем на закрытие тела запроса по окончании использования функции
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return fmt.Errorf("неудачный статус ответа: %d", resp.StatusCode)
	}

	return nil
}
