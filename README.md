# Pcity

Agregador de bares e restaurantes — descubra lugares perto de você com informações que importam (espaço kids, pet friendly, etc.).

**Cidade piloto:** Franca, SP

## Stack

| Camada | Tecnologia |
|---|---|
| Mobile | Expo + TypeScript |
| API | Go |
| DB + Auth | Supabase (Postgres + PostGIS) |
| Dados iniciais | Google Places API |

## Estrutura

```
pcity/
├── api/                  # API Go
├── mobile/               # App Expo (React Native)
├── docs/                 # Specs, OpenAPI, fluxos UX
├── supabase/
│   └── migrations/       # Migrations SQL
└── docker-compose.yml    # Postgres local (dev)
```

## Setup local

### 1. Banco (Postgres + PostGIS)

```bash
docker compose up -d
# migrations rodam automaticamente na primeira subida
# Postgres exposto em localhost:5433 (evita conflito com Postgres local na 5432)
```

Comandos úteis: `make up`, `make down`, `make reset` (apaga volume e reinicializa), `make psql`.

**Primeira subida** aplica em ordem:
1. `000000_auth_stub.sql` — schema `auth` + `auth.uid()` stub
2. `000001_initial_schema.sql` — schema público
3. `supabase/seed/dev.sql` — perfil do usuário dev

### 2. API

```bash
cd api
cp .env.example .env
go run ./cmd/server
```

Servidor em `http://localhost:8080/v1`.

### 3. Seed de lugares (opcional)

```bash
curl -X POST http://localhost:8080/v1/admin/cities/franca-sp/index \
  -H "X-Admin-Key: dev-admin-key"
```

Requer `GOOGLE_PLACES_API_KEY` no `.env`. Tutorial completo: [docs/google-places-setup.md](docs/google-places-setup.md).

### 4. Mobile

```bash
cd mobile
cp .env.example .env
npm install
npm start
```

Ver [mobile/README.md](mobile/README.md) para URL da API no emulador/dispositivo.

## Status

MVP em implementação — spec completa; API Go inicial pronta.

| Doc | Descrição |
|---|---|
| [docs/data-model.md](docs/data-model.md) | Modelo de dados, RLS, fluxos |
| [docs/api.openapi.yaml](docs/api.openapi.yaml) | Contrato REST da API (OpenAPI 3.0) |
| [docs/ux-flows.md](docs/ux-flows.md) | Fluxos de UX, wireframes, copy |
| [api/README.md](api/README.md) | API Go |
| [docs/google-places-setup.md](docs/google-places-setup.md) | Tutorial Google Places API key |
| [mobile/README.md](mobile/README.md) | App Expo |
| `supabase/migrations/` | Migrations SQL |
| `supabase/seed/dev.sql` | Seed local (só Docker) |
| `docker-compose.yml` | Postgres + PostGIS local (porta **5433**) |
