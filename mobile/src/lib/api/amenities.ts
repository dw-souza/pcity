import { apiData, apiRequest, ApiClientError } from '@/lib/api/client';
import type { AmenityType, PlaceAmenity } from '@/lib/types/api';

export async function listAmenityTypes(): Promise<AmenityType[]> {
  const json = await apiRequest<{ data: AmenityType[] }>('/amenity-types');
  return apiData(json);
}

export async function reportAmenity(
  placeId: string,
  amenityTypeId: string,
  token: string,
): Promise<PlaceAmenity> {
  try {
    const json = await apiRequest<{ data: PlaceAmenity }>(`/places/${placeId}/amenities`, {
      method: 'POST',
      body: { amenity_type_id: amenityTypeId },
      token,
    });
    return apiData(json);
  } catch (err) {
    if (err instanceof ApiClientError && err.code === 'amenity_already_exists' && err.body?.existing) {
      return err.body.existing as PlaceAmenity;
    }
    throw err;
  }
}

export async function confirmAmenity(
  placeId: string,
  amenityTypeId: string,
  token: string,
): Promise<PlaceAmenity> {
  const json = await apiRequest<{ data: PlaceAmenity }>(
    `/places/${placeId}/amenities/${amenityTypeId}/confirm`,
    { method: 'POST', token },
  );
  return apiData(json);
}
