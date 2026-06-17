import { useQuery } from '@tanstack/react-query';
import { router, useLocalSearchParams } from 'expo-router';
import { Pressable, ScrollView, StyleSheet, View } from 'react-native';

import { EmptyState } from '@/components/EmptyState';
import { LoadingView } from '@/components/LoadingView';
import { ThemedText } from '@/components/themed-text';
import { useLocation } from '@/hooks/useLocation';
import { getPlace } from '@/lib/api/places';
import { getAccessToken } from '@/lib/auth';
import { COPY } from '@/lib/constants';

export default function PlaceDetailScreen() {
  const { id } = useLocalSearchParams<{ id: string }>();
  const { coords } = useLocation();

  const query = useQuery({
    queryKey: ['place', id, coords?.lat, coords?.lng],
    enabled: Boolean(id),
    queryFn: async () => {
      const token = await getAccessToken();
      return getPlace(id!, coords ?? undefined, token);
    },
  });

  if (query.isLoading) {
    return <LoadingView />;
  }

  if (query.isError || !query.data) {
    return <EmptyState message={COPY.networkError} />;
  }

  const place = query.data;

  return (
    <ScrollView contentContainerStyle={styles.container}>
      <ThemedText type="title">{place.name}</ThemedText>
      <ThemedText type="small" themeColor="textSecondary">
        {place.category} · {place.address}
      </ThemedText>

      <ThemedText type="defaultSemiBold" style={styles.section}>
        Comodidades
      </ThemedText>

      {place.amenities.length === 0 ? (
        <ThemedText themeColor="textSecondary">{COPY.emptyAmenities}</ThemedText>
      ) : (
        place.amenities.map((a) => (
          <View key={a.id} style={styles.amenity}>
            <ThemedText>
              {a.amenity_type?.icon} {a.amenity_type?.label}
            </ThemedText>
            <ThemedText type="small" themeColor="textSecondary">
              {a.is_verified
                ? 'Verificado ✓'
                : `${a.confirmation_count} pessoa(s) confirmaram`}
            </ThemedText>
          </View>
        ))
      )}

      <Pressable
        style={styles.reportButton}
        onPress={() => router.push({ pathname: '/auth/login', params: { returnTo: `/place/${id}` } })}
      >
        <ThemedText type="defaultSemiBold" style={styles.reportText}>
          + Reportar comodidade
        </ThemedText>
      </Pressable>
    </ScrollView>
  );
}

const styles = StyleSheet.create({
  container: { padding: 16, gap: 8, paddingBottom: 32 },
  section: { marginTop: 16, marginBottom: 4 },
  amenity: {
    paddingVertical: 10,
    borderBottomWidth: StyleSheet.hairlineWidth,
    borderBottomColor: '#ddd',
    gap: 2,
  },
  reportButton: {
    marginTop: 24,
    backgroundColor: '#f97316',
    padding: 14,
    borderRadius: 12,
    alignItems: 'center',
  },
  reportText: { color: '#fff' },
});
