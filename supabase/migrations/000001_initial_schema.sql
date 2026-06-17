-- Pcity MVP — schema inicial
-- Requer: Postgres 15+ com PostGIS

CREATE EXTENSION IF NOT EXISTS postgis;
CREATE EXTENSION IF NOT EXISTS pgcrypto;

-- ---------------------------------------------------------------------------
-- Enums
-- ---------------------------------------------------------------------------

CREATE TYPE indexing_status AS ENUM (
  'pending',
  'in_progress',
  'completed',
  'failed'
);

CREATE TYPE place_source AS ENUM (
  'google',
  'user'
);

CREATE TYPE place_category AS ENUM (
  'bar',
  'restaurant',
  'cafe',
  'other'
);

-- ---------------------------------------------------------------------------
-- cities
-- ---------------------------------------------------------------------------

CREATE TABLE cities (
  id               uuid PRIMARY KEY DEFAULT gen_random_uuid(),
  slug             text NOT NULL UNIQUE,
  name             text NOT NULL,
  state            text NOT NULL,
  country          text NOT NULL DEFAULT 'BR',
  location         geography(Point, 4326) NOT NULL,
  indexed_at       timestamptz,
  indexing_status  indexing_status NOT NULL DEFAULT 'pending',
  created_at       timestamptz NOT NULL DEFAULT now(),
  updated_at       timestamptz NOT NULL DEFAULT now()
);

CREATE INDEX cities_location_idx ON cities USING GIST (location);

-- ---------------------------------------------------------------------------
-- places
-- ---------------------------------------------------------------------------

CREATE TABLE places (
  id               uuid PRIMARY KEY DEFAULT gen_random_uuid(),
  city_id          uuid NOT NULL REFERENCES cities (id) ON DELETE CASCADE,
  google_place_id  text UNIQUE,
  name             text NOT NULL,
  address          text NOT NULL DEFAULT '',
  location         geography(Point, 4326) NOT NULL,
  phone            text,
  website          text,
  category         place_category NOT NULL DEFAULT 'other',
  source           place_source NOT NULL DEFAULT 'google',
  created_by       uuid REFERENCES auth.users (id) ON DELETE SET NULL,
  is_active        boolean NOT NULL DEFAULT true,
  created_at       timestamptz NOT NULL DEFAULT now(),
  updated_at       timestamptz NOT NULL DEFAULT now(),

  CONSTRAINT places_google_id_required_for_google_source
    CHECK (source != 'google' OR google_place_id IS NOT NULL)
);

CREATE INDEX places_city_id_idx ON places (city_id);
CREATE INDEX places_location_idx ON places USING GIST (location);
CREATE INDEX places_active_city_idx ON places (city_id) WHERE is_active = true;

-- ---------------------------------------------------------------------------
-- amenity_types (catálogo fixo)
-- ---------------------------------------------------------------------------

CREATE TABLE amenity_types (
  id          text PRIMARY KEY,
  label       text NOT NULL,
  icon        text NOT NULL,
  sort_order  smallint NOT NULL DEFAULT 0,
  is_active   boolean NOT NULL DEFAULT true
);

-- ---------------------------------------------------------------------------
-- place_amenities
-- ---------------------------------------------------------------------------

CREATE TABLE place_amenities (
  id                  uuid PRIMARY KEY DEFAULT gen_random_uuid(),
  place_id            uuid NOT NULL REFERENCES places (id) ON DELETE CASCADE,
  amenity_type_id     text NOT NULL REFERENCES amenity_types (id),
  confirmation_count  int NOT NULL DEFAULT 1,
  is_verified         boolean NOT NULL DEFAULT false,
  first_reported_by   uuid NOT NULL REFERENCES auth.users (id) ON DELETE CASCADE,
  first_reported_at   timestamptz NOT NULL DEFAULT now(),
  verified_at         timestamptz,
  created_at          timestamptz NOT NULL DEFAULT now(),
  updated_at          timestamptz NOT NULL DEFAULT now(),

  UNIQUE (place_id, amenity_type_id),
  CONSTRAINT place_amenities_count_positive CHECK (confirmation_count >= 1)
);

