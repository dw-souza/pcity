import { router } from 'expo-router';
import { useEffect, useState } from 'react';
import { Pressable, StyleSheet, View } from 'react-native';

import { ThemedText } from '@/components/themed-text';
import { ThemedView } from '@/components/themed-view';
import { useTheme } from '@/hooks/use-theme';
import { apiData, apiRequest } from '@/lib/api/client';
import { getAccessToken, isLoggedIn, signOut } from '@/lib/auth';
import type { Profile } from '@/lib/types/api';

export default function ProfileScreen() {
  const theme = useTheme();
  const [profile, setProfile] = useState<Profile | null>(null);
  const [loggedIn, setLoggedIn] = useState(false);

  useEffect(() => {
    void (async () => {
      const ok = await isLoggedIn();
      setLoggedIn(ok);
      if (ok) {
        try {
          const token = await getAccessToken();
          const json = await apiRequest<{ data: Profile }>('/me', { token });
          setProfile(apiData(json));
        } catch {
          setProfile(null);
        }
      }
    })();
  }, []);

  async function onSignOut() {
    await signOut();
    setLoggedIn(false);
    setProfile(null);
    router.replace('/');
  }

  if (!loggedIn) {
    return (
      <ThemedView style={styles.container}>
        <ThemedText style={styles.message}>Entre para contribuir com o que os lugares têm.</ThemedText>
        <Pressable
          style={[styles.button, { backgroundColor: theme.primary }]}
          onPress={() => router.push('/auth/login')}
        >
          <ThemedText type="defaultSemiBold" style={styles.buttonText}>
            Entrar ou criar conta
          </ThemedText>
        </Pressable>
      </ThemedView>
    );
  }

  return (
    <ThemedView style={styles.container}>
      <ThemedText type="subtitle">{profile?.display_name ?? 'Usuário'}</ThemedText>
      <ThemedText themeColor="textSecondary">
        Você pode reportar e confirmar comodidades nos lugares.
      </ThemedText>
      <Pressable
        style={[styles.outline, { borderColor: theme.border }]}
        onPress={() => void onSignOut()}
      >
        <ThemedText>Sair</ThemedText>
      </Pressable>
    </ThemedView>
  );
}

const styles = StyleSheet.create({
  container: { flex: 1, padding: 24, gap: 16 },
  message: { textAlign: 'center', marginTop: 40 },
  button: {
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
    alignItems: 'center',
  },
});
