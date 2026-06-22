import { apiData, apiRequest } from '@/lib/api/client';
import { CITY_SLUG } from '@/lib/constants';
import type { PlaceDetail, PlaceListResponse } from '@/lib/types/api';

type ListParams = {
  lat?: number;
  lng?: number;
  q?: string;
  category?: string;
  limit?: number;
  offset?: number;
  token?: string | null;
};

export async function listCityPlaces(params: ListParams = {}): Promise<PlaceListResponse> {
  const search = new URLSearchParams();
  if (params.lat != null) search.set('lat', String(params.lat));
  if (params.lng != null) search.set('lng', String(params.lng));
  if (params.q) search.set('q', params.q);
  if (params.category) search.set('category', params.category);
  if (params.limit != null) search.set('limit', String(params.limit));
  if (params.offset != null) search.set('offset', String(params.offset));

  const qs = search.toString();
  const path = `/cities/${CITY_SLUG}/places${qs ? `?${qs}` : ''}`;

  return apiRequest<PlaceListResponse>(path, { token: params.token });
}

export async function getPlace(
  id: string,
  coords?: { lat: number; lng: number },
  token?: string | null,
): Promise<PlaceDetail> {
  const search = new URLSearchParams();
  if (coords) {
    search.set('lat', String(coords.lat));
    search.set('lng', String(coords.lng));
  }
  const qs = search.toString();
  const json = await apiRequest<{ data: PlaceDetail }>(
    `/places/${id}${qs ? `?${qs}` : ''}`,
    { token },
  );
  return apiData(json);
}
