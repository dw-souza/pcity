import { StyleSheet, View } from 'react-native';

import { ThemedText } from '@/components/themed-text';
import { useTheme } from '@/hooks/use-theme';
import { categoryIcon } from '@/lib/place-labels';
import type { PlaceCategory } from '@/lib/types/api';

type Props = {
  category: PlaceCategory;
  size?: 'sm' | 'md';
};

export function PlaceThumbnail({ category, size = 'md' }: Props) {
  const theme = useTheme();
  const dim = size === 'sm' ? 48 : 56;

  return (
    <View
      style={[
        styles.wrap,
        {
          width: dim,
          height: dim,
          borderRadius: 14,
          backgroundColor: theme.backgroundSelected,
        },
      ]}
    >
      <ThemedText style={[styles.icon, size === 'sm' && styles.iconSm]}>
        {categoryIcon(category)}
      </ThemedText>
    </View>
  );
}

const styles = StyleSheet.create({
  wrap: {
    alignItems: 'center',
    justifyContent: 'center',
  },
  icon: { fontSize: 26 },
  iconSm: { fontSize: 22 },
});
