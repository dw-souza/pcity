import { router, useLocalSearchParams } from 'expo-router';
import { useState } from 'react';
import { Alert, Pressable, StyleSheet, TextInput, View } from 'react-native';

import { ThemedText } from '@/components/themed-text';
import { COPY } from '@/lib/constants';
import { config, isSupabaseConfigured } from '@/lib/config';
import { getSupabase } from '@/lib/auth';

export default function LoginScreen() {
  const { returnTo } = useLocalSearchParams<{ returnTo?: string }>();
  const [email, setEmail] = useState('');
  const [password, setPassword] = useState('');
  const [loading, setLoading] = useState(false);

  async function onLogin() {
    if (!isSupabaseConfigured()) {
      Alert.alert(
        'Supabase não configurado',
        'Defina EXPO_PUBLIC_SUPABASE_URL e EXPO_PUBLIC_SUPABASE_ANON_KEY no .env',
      );
      return;
    }

    setLoading(true);
    try {
      const { error } = await getSupabase().auth.signInWithPassword({ email, password });
      if (error) throw error;
      if (returnTo && typeof returnTo === 'string') {
        router.replace(returnTo as `/place/${string}`);
      } else {
        router.replace('/');
      }
    } catch (e) {
      Alert.alert('Erro ao entrar', e instanceof Error ? e.message : 'Tente novamente');
    } finally {
      setLoading(false);
    }
  }

  return (
    <View style={styles.container}>
      <ThemedText type="title" style={styles.logo}>
        Pcity
      </ThemedText>
      <ThemedText style={styles.subtitle}>{COPY.loginCta}</ThemedText>

      <TextInput
        placeholder="E-mail"
        autoCapitalize="none"
        keyboardType="email-address"
        value={email}
        onChangeText={setEmail}
        style={styles.input}
      />
      <TextInput
        placeholder="Senha"
        secureTextEntry
        value={password}
        onChangeText={setPassword}
        style={styles.input}
      />

      <Pressable
        style={[styles.button, loading && styles.disabled]}
        onPress={() => void onLogin()}
        disabled={loading}
      >
        <ThemedText type="defaultSemiBold" style={styles.buttonText}>
          {loading ? 'Entrando...' : 'Entrar'}
        </ThemedText>
      </Pressable>

      {!isSupabaseConfigured() ? (
        <ThemedText type="small" themeColor="textSecondary" style={styles.hint}>
          Dev sem Supabase: API em {config.apiUrl}
        </ThemedText>
      ) : null}
    </View>
  );
}

const styles = StyleSheet.create({
  container: { flex: 1, padding: 24, justifyContent: 'center', gap: 12 },
  logo: { textAlign: 'center', marginBottom: 8 },
  subtitle: { textAlign: 'center', marginBottom: 16 },
  input: {
    borderWidth: StyleSheet.hairlineWidth,
    borderColor: '#ccc',
    borderRadius: 10,
    padding: 12,
    fontSize: 16,
  },
  button: {
    backgroundColor: '#f97316',
    padding: 14,
    borderRadius: 12,
    alignItems: 'center',
    marginTop: 8,
  },
  disabled: { opacity: 0.6 },
  buttonText: { color: '#fff' },
  hint: { textAlign: 'center', marginTop: 16 },
});
