import { Pressable, StyleSheet, View } from 'react-native';

import { ThemedText } from '@/components/themed-text';
import type { PlaceSummary } from '@/lib/types/api';

type Props = {
  place: PlaceSummary;
  onPress: () => void;
};

function formatDistance(m?: number | null): string | null {
  if (m == null) return null;
  if (m < 1000) return `${Math.round(m)} m`;
  return `${(m / 1000).toFixed(1)} km`;
}

export function PlaceCard({ place, onPress }: Props) {
  const distance = formatDistance(place.distance_m);
  const icons = (place.amenities ?? [])
    .slice(0, 5)
    .map((a) => a.amenity_type?.icon ?? '')
    .join(' ');

  return (
    <Pressable onPress={onPress} style={({ pressed }) => [styles.card, pressed && styles.pressed]}>
      <ThemedText type="defaultSemiBold">{place.name}</ThemedText>
      <ThemedText type="small" themeColor="textSecondary">
        {place.address}
        {distance ? ` · ${distance}` : ''}
      </ThemedText>
      {icons ? (
        <View style={styles.icons}>
          <ThemedText>{icons}</ThemedText>
        </View>
      ) : null}
    </Pressable>
  );
}

const styles = StyleSheet.create({
  card: {
    padding: 16,
    borderRadius: 12,
    borderWidth: StyleSheet.hairlineWidth,
    borderColor: '#ccc',
    gap: 4,
  },
  pressed: { opacity: 0.85 },
  icons: { marginTop: 4 },
});
