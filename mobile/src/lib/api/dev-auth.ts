import { apiData, apiRequest } from '@/lib/api/client';
import { setDevAccessToken } from '@/lib/auth';
import { isSupabaseConfigured } from '@/lib/config';

type DevAuthResponse = {
  access_token: string;
  user: {
    id: string;
    email: string;
    display_name: string;
  };
};

export async function devRegister(email: string, displayName: string): Promise<DevAuthResponse> {
  const json = await apiRequest<{ data: DevAuthResponse }>('/dev/auth', {
    method: 'POST',
    body: { email, display_name: displayName },
  });
  const data = apiData(json);
  await setDevAccessToken(data.access_token);
  return data;
}

export function supportsDevAuth(): boolean {
  return !isSupabaseConfigured();
}
