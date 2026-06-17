.PHONY: up down reset logs psql migrate

up:
	docker compose up -d

down:
	docker compose down

reset:
	docker compose down -v

logs:
	docker compose logs -f db

psql:
	docker compose exec db psql -U postgres -d pcity

# Reaplica migrations manualmente (se o volume já existia ou init falhou)
migrate:
	docker compose exec -T db psql -U postgres -d pcity < supabase/migrations/000000_auth_stub.sql
	docker compose exec -T db psql -U postgres -d pcity < supabase/migrations/000001_initial_schema.sql
	docker compose exec -T db psql -U postgres -d pcity < supabase/seed/dev.sql
