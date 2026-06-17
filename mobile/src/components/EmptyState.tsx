import { StyleSheet, View } from 'react-native';

import { ThemedText } from '@/components/themed-text';

type Props = {
  message: string;
  actionLabel?: string;
  onAction?: () => void;
};

export function EmptyState({ message }: Props) {
  return (
    <View style={styles.container}>
      <ThemedText type="subtitle" style={styles.text}>
        {message}
      </ThemedText>
    </View>
  );
}

const styles = StyleSheet.create({
  container: {
    flex: 1,
    alignItems: 'center',
    justifyContent: 'center',
    padding: 24,
  },
  text: { textAlign: 'center' },
});
