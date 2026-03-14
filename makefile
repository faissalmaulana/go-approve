include .env

API_PATH = api
MIGRATE_PATH = $(API_PATH)/migrations
DB_URL = mysql://root:$(DB_MIGRATE_AUTH)@/$(DATABASE_NAME)

.PHONY: migrate migrate-up migrate-down migrate-back run-app

migrate:
	cd $(API_PATH) && migrate create -ext sql -dir migrations -seq $(name)

migrate-up:
	migrate -path=$(MIGRATE_PATH) -database=$(DB_URL) up

migrate-down:
	migrate -path=$(MIGRATE_PATH) -database=$(DB_URL) down

migrate-back:
	migrate -path=$(MIGRATE_PATH) -database=$(DB_URL) force $(version)

run-server:
	cd api && air

run-client:
	cd app && pnpm run dev

run-app:
	make -j 2 run-server run-client