CREATE INDEX place_amenities_place_id_idx ON place_amenities (place_id);

-- ---------------------------------------------------------------------------
-- amenity_confirmations
-- ---------------------------------------------------------------------------

CREATE TABLE amenity_confirmations (
  id                uuid PRIMARY KEY DEFAULT gen_random_uuid(),
  place_amenity_id  uuid NOT NULL REFERENCES place_amenities (id) ON DELETE CASCADE,
  user_id           uuid NOT NULL REFERENCES auth.users (id) ON DELETE CASCADE,
  created_at        timestamptz NOT NULL DEFAULT now(),

  UNIQUE (place_amenity_id, user_id)
);

CREATE INDEX amenity_confirmations_user_id_idx ON amenity_confirmations (user_id);

-- ---------------------------------------------------------------------------
-- profiles
-- ---------------------------------------------------------------------------

CREATE TABLE profiles (
  id            uuid PRIMARY KEY REFERENCES auth.users (id) ON DELETE CASCADE,
  display_name  text,
  avatar_url    text,
  created_at    timestamptz NOT NULL DEFAULT now(),
  updated_at    timestamptz NOT NULL DEFAULT now()
);

-- ---------------------------------------------------------------------------
-- Triggers: updated_at
-- ---------------------------------------------------------------------------

CREATE OR REPLACE FUNCTION set_updated_at()
RETURNS trigger
LANGUAGE plpgsql
AS $$
BEGIN
  NEW.updated_at = now();
  RETURN NEW;
END;
$$;

CREATE TRIGGER cities_set_updated_at
  BEFORE UPDATE ON cities
  FOR EACH ROW EXECUTE FUNCTION set_updated_at();

CREATE TRIGGER places_set_updated_at
  BEFORE UPDATE ON places
  FOR EACH ROW EXECUTE FUNCTION set_updated_at();

CREATE TRIGGER place_amenities_set_updated_at
  BEFORE UPDATE ON place_amenities
  FOR EACH ROW EXECUTE FUNCTION set_updated_at();

CREATE TRIGGER profiles_set_updated_at
  BEFORE UPDATE ON profiles
  FOR EACH ROW EXECUTE FUNCTION set_updated_at();

-- ---------------------------------------------------------------------------
-- Trigger: auto-profile on signup
-- ---------------------------------------------------------------------------

CREATE OR REPLACE FUNCTION handle_new_user()
RETURNS trigger
LANGUAGE plpgsql
SECURITY DEFINER
SET search_path = public
AS $$
BEGIN
  INSERT INTO public.profiles (id, display_name, avatar_url)
  VALUES (
    NEW.id,
    COALESCE(NEW.raw_user_meta_data ->> 'display_name', NEW.raw_user_meta_data ->> 'full_name'),
    NEW.raw_user_meta_data ->> 'avatar_url'
  );
  RETURN NEW;
END;
$$;

CREATE TRIGGER on_auth_user_created
  AFTER INSERT ON auth.users
  FOR EACH ROW EXECUTE FUNCTION handle_new_user();

-- ---------------------------------------------------------------------------
-- Trigger: increment confirmation count + verify at threshold
-- ---------------------------------------------------------------------------

CREATE OR REPLACE FUNCTION handle_amenity_confirmation()
RETURNS trigger
LANGUAGE plpgsql
AS $$
DECLARE
  new_count int;
  threshold constant int := 3;
BEGIN
  UPDATE place_amenities
  SET
    confirmation_count = confirmation_count + 1,
    is_verified = CASE
      WHEN confirmation_count + 1 >= threshold THEN true
      ELSE is_verified
    END,
    verified_at = CASE
      WHEN confirmation_count + 1 >= threshold AND verified_at IS NULL THEN now()
      ELSE verified_at
    END
  WHERE id = NEW.place_amenity_id
  RETURNING confirmation_count INTO new_count;

  RETURN NEW;
