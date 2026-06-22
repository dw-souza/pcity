import AsyncStorage from '@react-native-async-storage/async-storage';
import { useInfiniteQuery } from '@tanstack/react-query';
import { router, Stack } from 'expo-router';
import { useCallback, useEffect, useMemo, useState } from 'react';
import {
  ActivityIndicator,
  FlatList,
  Pressable,
  RefreshControl,
  StyleSheet,
  View,
} from 'react-native';
import { SafeAreaView } from 'react-native-safe-area-context';

import {
  CategoryFilterBar,
  type CategoryFilter,
} from '@/components/CategoryFilterBar';
import { EmptyState } from '@/components/EmptyState';
import { LoadingView } from '@/components/LoadingView';
import { PlaceCard } from '@/components/PlaceCard';
import { ResultsToolbar } from '@/components/ResultsToolbar';
import { ScreenHeader } from '@/components/ScreenHeader';
import { SearchField } from '@/components/SearchField';
import type { SortOption } from '@/components/SortSelector';
import { ThemedText } from '@/components/themed-text';
import { ThemedView } from '@/components/themed-view';
import { MaxContentWidth, Spacing } from '@/constants/theme';
import { useLocation } from '@/hooks/useLocation';
import { useTheme } from '@/hooks/use-theme';
import { ApiClientError } from '@/lib/api/client';
import { listCityPlaces } from '@/lib/api/places';
import { getAccessToken } from '@/lib/auth';
import { CITY_LABEL, COPY, STORAGE_KEYS } from '@/lib/constants';
import type { PlaceSummary } from '@/lib/types/api';

const PAGE_SIZE = 30;

function sortPlaces(places: PlaceSummary[], sort: SortOption): PlaceSummary[] {
  if (sort === 'name') {
    return [...places].sort((a, b) => a.name.localeCompare(b.name, 'pt-BR'));
  }
  return places;
}

export default function HomeScreen() {
  const theme = useTheme();
  const { coords, loading: locationLoading } = useLocation();
  const [search, setSearch] = useState('');
  const [debouncedSearch, setDebouncedSearch] = useState('');
  const [category, setCategory] = useState<CategoryFilter>('all');
  const [sort, setSort] = useState<SortOption>('distance');
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

  const query = useInfiniteQuery({
    queryKey: ['places', coords?.lat, coords?.lng, debouncedSearch, category],
    enabled: ready && !locationLoading,
    initialPageParam: 0,
    queryFn: async ({ pageParam }) => {
      const token = await getAccessToken();
      return listCityPlaces({
        lat: coords?.lat,
        lng: coords?.lng,
        q: debouncedSearch || undefined,
        category: category === 'all' ? undefined : category,
        limit: PAGE_SIZE,
        offset: pageParam,
        token,
      });
    },
    getNextPageParam: (lastPage, pages) => {
      const loaded = pages.reduce((sum, page) => sum + page.data.length, 0);
      return loaded < lastPage.meta.total ? loaded : undefined;
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

  const places = useMemo(() => {
    const flat = query.data?.pages.flatMap((page) => page.data) ?? [];
    return sortPlaces(flat, sort);
  }, [query.data, sort]);

  const total = query.data?.pages[0]?.meta.total ?? 0;

  const onRefresh = useCallback(() => {
    void query.refetch();
  }, [query]);

  const onEndReached = useCallback(() => {
    if (query.hasNextPage && !query.isFetchingNextPage) {
      void query.fetchNextPage();
    }
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
        <ThemedView style={styles.safe}>
          <SafeAreaView style={styles.safeInner}>
            <EmptyState message={COPY.networkError} />
            <Pressable onPress={onRefresh} style={styles.retry}>
              <ThemedText type="link">Tentar de novo</ThemedText>
            </Pressable>
          </SafeAreaView>
        </ThemedView>
      </>
    );
  }

  return (
    <>
      <Stack.Screen options={{ headerShown: false }} />
      <ThemedView style={styles.safe}>
        <SafeAreaView style={styles.safeInner} edges={['top']}>
          <View style={styles.content}>
            <ScreenHeader
              title="Pcity"
              subtitle={`📍 ${CITY_LABEL}`}
              actionLabel="Perfil"
              onActionPress={() => router.push('/profile')}
            />

            <View style={styles.controls}>
              <SearchField value={search} onChangeText={setSearch} />
              <CategoryFilterBar value={category} onChange={setCategory} />
              <ResultsToolbar
                sort={sort}
                onSortChange={setSort}
                shown={places.length}
                total={total}
              />
            </View>

            <FlatList
              data={places}
              keyExtractor={(item) => item.id}
              contentContainerStyle={places.length === 0 ? styles.emptyList : styles.list}
              refreshControl={
                <RefreshControl
                  refreshing={query.isFetching && !query.isFetchingNextPage}
                  onRefresh={onRefresh}
                  tintColor={theme.primary}
                />
              }
              onEndReached={onEndReached}
              onEndReachedThreshold={0.4}
              ListEmptyComponent={<EmptyState message={COPY.emptyList} />}
              ListFooterComponent={
                query.isFetchingNextPage ? (
                  <View style={styles.footerLoader}>
                    <ActivityIndicator color={theme.primary} />
                  </View>
                ) : null
              }
              renderItem={({ item }) => (
                <PlaceCard place={item} onPress={() => router.push(`/place/${item.id}`)} />
              )}
            />
          </View>
        </SafeAreaView>
      </ThemedView>
    </>
  );
}

const styles = StyleSheet.create({
  safe: { flex: 1 },
  safeInner: { flex: 1 },
  content: {
    flex: 1,
    width: '100%',
    maxWidth: MaxContentWidth,
    alignSelf: 'center',
  },
  controls: {
    paddingHorizontal: Spacing.three,
    gap: 14,
    paddingBottom: 12,
  },
  list: { gap: 18, paddingHorizontal: Spacing.three, paddingTop: 4, paddingBottom: 32 },
  emptyList: { flexGrow: 1, paddingHorizontal: Spacing.three },
  footerLoader: { paddingVertical: 16 },
  retry: { alignItems: 'center', padding: 16 },
});
