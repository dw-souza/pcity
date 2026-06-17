-- Stub mínimo do Supabase Auth para desenvolvimento local (Docker).
-- Em produção o schema auth já existe no Supabase; IF NOT EXISTS torna isso seguro.

CREATE SCHEMA IF NOT EXISTS auth;

CREATE TABLE IF NOT EXISTS auth.users (
  id                 uuid PRIMARY KEY DEFAULT gen_random_uuid(),
  email              text,
  raw_user_meta_data jsonb NOT NULL DEFAULT '{}'::jsonb,
  created_at         timestamptz NOT NULL DEFAULT now()
);

-- Compatível com Supabase RLS (local: null quando não há JWT no contexto)
CREATE OR REPLACE FUNCTION auth.uid()
RETURNS uuid
LANGUAGE sql
STABLE
AS $$
  SELECT NULLIF(current_setting('request.jwt.claim.sub', true), '')::uuid;
$$;

-- Usuário fixo para testes locais (JWT sub = este id)
INSERT INTO auth.users (id, email, raw_user_meta_data)
VALUES (
  '00000000-0000-4000-8000-000000000001',
  'dev@pcity.local',
  '{"display_name": "Dev User"}'::jsonb
)
ON CONFLICT (id) DO NOTHING;
