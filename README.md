# Тестовое задание 

Создание буфера прилетающих фактов

## Установка

1. Установите Go
2. Установите Redis: https://redis.io/docs/latest/operate/oss_and_stack/install/install-redis/install-redis-on-windows/
3. Клонируйте проект: git clone https://github.com/Svatoslav191/TestovoeBuffer.git`
3. Перейдите в директорию проекта: `cd ваша_директория`
4. Установите зависимости: `go mod download`

## Использование
1. На машине с установленным Redis:
	1) Перейдите в файл конфигурации Redis (vim /etc/redis/redis.conf)
	2) Укажите адрес Redis, например "bind 192.168.31.6"
	3) Запустите сервер "service redis-server start"
	4) Зайдите на клиента: redis-cli -h 192.168.31.6 -p 6379
	5) Убедитесь, что сервер отвечает на ping PONG`ом
2. Запустите приложение:
	Из-под Go
		1) go run cmd/main.go -ip 192.168.31.6 -port 6379
		2) go run cmd/main.go (измените стандартные значения на 24 и 25 строчках internal/app/app.go на ваши адрес и порт
	exe:
		1) main.exe -ip (YOUR_IP) -port (YOUR_PORT)
		2) Скомпилируйте программу с измененными адресом и портом, запустите даблкликом

## Дополнительные идеи
1. Стоит подумать над Redis Sentinel - Мониторит ведущие и подчиненные узлы, отправляет уведомления о происшествиях,
   восстановление после отказа запуском процесса восстановления после кворума узлов о недоступности главной ноды редиса
2. Обращаться к нему для определения доступности бд
3. Рекомендуется составлять по 3 узла с кворумом из 2х нод