END;
$$;

CREATE OR REPLACE FUNCTION prevent_self_confirmation()
RETURNS trigger
LANGUAGE plpgsql
AS $$
DECLARE
  reporter uuid;
BEGIN
  SELECT first_reported_by INTO reporter
  FROM place_amenities
  WHERE id = NEW.place_amenity_id;

  IF reporter = NEW.user_id THEN
    RAISE EXCEPTION 'Cannot confirm your own report';
  END IF;

  RETURN NEW;
END;
$$;

CREATE TRIGGER amenity_confirmation_prevent_self
  BEFORE INSERT ON amenity_confirmations
  FOR EACH ROW EXECUTE FUNCTION prevent_self_confirmation();

CREATE TRIGGER amenity_confirmation_increment
  AFTER INSERT ON amenity_confirmations
  FOR EACH ROW EXECUTE FUNCTION handle_amenity_confirmation();

-- ---------------------------------------------------------------------------
-- Row Level Security
-- ---------------------------------------------------------------------------

ALTER TABLE cities ENABLE ROW LEVEL SECURITY;
ALTER TABLE places ENABLE ROW LEVEL SECURITY;
ALTER TABLE amenity_types ENABLE ROW LEVEL SECURITY;
ALTER TABLE place_amenities ENABLE ROW LEVEL SECURITY;
ALTER TABLE amenity_confirmations ENABLE ROW LEVEL SECURITY;
ALTER TABLE profiles ENABLE ROW LEVEL SECURITY;

-- Leitura pública
CREATE POLICY cities_select_public ON cities
  FOR SELECT USING (true);

CREATE POLICY places_select_public ON places
  FOR SELECT USING (is_active = true);

CREATE POLICY amenity_types_select_public ON amenity_types
  FOR SELECT USING (is_active = true);

CREATE POLICY place_amenities_select_public ON place_amenities
  FOR SELECT USING (true);

CREATE POLICY profiles_select_public ON profiles
  FOR SELECT USING (true);

-- Contribuições autenticadas
CREATE POLICY place_amenities_insert_authenticated ON place_amenities
  FOR INSERT
  WITH CHECK (auth.uid() = first_reported_by);

CREATE POLICY amenity_confirmations_select_own ON amenity_confirmations
  FOR SELECT USING (auth.uid() = user_id);

CREATE POLICY amenity_confirmations_insert_own ON amenity_confirmations
  FOR INSERT
  WITH CHECK (auth.uid() = user_id);

CREATE POLICY profiles_update_own ON profiles
  FOR UPDATE USING (auth.uid() = id);

-- ---------------------------------------------------------------------------
-- Seed: amenity types
-- ---------------------------------------------------------------------------

INSERT INTO amenity_types (id, label, icon, sort_order) VALUES
  ('kids_area',    'Espaço kids',    '👶', 1),
  ('pet_friendly', 'Pet friendly',   '🐕', 2),
  ('outdoor',      'Área externa',   '☀️', 3),
  ('live_music',   'Música ao vivo', '🎵', 4),
  ('parking',      'Estacionamento', '🅿️', 5),
  ('wifi',         'Wi-Fi',          '📶', 6),
  ('accessible',   'Acessível',      '♿', 7);

-- ---------------------------------------------------------------------------
-- Seed: cidade piloto — Franca, SP
-- Centro aproximado: -20.5386, -47.4008
-- ---------------------------------------------------------------------------

INSERT INTO cities (slug, name, state, location, indexing_status)
VALUES (
  'franca-sp',
  'Franca',
  'SP',
  ST_SetSRID(ST_MakePoint(-47.4008, -20.5386), 4326)::geography,
  'pending'
);
