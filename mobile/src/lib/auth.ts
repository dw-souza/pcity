import AsyncStorage from '@react-native-async-storage/async-storage';
import { createClient, type Session } from '@supabase/supabase-js';
import { config, isSupabaseConfigured } from '@/lib/config';
import { STORAGE_KEYS } from '@/lib/constants';

let client: ReturnType<typeof createClient> | null = null;

export function getSupabase() {
  if (!isSupabaseConfigured()) {
    throw new Error('Supabase não configurado. Defina EXPO_PUBLIC_SUPABASE_URL e EXPO_PUBLIC_SUPABASE_ANON_KEY.');
  }
  if (!client) {
    client = createClient(config.supabaseUrl, config.supabaseAnonKey, {
      auth: {
        storage: AsyncStorage,
        autoRefreshToken: true,
        persistSession: true,
        detectSessionInUrl: false,
      },
    });
  }
  return client;
}

export async function getAccessToken(): Promise<string | null> {
  if (isSupabaseConfigured()) {
    const { data } = await getSupabase().auth.getSession();
    return data.session?.access_token ?? null;
  }
  return AsyncStorage.getItem(STORAGE_KEYS.devAccessToken);
}

export async function setDevAccessToken(token: string): Promise<void> {
  await AsyncStorage.setItem(STORAGE_KEYS.devAccessToken, token);
}

export async function getSession(): Promise<Session | null> {
  if (isSupabaseConfigured()) {
    const { data } = await getSupabase().auth.getSession();
    return data.session;
  }
  const token = await AsyncStorage.getItem(STORAGE_KEYS.devAccessToken);
  if (!token) return null;
  return { access_token: token } as Session;
}

export async function isLoggedIn(): Promise<boolean> {
  return Boolean(await getAccessToken());
}

export async function signOut(): Promise<void> {
  if (isSupabaseConfigured()) {
    await getSupabase().auth.signOut();
  }
  await AsyncStorage.removeItem(STORAGE_KEYS.devAccessToken);
}
