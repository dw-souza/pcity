# API Go

Servidor HTTP do Pcity. Contrato: [`docs/api.openapi.yaml`](../docs/api.openapi.yaml).

## Requisitos

- Go 1.25+
- Postgres com PostGIS — `docker compose up` na raiz do monorepo
- Supabase Auth (JWT) para endpoints autenticados em produção; ver `.env.example` para dev local

## Setup

```bash
# Na raiz do monorepo
docker compose up -d

cd api
cp .env.example .env
go run ./cmd/server
```

Servidor em `http://localhost:8080/v1`.

## Endpoints principais

| Método | Path |
|---|---|
| GET | `/v1/health` |
| GET | `/v1/cities/{slug}/places` |
| GET | `/v1/places/{id}` |
| POST | `/v1/places/{id}/amenities` |
| POST | `/v1/admin/cities/{slug}/index` |

## Testes

```bash
go test ./...
```

## Estrutura

```
cmd/server/          # entrypoint
internal/
  handlers/          # HTTP
  services/          # regras de negócio
  repository/        # pgx + PostGIS
  middleware/        # auth, admin, logging
  googleplaces/      # client Google Places API
```
