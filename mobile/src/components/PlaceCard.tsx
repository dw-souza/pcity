import { Pressable, StyleSheet, View } from 'react-native';

import { CategoryBadge } from '@/components/CategoryBadge';
import { PlaceThumbnail } from '@/components/PlaceThumbnail';
import { ThemedText } from '@/components/themed-text';
import { useTheme } from '@/hooks/use-theme';
import { formatDistance, shortAddress } from '@/lib/place-labels';
import type { PlaceSummary } from '@/lib/types/api';

type Props = {
  place: PlaceSummary;
  onPress: () => void;
};

export function PlaceCard({ place, onPress }: Props) {
  const theme = useTheme();
  const distance = formatDistance(place.distance_m);
  const amenities = (place.amenities ?? []).slice(0, 3);

  return (
    <Pressable
      onPress={onPress}
      style={({ pressed }) => [
        styles.card,
        {
          backgroundColor: theme.backgroundElement,
          shadowColor: '#1a1a1a',
        },
        pressed && styles.pressed,
      ]}
    >
      <PlaceThumbnail category={place.category} />

      <View style={styles.body}>
        <View style={styles.titleRow}>
          <ThemedText type="smallBold" numberOfLines={2} style={styles.name}>
            {place.name}
          </ThemedText>
          {distance ? (
            <View style={[styles.distance, { backgroundColor: theme.primaryMuted }]}>
              <ThemedText type="smallBold" style={{ color: theme.text }}>
                {distance}
              </ThemedText>
            </View>
          ) : null}
        </View>

        <CategoryBadge category={place.category} />

        <ThemedText type="small" themeColor="textSecondary" numberOfLines={2} style={styles.address}>
          {shortAddress(place.address)}
        </ThemedText>

        {amenities.length > 0 ? (
          <View style={styles.amenities}>
            {amenities.map((a) => (
              <ThemedText key={a.id} style={styles.amenityIcon}>
                {a.amenity_type?.icon}
              </ThemedText>
            ))}
          </View>
        ) : null}
      </View>

      <View style={[styles.chevron, { backgroundColor: theme.backgroundSelected }]}>
        <ThemedText style={styles.chevronIcon}>›</ThemedText>
      </View>
    </Pressable>
  );
}

const styles = StyleSheet.create({
  card: {
    flexDirection: 'row',
    alignItems: 'flex-start',
    padding: 18,
    borderRadius: 18,
    gap: 14,
    shadowOpacity: 0.08,
    shadowRadius: 16,
    shadowOffset: { width: 0, height: 4 },
    elevation: 3,
  },
  pressed: { opacity: 0.94, transform: [{ scale: 0.996 }] },
  body: { flex: 1, gap: 8, minWidth: 0 },
  titleRow: {
    flexDirection: 'row',
    alignItems: 'flex-start',
    justifyContent: 'space-between',
    gap: 10,
  },
  name: {
    flex: 1,
    fontSize: 17,
    lineHeight: 22,
    fontWeight: '700',
  },
  address: { lineHeight: 19 },
  distance: {
    paddingHorizontal: 10,
    paddingVertical: 5,
    borderRadius: 10,
    flexShrink: 0,
  },
  amenities: { flexDirection: 'row', gap: 6, marginTop: 2 },
  amenityIcon: { fontSize: 16 },
  chevron: {
    width: 32,
    height: 32,
    borderRadius: 16,
    alignItems: 'center',
    justifyContent: 'center',
    alignSelf: 'center',
  },
  chevronIcon: {
    fontSize: 22,
    lineHeight: 24,
    fontWeight: '600',
    marginTop: -2,
  },
});
