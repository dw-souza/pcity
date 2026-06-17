-- Seed de desenvolvimento local (só Docker).
-- O usuário dev é criado em 000000 antes do trigger de profiles existir.

INSERT INTO public.profiles (id, display_name)
VALUES ('00000000-0000-4000-8000-000000000001', 'Dev User')
ON CONFLICT (id) DO NOTHING;
