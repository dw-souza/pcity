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

## Status

MVP em planejamento.
