import { router } from 'expo-router';
import { useEffect, useState } from 'react';
import { Pressable, StyleSheet, View } from 'react-native';

import { ThemedText } from '@/components/themed-text';
import { getSession, getSupabase } from '@/lib/auth';
import { isSupabaseConfigured } from '@/lib/config';
import type { Profile } from '@/lib/types/api';
import { apiData, apiRequest } from '@/lib/api/client';

export default function ProfileScreen() {
  const [profile, setProfile] = useState<Profile | null>(null);
  const [loggedIn, setLoggedIn] = useState(false);

  useEffect(() => {
    void (async () => {
      const session = await getSession();
      setLoggedIn(Boolean(session));
      if (session) {
        try {
          const json = await apiRequest<{ data: Profile }>('/me', {
            token: session.access_token,
          });
          setProfile(apiData(json));
        } catch {
          setProfile(null);
        }
      }
    })();
  }, []);

  async function onSignOut() {
    if (isSupabaseConfigured()) {
      await getSupabase().auth.signOut();
    }
    setLoggedIn(false);
    setProfile(null);
    router.replace('/');
  }

  if (!loggedIn) {
    return (
      <View style={styles.container}>
        <ThemedText style={styles.message}>Entre para contribuir com o que os lugares têm.</ThemedText>
        <Pressable style={styles.button} onPress={() => router.push('/auth/login')}>
          <ThemedText type="defaultSemiBold" style={styles.buttonText}>
            Entrar
          </ThemedText>
        </Pressable>
      </View>
    );
  }

  return (
    <View style={styles.container}>
      <ThemedText type="title">{profile?.display_name ?? 'Usuário'}</ThemedText>
      <Pressable style={styles.outline} onPress={() => void onSignOut()}>
        <ThemedText>Sair</ThemedText>
      </Pressable>
    </View>
  );
}

const styles = StyleSheet.create({
  container: { flex: 1, padding: 24, gap: 16 },
  message: { textAlign: 'center', marginTop: 40 },
  button: {
    backgroundColor: '#f97316',
    padding: 14,
    borderRadius: 12,
    alignItems: 'center',
  },
  buttonText: { color: '#fff' },
  outline: {
    marginTop: 24,
    padding: 14,
    borderRadius: 12,
    borderWidth: StyleSheet.hairlineWidth,
    borderColor: '#ccc',
    alignItems: 'center',
  },
});
