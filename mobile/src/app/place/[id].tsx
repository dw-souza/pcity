import { useQuery } from '@tanstack/react-query';
import { router, Stack, useLocalSearchParams } from 'expo-router';
import { Linking, Pressable, ScrollView, StyleSheet, View } from 'react-native';
import { useSafeAreaInsets } from 'react-native-safe-area-context';

import { CategoryBadge } from '@/components/CategoryBadge';
import { EmptyState } from '@/components/EmptyState';
import { LoadingView } from '@/components/LoadingView';
import { PlaceAmenitiesSection } from '@/components/PlaceAmenitiesSection';
import { ThemedText } from '@/components/themed-text';
import { ThemedView } from '@/components/themed-view';
import { MaxContentWidth, Spacing } from '@/constants/theme';
import { useLocation } from '@/hooks/useLocation';
import { useTheme } from '@/hooks/use-theme';
import { getPlace } from '@/lib/api/places';
import { getAccessToken, isLoggedIn } from '@/lib/auth';
import { COPY } from '@/lib/constants';
import { categoryLabel, formatDistance } from '@/lib/place-labels';
import { useEffect, useState } from 'react';

export default function PlaceDetailScreen() {
  const { id } = useLocalSearchParams<{ id: string }>();
  const { coords } = useLocation();
  const theme = useTheme();
  const insets = useSafeAreaInsets();
  const [loggedIn, setLoggedIn] = useState(false);

  const query = useQuery({
    queryKey: ['place', id, coords?.lat, coords?.lng],
    enabled: Boolean(id),
    queryFn: async () => {
      const token = await getAccessToken();
      return getPlace(id!, coords ?? undefined, token);
    },
  });

  useEffect(() => {
    void isLoggedIn().then(setLoggedIn);
  }, [query.dataUpdatedAt]);

  if (query.isLoading) {
    return <LoadingView />;
  }

  if (query.isError || !query.data) {
    return <EmptyState message={COPY.networkError} />;
  }

  const place = query.data;
  const distance = formatDistance(place.distance_m);

  function openMaps() {
    const q = encodeURIComponent(`${place.name}, ${place.address}`);
    void Linking.openURL(`https://www.google.com/maps/search/?api=1&query=${q}`);
  }

  return (
    <>
      <Stack.Screen
        options={{
          title: '',
          headerBackTitle: 'Voltar',
          headerTintColor: theme.primary,
        }}
      />
      <ThemedView style={styles.safe}>
        <ScrollView
          contentContainerStyle={[
            styles.container,
            { paddingBottom: insets.bottom + 24 },
          ]}
        >
          <View style={styles.content}>
            <View style={[styles.hero, { backgroundColor: theme.primaryMuted }]}>
              <CategoryBadge category={place.category} />
              <ThemedText type="subtitle" style={styles.name}>
                {place.name}
              </ThemedText>
              <ThemedText themeColor="textSecondary" style={styles.address}>
                {place.address}
              </ThemedText>
              {distance ? (
                <ThemedText type="smallBold" style={{ color: theme.primary }}>
                  📍 {distance} de você
                </ThemedText>
              ) : null}
            </View>

            <View style={[styles.section, { borderColor: theme.border }]}>
              <ThemedText type="defaultSemiBold">Informações</ThemedText>
              <InfoRow label="Tipo" value={categoryLabel(place.category)} />
              {place.phone ? (
                <Pressable onPress={() => void Linking.openURL(`tel:${place.phone}`)}>
                  <InfoRow label="Telefone" value={place.phone} link />
                </Pressable>
              ) : null}
              {place.website ? (
                <Pressable onPress={() => void Linking.openURL(place.website!)}>
                  <InfoRow label="Site" value="Abrir site" link />
                </Pressable>
              ) : null}
              <Pressable onPress={openMaps}>
                <InfoRow label="Mapa" value="Abrir no Google Maps" link />
              </Pressable>
            </View>

            <PlaceAmenitiesSection place={place} onUpdated={() => void query.refetch()} />

            {!loggedIn ? (
              <Pressable
                style={[styles.reportButton, { backgroundColor: theme.primary }]}
                onPress={() =>
                  router.push({ pathname: '/auth/login', params: { returnTo: `/place/${id}` } })
                }
              >
                <ThemedText type="defaultSemiBold" style={styles.reportText}>
                  Entrar para reportar comodidades
                </ThemedText>
              </Pressable>
            ) : null}
          </View>
        </ScrollView>
      </ThemedView>
    </>
  );
}

function InfoRow({
  label,
  value,
  link,
}: {
  label: string;
  value: string;
  link?: boolean;
}) {
  return (
    <View style={styles.infoRow}>
      <ThemedText type="small" themeColor="textSecondary">
        {label}
      </ThemedText>
      <ThemedText type="defaultSemiBold" style={link ? styles.linkValue : undefined}>
        {value}
      </ThemedText>
    </View>
  );
}

const styles = StyleSheet.create({
  safe: { flex: 1 },
  container: { flexGrow: 1 },
  content: {
    width: '100%',
    maxWidth: MaxContentWidth,
    alignSelf: 'center',
    padding: Spacing.three,
    gap: 16,
  },
  hero: {
    borderRadius: 20,
    padding: 20,
    gap: 10,
  },
  name: { fontSize: 28, lineHeight: 34 },
  address: { lineHeight: 22 },
  section: {
    borderWidth: 1,
    borderRadius: 16,
    padding: 16,
    gap: 12,
  },
  infoRow: { gap: 2 },
  linkValue: { color: '#2563eb' },
  reportButton: {
    marginTop: 8,
    padding: 16,
    borderRadius: 14,
    alignItems: 'center',
  },
  reportText: { color: '#fff' },
});
