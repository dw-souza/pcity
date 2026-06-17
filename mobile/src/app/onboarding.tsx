import AsyncStorage from '@react-native-async-storage/async-storage';
import { router } from 'expo-router';
import { Pressable, StyleSheet, View } from 'react-native';
import { SafeAreaView } from 'react-native-safe-area-context';

import { ThemedText } from '@/components/themed-text';
import { COPY, STORAGE_KEYS } from '@/lib/constants';

export default function OnboardingScreen() {
  async function onStart() {
    await AsyncStorage.setItem(STORAGE_KEYS.onboardingSeen, '1');
    router.replace('/');
  }

  return (
    <SafeAreaView style={styles.safe}>
      <View style={styles.content}>
        <ThemedText type="title" style={styles.logo}>
          Pcity
        </ThemedText>
        <ThemedText type="subtitle" style={styles.tagline}>
          {COPY.tagline}
        </ThemedText>
        <ThemedText style={styles.body}>
          Bares e restaurantes de Franca com o que importa: espaço kids, pet friendly e muito mais.
        </ThemedText>
        <Pressable onPress={() => void onStart()} style={styles.button}>
          <ThemedText type="defaultSemiBold" style={styles.buttonText}>
            Começar
          </ThemedText>
        </Pressable>
      </View>
    </SafeAreaView>
  );
}

const styles = StyleSheet.create({
  safe: { flex: 1 },
  content: {
    flex: 1,
    justifyContent: 'center',
    alignItems: 'center',
    padding: 24,
    gap: 16,
  },
  logo: { fontSize: 40 },
  tagline: { textAlign: 'center' },
  body: { textAlign: 'center', lineHeight: 22 },
  button: {
    marginTop: 16,
    backgroundColor: '#f97316',
    paddingHorizontal: 32,
    paddingVertical: 14,
    borderRadius: 12,
  },
  buttonText: { color: '#fff' },
});
