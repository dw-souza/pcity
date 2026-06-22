import { StyleSheet, View } from 'react-native';

import { ThemedText } from '@/components/themed-text';
import { useTheme } from '@/hooks/use-theme';
import { categoryIcon, categoryLabel } from '@/lib/place-labels';
import type { PlaceCategory } from '@/lib/types/api';

type Props = {
  category: PlaceCategory;
};

export function CategoryBadge({ category }: Props) {
  const theme = useTheme();

  return (
    <View style={[styles.badge, { backgroundColor: theme.backgroundSelected }]}>
      <ThemedText style={styles.icon}>{categoryIcon(category)}</ThemedText>
      <ThemedText type="small" themeColor="textSecondary">
        {categoryLabel(category)}
      </ThemedText>
    </View>
  );
}

const styles = StyleSheet.create({
  badge: {
    flexDirection: 'row',
    alignItems: 'center',
    alignSelf: 'flex-start',
    gap: 4,
    paddingHorizontal: 10,
    paddingVertical: 5,
    borderRadius: 10,
  },
  icon: { fontSize: 13 },
});
