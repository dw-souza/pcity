# Pcity Mobile

App Expo (React Native) do Pcity.

## Setup

```bash
cd mobile
cp .env.example .env
npm install
npm start
```

## Variáveis de ambiente

| Variável | Descrição |
|---|---|
| `EXPO_PUBLIC_API_URL` | Base da API Go (`/v1`) |
| `EXPO_PUBLIC_SUPABASE_URL` | URL do projeto Supabase |
| `EXPO_PUBLIC_SUPABASE_ANON_KEY` | Anon key do Supabase |

### API URL por plataforma

| Onde roda | URL típica |
|---|---|
| iOS Simulator | `http://localhost:8080/v1` |
| Android Emulator | `http://10.0.2.2:8080/v1` |
| Celular físico | `http://<IP-da-sua-máquina>:8080/v1` |

## Estrutura

```
src/
  app/           # Expo Router (telas)
  components/    # UI
  lib/api/       # Cliente HTTP
  hooks/         # useLocation, etc.
```

## Telas (MVP)

| Rota | Tela |
|---|---|
| `/` | Lista de lugares |
| `/onboarding` | Boas-vindas (1ª vez) |
| `/place/[id]` | Detalhe |
| `/auth/login` | Login |
| `/profile` | Perfil |

## Stack

- Expo SDK 56 + Expo Router
- TanStack Query
- expo-location
- Supabase Auth (quando configurado)

## Google Places

A key do Google fica **somente na API Go** — ver [docs/google-places-setup.md](../docs/google-places-setup.md).
