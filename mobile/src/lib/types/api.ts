export type GeoPoint = { lat: number; lng: number };

export type PlaceCategory = 'bar' | 'restaurant' | 'cafe' | 'other';
export type PlaceSource = 'google' | 'user';
export type IndexingStatus = 'pending' | 'in_progress' | 'completed' | 'failed';

export type AmenityType = {
  id: string;
  label: string;
  icon: string;
  sort_order: number;
};

export type PlaceAmenity = {
  id: string;
  amenity_type_id: string;
  amenity_type?: AmenityType;
  confirmation_count: number;
  is_verified: boolean;
  first_reported_at: string;
  verified_at?: string | null;
  user_has_confirmed?: boolean;
};

export type PlaceSummary = {
  id: string;
  name: string;
  address: string;
  location: GeoPoint;
  phone?: string | null;
  category: PlaceCategory;
  source: PlaceSource;
  distance_m?: number | null;
  amenities?: PlaceAmenity[];
};

export type PlaceDetail = PlaceSummary & {
  website?: string | null;
  amenities: PlaceAmenity[];
  created_at: string;
  updated_at: string;
};

export type City = {
  id: string;
  slug: string;
  name: string;
  state: string;
  country: string;
  location: GeoPoint;
  indexing_status: IndexingStatus;
  indexed_at?: string | null;
};

export type PlaceListResponse = {
  data: PlaceSummary[];
  meta: { limit: number; offset: number; total: number };
  city: City;
};

export type Profile = {
  id: string;
  display_name?: string | null;
  avatar_url?: string | null;
  created_at: string;
  updated_at: string;
};

export type ApiError = {
  error: {
    code: string;
    message: string;
    city?: City;
    existing?: PlaceAmenity;
  };
};
