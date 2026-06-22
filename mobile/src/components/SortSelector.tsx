import { Pressable, StyleSheet, View } from 'react-native';

import { ThemedText } from '@/components/themed-text';
import { useTheme } from '@/hooks/use-theme';

export type SortOption = 'distance' | 'name';

const OPTIONS: { id: SortOption; label: string }[] = [
  { id: 'distance', label: 'Mais perto' },
  { id: 'name', label: 'Nome A–Z' },
];

type Props = {
  value: SortOption;
  onChange: (value: SortOption) => void;
};

export function SortSelector({ value, onChange }: Props) {
  const theme = useTheme();

  return (
    <View style={styles.wrap}>
      {OPTIONS.map((option) => {
        const active = value === option.id;
        return (
          <Pressable
            key={option.id}
            onPress={() => onChange(option.id)}
            style={[
              styles.option,
              active && { backgroundColor: theme.backgroundSelected },
            ]}
          >
            <ThemedText
              type="smallBold"
              style={{ color: active ? theme.primary : theme.textSecondary }}
            >
              {option.label}
            </ThemedText>
          </Pressable>
        );
      })}
    </View>
  );
}

const styles = StyleSheet.create({
  wrap: {
    flexDirection: 'row',
    borderRadius: 12,
    padding: 4,
    gap: 4,
    backgroundColor: 'transparent',
  },
  option: {
    flex: 1,
    alignItems: 'center',
    paddingVertical: 8,
    borderRadius: 8,
  },
});
