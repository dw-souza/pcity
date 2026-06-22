import { router, useLocalSearchParams } from 'expo-router';
import { useState } from 'react';
import { Alert, Pressable, StyleSheet, TextInput, View } from 'react-native';

import { ThemedText } from '@/components/themed-text';
import { ThemedView } from '@/components/themed-view';
import { devRegister, supportsDevAuth } from '@/lib/api/dev-auth';
import { config, isSupabaseConfigured } from '@/lib/config';
import { getSupabase } from '@/lib/auth';
import { COPY } from '@/lib/constants';
import { useTheme } from '@/hooks/use-theme';

type Mode = 'login' | 'register';

export default function LoginScreen() {
  const { returnTo } = useLocalSearchParams<{ returnTo?: string }>();
  const theme = useTheme();
  const [mode, setMode] = useState<Mode>('register');
  const [displayName, setDisplayName] = useState('');
  const [email, setEmail] = useState('');
  const [password, setPassword] = useState('');
  const [loading, setLoading] = useState(false);
  const devAuth = supportsDevAuth();

  function goNext() {
    if (returnTo && typeof returnTo === 'string') {
      router.replace(returnTo as `/place/${string}`);
    } else {
      router.replace('/');
    }
  }

  async function onSubmit() {
    const trimmedEmail = email.trim().toLowerCase();
    if (!trimmedEmail) {
      Alert.alert('E-mail obrigatório', 'Informe seu e-mail.');
      return;
    }
    if (mode === 'register' && !displayName.trim()) {
      Alert.alert('Nome obrigatório', 'Como você quer aparecer no app?');
      return;
    }
    if (isSupabaseConfigured() && !password) {
      Alert.alert('Senha obrigatória', 'Informe sua senha.');
      return;
    }

    setLoading(true);
    try {
      if (isSupabaseConfigured()) {
        if (mode === 'register') {
          const { error } = await getSupabase().auth.signUp({
            email: trimmedEmail,
            password,
            options: { data: { display_name: displayName.trim() } },
          });
          if (error) throw error;
          Alert.alert('Conta criada', 'Se pedir confirmação por e-mail, confira sua caixa de entrada.');
        } else {
          const { error } = await getSupabase().auth.signInWithPassword({
            email: trimmedEmail,
            password,
          });
          if (error) throw error;
        }
      } else if (devAuth) {
        const name =
          mode === 'register'
            ? displayName.trim()
            : displayName.trim() || trimmedEmail.split('@')[0] || 'Usuário';
        await devRegister(trimmedEmail, name);
      } else {
        Alert.alert(
          'Auth não configurado',
          'Configure Supabase no mobile/.env ou DEV_AUTH=true na API.',
        );
        return;
      }
      goNext();
    } catch (e) {
      Alert.alert(
        mode === 'register' ? 'Erro ao criar conta' : 'Erro ao entrar',
        e instanceof Error ? e.message : 'Tente novamente',
      );
    } finally {
      setLoading(false);
    }
  }

  return (
    <ThemedView style={styles.container}>
      <ThemedText type="subtitle" style={styles.logo}>
        Pcity
      </ThemedText>
      <ThemedText themeColor="textSecondary" style={styles.subtitle}>
        {COPY.loginCta}
      </ThemedText>

      <View style={[styles.tabs, { backgroundColor: theme.backgroundElement, borderColor: theme.border }]}>
        <Pressable
          style={[styles.tab, mode === 'login' && { backgroundColor: theme.backgroundSelected }]}
          onPress={() => setMode('login')}
        >
          <ThemedText type="smallBold">Entrar</ThemedText>
        </Pressable>
        <Pressable
          style={[styles.tab, mode === 'register' && { backgroundColor: theme.backgroundSelected }]}
          onPress={() => setMode('register')}
        >
          <ThemedText type="smallBold">Criar conta</ThemedText>
        </Pressable>
      </View>

      {mode === 'register' || !devAuth ? (
        <TextInput
          placeholder="Seu nome"
          value={displayName}
          onChangeText={setDisplayName}
          style={[styles.input, { borderColor: theme.border, color: theme.text }]}
          placeholderTextColor={theme.textSecondary}
        />
      ) : null}

      <TextInput
        placeholder="E-mail"
        autoCapitalize="none"
        keyboardType="email-address"
        value={email}
        onChangeText={setEmail}
        style={[styles.input, { borderColor: theme.border, color: theme.text }]}
        placeholderTextColor={theme.textSecondary}
      />

      {isSupabaseConfigured() ? (
        <TextInput
          placeholder="Senha"
          secureTextEntry
          value={password}
          onChangeText={setPassword}
          style={[styles.input, { borderColor: theme.border, color: theme.text }]}
          placeholderTextColor={theme.textSecondary}
        />
      ) : null}

      <Pressable
        style={[styles.button, { backgroundColor: theme.primary }, loading && styles.disabled]}
        onPress={() => void onSubmit()}
        disabled={loading}
      >
        <ThemedText type="defaultSemiBold" style={styles.buttonText}>
          {loading ? 'Aguarde...' : mode === 'register' ? 'Criar conta' : 'Entrar'}
        </ThemedText>
      </Pressable>

      {devAuth ? (
        <ThemedText type="small" themeColor="textSecondary" style={styles.hint}>
          Dev local: conta salva no Postgres via {config.apiUrl}/dev/auth
        </ThemedText>
      ) : null}
    </ThemedView>
  );
}

const styles = StyleSheet.create({
  container: { flex: 1, padding: 24, justifyContent: 'center', gap: 12 },
  logo: { textAlign: 'center', fontSize: 32, lineHeight: 38 },
  subtitle: { textAlign: 'center', marginBottom: 8 },
  tabs: {
    flexDirection: 'row',
    borderWidth: 1,
    borderRadius: 12,
    padding: 4,
    gap: 4,
  },
  tab: {
    flex: 1,
    alignItems: 'center',
    paddingVertical: 10,
    borderRadius: 8,
  },
  input: {
    borderWidth: 1,
    borderRadius: 12,
    padding: 14,
    fontSize: 16,
  },
  button: {
    padding: 16,
    borderRadius: 14,
    alignItems: 'center',
    marginTop: 8,
  },
  disabled: { opacity: 0.6 },
  buttonText: { color: '#fff' },
  hint: { textAlign: 'center', marginTop: 8, lineHeight: 20 },
});
