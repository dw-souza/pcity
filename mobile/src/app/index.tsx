import AsyncStorage from '@react-native-async-storage/async-storage';
import { useQuery } from '@tanstack/react-query';
import { router, Stack } from 'expo-router';
import { useCallback, useEffect, useState } from 'react';
import {
  FlatList,
  Pressable,
  RefreshControl,
  StyleSheet,
  TextInput,
  View,
} from 'react-native';
import { SafeAreaView } from 'react-native-safe-area-context';

import { EmptyState } from '@/components/EmptyState';
import { LoadingView } from '@/components/LoadingView';
import { PlaceCard } from '@/components/PlaceCard';
import { ThemedText } from '@/components/themed-text';
import { useLocation } from '@/hooks/useLocation';
import { ApiClientError } from '@/lib/api/client';
import { listCityPlaces } from '@/lib/api/places';
import { getAccessToken } from '@/lib/auth';
import { CITY_LABEL, COPY, STORAGE_KEYS } from '@/lib/constants';

export default function HomeScreen() {
  const { coords, loading: locationLoading } = useLocation();
  const [search, setSearch] = useState('');
  const [debouncedSearch, setDebouncedSearch] = useState('');
  const [ready, setReady] = useState(false);

  useEffect(() => {
    const t = setTimeout(() => setDebouncedSearch(search.trim()), 300);
    return () => clearTimeout(t);
  }, [search]);

  useEffect(() => {
    AsyncStorage.getItem(STORAGE_KEYS.onboardingSeen).then((v) => {
      if (!v) {
        router.replace('/onboarding');
      } else {
        setReady(true);
      }
    });
  }, []);

  const query = useQuery({
    queryKey: ['places', coords?.lat, coords?.lng, debouncedSearch],
    enabled: ready && !locationLoading,
    queryFn: async () => {
      const token = await getAccessToken();
      return listCityPlaces({
        lat: coords?.lat,
        lng: coords?.lng,
        q: debouncedSearch || undefined,
        token,
      });
    },
    retry: (count, err) => {
      if (err instanceof ApiClientError && err.code === 'indexing_in_progress') {
        return count < 20;
      }
      return count < 2;
    },
    retryDelay: (attempt, err) => {
      if (err instanceof ApiClientError && err.code === 'indexing_in_progress') {
        return 3000;
      }
      return 1000;
    },
  });

  const onRefresh = useCallback(() => {
    void query.refetch();
  }, [query]);

  if (!ready || locationLoading || query.isLoading) {
    const indexing =
      query.error instanceof ApiClientError && query.error.code === 'indexing_in_progress';
    return (
      <>
        <Stack.Screen options={{ headerShown: false }} />
        <LoadingView message={indexing ? COPY.indexing : COPY.loadingPlaces} />
      </>
    );
  }

  if (query.isError && !(query.error instanceof ApiClientError && query.error.code === 'indexing_in_progress')) {
    return (
      <>
        <Stack.Screen options={{ headerShown: false }} />
        <SafeAreaView style={styles.safe}>
          <EmptyState message={COPY.networkError} />
          <Pressable onPress={onRefresh} style={styles.retry}>
            <ThemedText type="link">Tentar de novo</ThemedText>
          </Pressable>
        </SafeAreaView>
      </>
    );
  }

  const places = query.data?.data ?? [];

  return (
    <>
      <Stack.Screen options={{ headerShown: false }} />
      <SafeAreaView style={styles.safe}>
        <View style={styles.header}>
          <ThemedText type="title">Pcity</ThemedText>
          <Pressable onPress={() => router.push('/profile')}>
            <ThemedText type="link">Perfil</ThemedText>
          </Pressable>
        </View>

        <ThemedText type="small" themeColor="textSecondary" style={styles.city}>
          📍 {CITY_LABEL}
        </ThemedText>

        <TextInput
          placeholder="Buscar lugar..."
          value={search}
          onChangeText={setSearch}
          style={styles.search}
          placeholderTextColor="#888"
        />

        <FlatList
          data={places}
          keyExtractor={(item) => item.id}
          contentContainerStyle={places.length === 0 ? styles.emptyList : styles.list}
          refreshControl={<RefreshControl refreshing={query.isFetching} onRefresh={onRefresh} />}
          ListEmptyComponent={<EmptyState message={COPY.emptyList} />}
          renderItem={({ item }) => (
            <PlaceCard place={item} onPress={() => router.push(`/place/${item.id}`)} />
          )}
        />
      </SafeAreaView>
    </>
  );
}

const styles = StyleSheet.create({
  safe: { flex: 1, paddingHorizontal: 16 },
  header: {
    flexDirection: 'row',
    justifyContent: 'space-between',
    alignItems: 'center',
    paddingVertical: 8,
  },
  city: { marginBottom: 8 },
  search: {
    borderWidth: StyleSheet.hairlineWidth,
    borderColor: '#ccc',
    borderRadius: 10,
    paddingHorizontal: 12,
    paddingVertical: 10,
    marginBottom: 12,
    fontSize: 16,
  },
  list: { gap: 12, paddingBottom: 24 },
  emptyList: { flexGrow: 1 },
  retry: { alignItems: 'center', padding: 16 },
});
