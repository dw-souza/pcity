import { useMutation, useQuery } from '@tanstack/react-query';
import { Alert, Pressable, StyleSheet, View } from 'react-native';

import { ThemedText } from '@/components/themed-text';
import { confirmAmenity, listAmenityTypes, reportAmenity } from '@/lib/api/amenities';
import { ApiClientError } from '@/lib/api/client';
import { getAccessToken, isLoggedIn } from '@/lib/auth';
import type { PlaceAmenity, PlaceDetail } from '@/lib/types/api';
import { useTheme } from '@/hooks/use-theme';
import { useEffect, useState } from 'react';

type Props = {
  place: PlaceDetail;
  onUpdated: () => void;
};

export function PlaceAmenitiesSection({ place, onUpdated }: Props) {
  const theme = useTheme();
  const [loggedIn, setLoggedIn] = useState(false);

  useEffect(() => {
    void isLoggedIn().then(setLoggedIn);
  }, []);

  const typesQuery = useQuery({
    queryKey: ['amenity-types'],
    enabled: loggedIn,
    queryFn: listAmenityTypes,
  });

  const reportMutation = useMutation({
    mutationFn: async (amenityTypeId: string) => {
      const token = await getAccessToken();
      if (!token) throw new Error('Faça login para continuar');
      return reportAmenity(place.id, amenityTypeId, token);
    },
    onSuccess: () => onUpdated(),
    onError: (err) => {
      const msg =
        err instanceof ApiClientError
          ? err.message
          : err instanceof Error
            ? err.message
            : 'Não foi possível reportar';
      Alert.alert('Ops', msg);
    },
  });

  const confirmMutation = useMutation({
    mutationFn: async (amenityTypeId: string) => {
      const token = await getAccessToken();
      if (!token) throw new Error('Faça login para continuar');
      return confirmAmenity(place.id, amenityTypeId, token);
    },
    onSuccess: () => onUpdated(),
    onError: (err) => {
      const msg =
        err instanceof ApiClientError
          ? err.message
          : err instanceof Error
            ? err.message
            : 'Não foi possível confirmar';
      Alert.alert('Ops', msg);
    },
  });

  const reportedTypeIds = new Set(place.amenities.map((a) => a.amenity_type_id));
  const availableTypes = (typesQuery.data ?? []).filter((t) => !reportedTypeIds.has(t.id));

  return (
    <View style={[styles.section, { borderColor: theme.border }]}>
      <ThemedText type="defaultSemiBold">Comodidades</ThemedText>

      {place.amenities.length === 0 ? (
        <ThemedText themeColor="textSecondary" style={styles.empty}>
          Ninguém reportou comodidades aqui ainda. Seja o primeiro!
        </ThemedText>
      ) : (
        place.amenities.map((a) => (
          <AmenityRow
            key={a.id}
            amenity={a}
            loading={confirmMutation.isPending}
            onConfirm={() => confirmMutation.mutate(a.amenity_type_id)}
          />
        ))
      )}

      {loggedIn && availableTypes.length > 0 ? (
        <>
          <ThemedText type="smallBold" style={styles.reportTitle}>
            Reportar comodidade
          </ThemedText>
          <View style={styles.typeGrid}>
            {availableTypes.map((t) => (
              <Pressable
                key={t.id}
                disabled={reportMutation.isPending}
                onPress={() => reportMutation.mutate(t.id)}
                style={[
                  styles.typeChip,
                  { backgroundColor: theme.backgroundSelected, borderColor: theme.border },
                ]}
              >
                <ThemedText>{t.icon}</ThemedText>
                <ThemedText type="small">{t.label}</ThemedText>
              </Pressable>
            ))}
          </View>
        </>
      ) : null}
    </View>
  );
}

function AmenityRow({
  amenity,
  onConfirm,
  loading,
}: {
  amenity: PlaceAmenity;
  onConfirm: () => void;
  loading: boolean;
}) {
  const theme = useTheme();
  const canConfirm = !amenity.is_verified && !amenity.user_has_confirmed;

  return (
    <View style={[styles.amenity, { borderBottomColor: theme.border }]}>
      <View style={styles.amenityHead}>
        <ThemedText type="defaultSemiBold">
          {amenity.amenity_type?.icon} {amenity.amenity_type?.label}
        </ThemedText>
        {amenity.is_verified ? (
          <View style={[styles.verified, { backgroundColor: theme.primaryMuted }]}>
            <ThemedText type="smallBold" style={{ color: theme.success }}>
              Verificado ✓
            </ThemedText>
          </View>
        ) : null}
      </View>
      <ThemedText type="small" themeColor="textSecondary">
        {amenity.confirmation_count} confirmação(ões)
      </ThemedText>
      {canConfirm ? (
        <Pressable
          disabled={loading}
          onPress={onConfirm}
          style={[styles.confirmBtn, { borderColor: theme.primary }]}
        >
          <ThemedText type="smallBold" style={{ color: theme.primary }}>
            Confirmar
          </ThemedText>
        </Pressable>
      ) : amenity.user_has_confirmed ? (
        <ThemedText type="small" themeColor="textSecondary">
          Você já confirmou
        </ThemedText>
      ) : null}
    </View>
  );
}

const styles = StyleSheet.create({
  section: {
    borderWidth: 1,
    borderRadius: 16,
    padding: 16,
    gap: 12,
  },
  empty: { lineHeight: 22 },
  reportTitle: { marginTop: 4 },
  typeGrid: { flexDirection: 'row', flexWrap: 'wrap', gap: 8 },
  typeChip: {
    flexDirection: 'row',
    alignItems: 'center',
    gap: 6,
    paddingHorizontal: 12,
    paddingVertical: 10,
    borderRadius: 12,
    borderWidth: 1,
  },
  amenity: {
    paddingVertical: 12,
    borderBottomWidth: StyleSheet.hairlineWidth,
    gap: 6,
  },
  amenityHead: {
    flexDirection: 'row',
    alignItems: 'center',
    justifyContent: 'space-between',
    gap: 8,
  },
  verified: {
    paddingHorizontal: 8,
    paddingVertical: 4,
    borderRadius: 999,
  },
  confirmBtn: {
    alignSelf: 'flex-start',
    paddingHorizontal: 12,
    paddingVertical: 6,
    borderRadius: 999,
    borderWidth: 1,
  },
});
