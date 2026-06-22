import type { PlaceCategory } from '@/lib/types/api';

const CATEGORY_LABELS: Record<PlaceCategory, string> = {
  bar: 'Bar',
  restaurant: 'Restaurante',
  cafe: 'Café',
  other: 'Outro',
};

const CATEGORY_ICONS: Record<PlaceCategory, string> = {
  bar: '🍺',
  restaurant: '🍽️',
  cafe: '☕',
  other: '📍',
};

export function categoryLabel(category: PlaceCategory): string {
  return CATEGORY_LABELS[category] ?? category;
}

export function categoryIcon(category: PlaceCategory): string {
  return CATEGORY_ICONS[category] ?? '📍';
}

export function formatDistance(m?: number | null): string | null {
  if (m == null) return null;
  if (m < 1000) return `${Math.round(m)} m`;
  return `${(m / 1000).toFixed(1)} km`;
}

export function shortAddress(full: string): string {
  const parts = full.split(',').map((s) => s.trim());
  if (parts.length >= 3) {
    const tail = parts.slice(-3, -1).join(', ');
    if (tail) return tail;
  }
  if (full.length > 72) return `${full.slice(0, 69)}…`;
  return full;
}
