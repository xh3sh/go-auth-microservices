# Makefile для последовательной сборки и запуска микросервисов

.PHONY: default build-all up down restart logs ps

# По умолчанию запускаем последовательную сборку всех сервисов
default: build-all

# Сборка сервисов по одному (no-parallel)
build-all:
	@echo "--- Building Auth Service ---"
	docker compose build auth-service
	@echo "--- Building User Service ---"
	docker compose build user-service
	@echo "--- Building Log Consumer ---"
	docker compose build log-consumer
	@echo "--- Building Frontend Service ---"
	docker compose build frontend-service
	@echo "--- Building API Gateway ---"
	docker compose build api-gateway

# Полный цикл: сборка и запуск в фоне
up: build-all
	@echo "--- Starting all services ---"
	docker compose up -d

# Остановка и удаление контейнеров
down:
	@echo "--- Stopping services ---"
	docker compose down

# Перезапуск всего стека
restart: down up

# Просмотр логов в реальном времени
logs:
	docker compose logs -f

# Статус контейнеров
ps:
	docker compose ps
