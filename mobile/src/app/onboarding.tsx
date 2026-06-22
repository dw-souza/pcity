import AsyncStorage from '@react-native-async-storage/async-storage';
import { router } from 'expo-router';
import { Pressable, StyleSheet, View } from 'react-native';
import { SafeAreaView } from 'react-native-safe-area-context';

import { ThemedText } from '@/components/themed-text';
import { ThemedView } from '@/components/themed-view';
import { MaxContentWidth } from '@/constants/theme';
import { useTheme } from '@/hooks/use-theme';
import { COPY, STORAGE_KEYS } from '@/lib/constants';

export default function OnboardingScreen() {
  const theme = useTheme();

  async function onStart() {
    await AsyncStorage.setItem(STORAGE_KEYS.onboardingSeen, '1');
    router.replace('/');
  }

  return (
    <ThemedView style={styles.safe}>
      <SafeAreaView style={styles.safeInner}>
        <View style={styles.content}>
          <View style={[styles.badge, { backgroundColor: theme.primaryMuted }]}>
            <ThemedText type="smallBold" style={{ color: theme.primary }}>
              Franca, SP
            </ThemedText>
          </View>
          <ThemedText type="title" style={styles.logo}>
            Pcity
          </ThemedText>
          <ThemedText type="subtitle" style={styles.tagline}>
            {COPY.tagline}
          </ThemedText>
          <ThemedText themeColor="textSecondary" style={styles.body}>
            Bares e restaurantes de Franca com o que importa: espaço kids, pet friendly e muito
            mais.
          </ThemedText>
          <Pressable
            onPress={() => void onStart()}
            style={({ pressed }) => [
              styles.button,
              { backgroundColor: theme.primary },
              pressed && styles.pressed,
            ]}
          >
            <ThemedText type="defaultSemiBold" style={styles.buttonText}>
              Explorar lugares
            </ThemedText>
          </Pressable>
        </View>
      </SafeAreaView>
    </ThemedView>
  );
}

const styles = StyleSheet.create({
  safe: { flex: 1 },
  safeInner: { flex: 1 },
  content: {
    flex: 1,
    justifyContent: 'center',
    alignItems: 'center',
    padding: 24,
    gap: 16,
    maxWidth: MaxContentWidth,
    alignSelf: 'center',
    width: '100%',
  },
  badge: {
    paddingHorizontal: 12,
    paddingVertical: 6,
    borderRadius: 999,
  },
  logo: { fontSize: 44, lineHeight: 48 },
  tagline: { fontSize: 28, lineHeight: 34, textAlign: 'center' },
  body: { textAlign: 'center', lineHeight: 24, maxWidth: 320 },
  button: {
    marginTop: 12,
    paddingHorizontal: 28,
    paddingVertical: 16,
    borderRadius: 14,
    minWidth: 220,
    alignItems: 'center',
  },
  buttonText: { color: '#fff' },
  pressed: { opacity: 0.9 },
});
