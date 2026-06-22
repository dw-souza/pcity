# Google Places API — Guia completo de configuração

Tutorial para criar conta no Google Cloud, ativar a Places API, gerar uma API key segura e usar no Pcity (seed de bares/restaurantes em Franca).

**Tempo estimado:** 20–30 minutos  
**Custo:** Google oferece **US$ 200/mês de crédito gratuito** no Google Maps Platform. Indexar Franca uma vez custa ~US$ 10–15 (ver [estimativa](#estimativa-de-custo-para-franca)).

---

## O que você vai precisar

- Conta Google (Gmail)
- Cartão de crédito (obrigatório para ativar billing, mesmo usando crédito grátis)
- Acesso ao projeto Pcity (`api/.env`)

---

## Visão geral do fluxo

```
1. Criar projeto no Google Cloud
2. Ativar faturamento (billing)
3. Ativar Places API
4. Criar API key
5. Restringir a key (segurança)
6. Configurar alerta de budget
7. Colocar a key no api/.env
8. Testar indexação
```

---

## Passo 1 — Criar projeto no Google Cloud

1. Acesse [Google Cloud Console](https://console.cloud.google.com/)
2. No topo, clique no seletor de projetos → **Novo projeto**
3. Nome sugerido: `pcity` (ou `pcity-prod`)
4. Clique em **Criar**
5. Aguarde alguns segundos e **selecione o projeto** no seletor do topo

---

## Passo 2 — Ativar faturamento (billing)

A Places API exige billing ativo, mesmo com crédito gratuito.

1. Menu ☰ → **Faturamento** (ou [console.cloud.google.com/billing](https://console.cloud.google.com/billing))
2. Se não tiver conta de faturamento:
   - Clique em **Criar conta de faturamento**
   - Preencha país, endereço e **cartão de crédito**
   - Aceite os termos
3. Volte ao projeto `pcity` e **vincule a conta de faturamento** ao projeto

> O cartão não é cobrado automaticamente além do uso. Com crédito de US$ 200/mês, o MVP de uma cidade fica dentro do gratuito.

---

## Passo 3 — Ativar a Places API

O Pcity usa **Places API (New) — Text Search** via endpoint:
`POST https://places.googleapis.com/v1/places:searchText`

1. Menu ☰ → **APIs e serviços** → **Biblioteca**
2. Busque: `Places API (New)`
3. Clique em **Ativar**

**Atalho direto:** [Ativar Places API (New)](https://console.cloud.google.com/apis/library/places.googleapis.com)

> Não é necessário ativar a Places API legada (`places-backend.googleapis.com`).

---

## Passo 4 — Criar a API key

1. Menu ☰ → **APIs e serviços** → **Credenciais**
2. **+ Criar credenciais** → **Chave de API**
3. A key será exibida (formato `AIza...`) — **copie e guarde**

Não commite essa key no Git. Ela vai em `api/.env`:

```env
GOOGLE_PLACES_API_KEY=AIza...sua-key-aqui
```

---

## Passo 5 — Restringir a API key (importante)

Keys sem restrição podem ser usadas por terceiros se vazarem — e geram cobrança.

1. Em **Credenciais**, clique na key criada
2. Em **Restrições de aplicativo**:
   - Para **desenvolvimento local** (API Go no seu PC): escolha **Endereços IP**
     - Adicione seu IP público (busque em [whatismyip.com](https://www.whatismyip.com/))
     - Para servidor fixo depois: IP do Railway/Fly/etc.
   - **Não** use restrição de HTTP referrer para a API Go (isso é para frontends JS no browser)
3. Em **Restrições de API**:
   - Selecione **Restringir key**
   - Marque apenas: **Places API (New)**
4. **Salvar**

### Desenvolvimento local com IP dinâmico

Se seu IP muda (Wi‑Fi doméstico), opções:

| Opção | Quando usar |
|---|---|
| Atualizar IP na restrição | Testes esporádicos |
| Key separada sem restrição de IP (só Places API (New)) | Dev temporário — **apague depois** |
| Rodar indexação só em CI/servidor com IP fixo | Mais seguro |

---

## Passo 6 — Alerta de budget (recomendado)

1. Menu ☰ → **Faturamento** → **Orçamentos e alertas**
2. **Criar orçamento**
3. Nome: `Pcity alerta`
4. Valor: ex. **US$ 10** (ou US$ 50 se preferir margem)
5. Configure alertas em **50%**, **90%** e **100%**
6. Adicione seu e-mail

Assim você recebe aviso antes de qualquer cobrança inesperada.

---

## Passo 7 — Configurar no Pcity

### API Go (`api/.env`)

```env
GOOGLE_PLACES_API_KEY=AIza...sua-key
ADMIN_KEY=dev-admin-key
DATABASE_URL=postgresql://postgres:postgres@localhost:5433/pcity?sslmode=disable
```

### Subir stack local

```bash
# Na raiz do monorepo
docker compose up -d

cd api
go run ./cmd/server
```

### Testar a key (curl direto no Google)

Substitua `SUA_KEY`:

```bash
curl -s -X POST 'https://places.googleapis.com/v1/places:searchText' \
  -H 'Content-Type: application/json' \
  -H "X-Goog-Api-Key: SUA_KEY" \
  -H 'X-Goog-FieldMask: places.id,places.displayName,places.formattedAddress,places.location,nextPageToken' \
  -d '{"textQuery":"bares em Franca, SP, Brasil","languageCode":"pt-BR","pageSize":5}' \
  | jq '{total: (.places | length), exemplo: .places[0].displayName}'
```

Resposta esperada: lista em `places` (campo `displayName.text` com o nome do lugar).

Erros comuns:

| sinal | Causa | Solução |
|---|---|---|
| `PERMISSION_DENIED` / HTTP 403 | Key inválida, API não ativada ou billing inativo | Ver passos 2–4 |
| `RESOURCE_EXHAUSTED` | Cota excedida | Aguardar ou revisar billing |
| `INVALID_ARGUMENT` | FieldMask ausente ou body malformado | Ver curl acima |

### Indexar Franca pelo Pcity

```bash
curl -X POST http://localhost:8080/v1/admin/cities/franca-sp/index \
  -H "X-Admin-Key: dev-admin-key"
```

Resposta esperada:

```json
{
  "data": {
    "city": { "slug": "franca-sp", "indexing_status": "completed", ... },
    "places_imported": 150,
    "indexing_status": "completed"
  }
}
```

### Verificar no banco

```bash
make psql
# SELECT count(*) FROM places;
```

---

## Estimativa de custo para Franca

A indexação divide a cidade em **grade 3×3** (~9 células de 5 km) e busca **6 categorias** em cada célula (bares, restaurantes, lanchonetes, cafés, pizzarias, hamburguerias). A Google limita **60 resultados por busca**; a grade compensa esse teto para cobrir a cidade inteira.

| Operação | Qtd aproximada | Custo |
|---|---|---|
| Text Search (9 células × 6 categorias × ~2 páginas) | ~100 requests | ~US$ 3–8 |
| Lugares importados (deduplicados) | 300–600+ | incluído no Text Search |
| **Total indexação inicial** | | **~US$ 5–10** |
| Uso mensal após seed | ~0 | Dados ficam no Postgres |

O Pcity **cacheia no banco** — Google Places só é chamado de novo ao indexar **outra cidade**.

Crédito mensal: **US$ 200** → sobra ampla margem para o MVP.

---

## Segurança — checklist

- [ ] Key apenas em `api/.env` (nunca no mobile, nunca no Git)
- [ ] `.env` no `.gitignore` ✓ (já configurado)
- [ ] Restrição de API: só Places API (New)
- [ ] Restrição de IP em produção
- [ ] Alerta de budget configurado
- [ ] Key de dev separada da key de produção (recomendado)

---

## Produção (quando deployar a API)

1. Crie **outra API key** para produção
2. Restrinja ao IP do servidor (Railway, Fly.io, etc.)
3. Configure `GOOGLE_PLACES_API_KEY` nas variáveis de ambiente do host
4. Rode indexação **uma vez** antes de abrir o app:

```bash
curl -X POST https://api.pcity.app/v1/admin/cities/franca-sp/index \
  -H "X-Admin-Key: SEU_ADMIN_KEY_PRODUCAO"
```

---

## Troubleshooting

### `google places api key not configured`

A variável `GOOGLE_PLACES_API_KEY` está vazia no `api/.env`. Reinicie o servidor após editar.

### Indexação retorna `places_imported: 0`

- Key sem permissão / billing inativo
- Restrição de IP bloqueando seu servidor
- Teste o curl direto no Google (passo 7)

### `indexing_status: failed`

Veja logs da API Go. Causa usual: `PERMISSION_DENIED` do Google (API (New) não ativada ou key restrita).

### Quero usar OpenStreetMap em vez do Google

Possível como plano B (Overpass API, gratuito), mas qualidade em cidades menores varia. O código atual usa Google Places; trocar exigiria novo importer.

---

## Links úteis

- [Google Cloud Console](https://console.cloud.google.com/)
- [Places API (New) — Text Search](https://developers.google.com/maps/documentation/places/web-service/text-search)
- [Preços Google Maps Platform](https://mapsplatform.google.com/pricing/)
- [Credenciais — boas práticas](https://cloud.google.com/docs/authentication/api-keys)

---

## Próximo passo

Com a key funcionando e Franca indexada, o app mobile (`mobile/`) consome `GET /v1/cities/franca-sp/places` — sem precisar da key no celular.
