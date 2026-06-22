import { StyleSheet, View } from 'react-native';

import { SortSelector, type SortOption } from '@/components/SortSelector';
import { ThemedText } from '@/components/themed-text';
import { useTheme } from '@/hooks/use-theme';

type Props = {
  sort: SortOption;
  onSortChange: (value: SortOption) => void;
  shown: number;
  total: number;
};

export function ResultsToolbar({ sort, onSortChange, shown, total }: Props) {
  const theme = useTheme();

  return (
    <View style={[styles.wrap, { backgroundColor: theme.backgroundElement, borderColor: theme.border }]}>
      <View style={styles.countBlock}>
        <ThemedText type="smallBold">
          {total > 0 ? `${shown} de ${total} lugares` : 'Nenhum lugar'}
        </ThemedText>
        <ThemedText type="small" themeColor="textSecondary">
          em Franca
        </ThemedText>
      </View>
      <View style={styles.sortWrap}>
        <SortSelector value={sort} onChange={onSortChange} />
      </View>
    </View>
  );
}

const styles = StyleSheet.create({
  wrap: {
    borderWidth: 1,
    borderRadius: 16,
    padding: 12,
    gap: 12,
  },
  countBlock: { gap: 2 },
  sortWrap: { width: '100%' },
});
