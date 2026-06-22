import { Pressable, ScrollView, StyleSheet } from 'react-native';

import { ThemedText } from '@/components/themed-text';
import { useTheme } from '@/hooks/use-theme';
import { categoryIcon, categoryLabel } from '@/lib/place-labels';
import type { PlaceCategory } from '@/lib/types/api';

export type CategoryFilter = PlaceCategory | 'all';

const FILTERS: { id: CategoryFilter; label: string; icon: string }[] = [
  { id: 'all', label: 'Todos', icon: '✦' },
  { id: 'bar', label: categoryLabel('bar'), icon: categoryIcon('bar') },
  { id: 'restaurant', label: categoryLabel('restaurant'), icon: categoryIcon('restaurant') },
  { id: 'cafe', label: categoryLabel('cafe'), icon: categoryIcon('cafe') },
  { id: 'other', label: categoryLabel('other'), icon: categoryIcon('other') },
];

type Props = {
  value: CategoryFilter;
  onChange: (value: CategoryFilter) => void;
};

export function CategoryFilterBar({ value, onChange }: Props) {
  const theme = useTheme();

  return (
    <ScrollView
      horizontal
      showsHorizontalScrollIndicator={false}
      contentContainerStyle={styles.row}
    >
      {FILTERS.map((filter) => {
        const active = value === filter.id;
        return (
          <Pressable
            key={filter.id}
            onPress={() => onChange(filter.id)}
            style={[
              styles.chip,
              {
                backgroundColor: active ? theme.backgroundSelected : theme.backgroundElement,
                borderColor: active ? theme.primary : theme.border,
              },
            ]}
          >
            <ThemedText style={styles.chipIcon}>{filter.icon}</ThemedText>
            <ThemedText
              type="smallBold"
              style={{ color: active ? theme.primary : theme.textSecondary }}
            >
              {filter.label}
            </ThemedText>
          </Pressable>
        );
      })}
    </ScrollView>
  );
}

const styles = StyleSheet.create({
  row: { gap: 8, paddingVertical: 2 },
  chip: {
    flexDirection: 'row',
    alignItems: 'center',
    gap: 5,
    paddingHorizontal: 14,
    paddingVertical: 9,
    borderRadius: 12,
    borderWidth: 1,
  },
  chipIcon: { fontSize: 14 },
});
