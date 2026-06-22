package repository

import (
	"context"
	"errors"
	"fmt"
	"strings"

	"github.com/dw-souza/pcity/api/internal/apperrors"
	"github.com/dw-souza/pcity/api/internal/models"
	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgconn"
	"github.com/jackc/pgx/v5/pgxpool"
)

type Repository struct {
	pool *pgxpool.Pool
}

func New(pool *pgxpool.Pool) *Repository {
	return &Repository{pool: pool}
}

func (r *Repository) GetCityBySlug(ctx context.Context, slug string) (models.City, error) {
	const q = `
		SELECT id, slug, name, state, country,
		       ST_Y(location::geometry), ST_X(location::geometry),
		       indexing_status, indexed_at
		FROM cities WHERE slug = $1`

	var c models.City
	err := r.pool.QueryRow(ctx, q, slug).Scan(
		&c.ID, &c.Slug, &c.Name, &c.State, &c.Country,
		&c.Location.Lat, &c.Location.Lng,
		&c.IndexingStatus, &c.IndexedAt,
	)
	if errors.Is(err, pgx.ErrNoRows) {
		return c, apperrors.ErrNotFound
	}
	return c, err
}

func (r *Repository) TryClaimIndexing(ctx context.Context, slug string) (models.City, bool, error) {
	const q = `
		UPDATE cities
		SET indexing_status = 'in_progress', updated_at = now()
		WHERE slug = $1 AND indexing_status IN ('pending', 'failed', 'completed')
		RETURNING id, slug, name, state, country,
		          ST_Y(location::geometry), ST_X(location::geometry),
		          indexing_status, indexed_at`

	var c models.City
	err := r.pool.QueryRow(ctx, q, slug).Scan(
		&c.ID, &c.Slug, &c.Name, &c.State, &c.Country,
		&c.Location.Lat, &c.Location.Lng,
		&c.IndexingStatus, &c.IndexedAt,
	)
	if errors.Is(err, pgx.ErrNoRows) {
		city, err := r.GetCityBySlug(ctx, slug)
		if err != nil {
			return c, false, err
		}
		if city.IndexingStatus == models.IndexingInProgress {
			return city, false, apperrors.ErrIndexingInProgress
		}
		return city, false, nil
	}
	return c, true, err
}

func (r *Repository) FinishIndexing(ctx context.Context, cityID string, status models.IndexingStatus) (models.City, error) {
	const q = `
		UPDATE cities
		SET indexing_status = $2::indexing_status,
		    indexed_at = CASE WHEN $2::text = 'completed' THEN now() ELSE indexed_at END,
		    updated_at = now()
		WHERE id = $1
		RETURNING id, slug, name, state, country,
		          ST_Y(location::geometry), ST_X(location::geometry),
		          indexing_status, indexed_at`

	var c models.City
	err := r.pool.QueryRow(ctx, q, cityID, status).Scan(
		&c.ID, &c.Slug, &c.Name, &c.State, &c.Country,
		&c.Location.Lat, &c.Location.Lng,
		&c.IndexingStatus, &c.IndexedAt,
	)
	return c, err
}

func (r *Repository) ListAmenityTypes(ctx context.Context) ([]models.AmenityType, error) {
	const q = `SELECT id, label, icon, sort_order FROM amenity_types WHERE is_active = true ORDER BY sort_order`
	rows, err := r.pool.Query(ctx, q)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var out []models.AmenityType
	for rows.Next() {
		var a models.AmenityType
		if err := rows.Scan(&a.ID, &a.Label, &a.Icon, &a.SortOrder); err != nil {
			return nil, err
		}
		out = append(out, a)
	}
	return out, rows.Err()
}

func (r *Repository) CountPlaces(ctx context.Context, cityID string, category *models.PlaceCategory, search string) (int, error) {
	where, args := placeFilters(cityID, category, search)
	const base = `SELECT COUNT(*) FROM places p WHERE p.is_active = true`
	var total int
	err := r.pool.QueryRow(ctx, base+where, args...).Scan(&total)
	return total, err
}

func (r *Repository) ListPlaces(ctx context.Context, p models.ListPlacesParams) ([]models.PlaceSummary, error) {
	where, args := placeFilters(p.CityID, p.Category, p.Search)
	argN := len(args)

	orderBy := "p.name ASC"
	distanceSelect := "NULL::float8"
	if p.Lat != nil && p.Lng != nil {
		argN++
		latArg := argN
		argN++
		lngArg := argN
		args = append(args, *p.Lat, *p.Lng)
		distanceSelect = fmt.Sprintf(
			"ST_Distance(p.location, ST_SetSRID(ST_MakePoint($%d, $%d), 4326)::geography)",
			lngArg, latArg,
		)
		orderBy = "distance_m ASC NULLS LAST, p.name ASC"
	}

	argN++
	limitArg := argN
	args = append(args, p.Limit)
	argN++
	offsetArg := argN
	args = append(args, p.Offset)

	q := fmt.Sprintf(`
		SELECT p.id, p.name, p.address,
		       ST_Y(p.location::geometry), ST_X(p.location::geometry),
		       p.phone, p.category, p.source,
		       %s AS distance_m
		FROM places p
		WHERE p.is_active = true%s
		ORDER BY %s
		LIMIT $%d OFFSET $%d`,
		distanceSelect, where, orderBy, limitArg, offsetArg,
	)

	rows, err := r.pool.Query(ctx, q, args...)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var places []models.PlaceSummary
	var ids []string
	for rows.Next() {
		var pl models.PlaceSummary
		if err := rows.Scan(
			&pl.ID, &pl.Name, &pl.Address,
			&pl.Location.Lat, &pl.Location.Lng,
			&pl.Phone, &pl.Category, &pl.Source,
			&pl.DistanceM,
		); err != nil {
			return nil, err
		}
		places = append(places, pl)
		ids = append(ids, pl.ID)
	}
	if err := rows.Err(); err != nil {
		return nil, err
	}
	if len(ids) == 0 {
		return places, nil
	}

	amenitiesByPlace, err := r.amenitiesForPlaces(ctx, ids, p.UserID)
	if err != nil {
		return nil, err
	}
	for i := range places {
		places[i].Amenities = amenitiesByPlace[places[i].ID]
	}
	return places, nil
}

func placeFilters(cityID string, category *models.PlaceCategory, search string) (string, []any) {
	args := []any{cityID}
	where := " AND p.city_id = $1"
	if category != nil {
		args = append(args, *category)
		where += fmt.Sprintf(" AND p.category = $%d", len(args))
	}
	if search != "" {
		args = append(args, "%"+strings.ToLower(search)+"%")
		where += fmt.Sprintf(" AND lower(p.name) LIKE $%d", len(args))
	}
	return where, args
}

func (r *Repository) GetPlace(ctx context.Context, placeID string, lat, lng *float64, userID string) (models.PlaceDetail, error) {
	distanceSelect := "NULL::float8"
	args := []any{placeID}
	if lat != nil && lng != nil {
		args = append(args, *lng, *lat)
		distanceSelect = "ST_Distance(p.location, ST_SetSRID(ST_MakePoint($2, $3), 4326)::geography)"
	}

	q := fmt.Sprintf(`
		SELECT p.id, p.name, p.address,
		       ST_Y(p.location::geometry), ST_X(p.location::geometry),
		       p.phone, p.website, p.category, p.source,
		       %s, p.created_at, p.updated_at
		FROM places p
		WHERE p.id = $1 AND p.is_active = true`, distanceSelect)

	var d models.PlaceDetail
	err := r.pool.QueryRow(ctx, q, args...).Scan(
		&d.ID, &d.Name, &d.Address,
		&d.Location.Lat, &d.Location.Lng,
		&d.Phone, &d.Website, &d.Category, &d.Source,
		&d.DistanceM, &d.CreatedAt, &d.UpdatedAt,
	)
	if errors.Is(err, pgx.ErrNoRows) {
		return d, apperrors.ErrNotFound
	}
	if err != nil {
		return d, err
	}

	amenities, err := r.amenitiesForPlaces(ctx, []string{placeID}, userID)
	if err != nil {
		return d, err
	}
	d.Amenities = amenities[placeID]
	if d.Amenities == nil {
		d.Amenities = []models.PlaceAmenity{}
	}
	return d, nil
}

func (r *Repository) amenitiesForPlaces(ctx context.Context, placeIDs []string, userID string) (map[string][]models.PlaceAmenity, error) {
	if userID == "" {
		return r.amenitiesForPlacesPublic(ctx, placeIDs)
	}
	return r.amenitiesForPlacesWithUser(ctx, placeIDs, userID)
}

func (r *Repository) amenitiesForPlacesPublic(ctx context.Context, placeIDs []string) (map[string][]models.PlaceAmenity, error) {
	const q = `
		SELECT pa.id, pa.place_id, pa.amenity_type_id, at.label, at.icon, at.sort_order,
		       pa.confirmation_count, pa.is_verified, pa.first_reported_at, pa.verified_at
		FROM place_amenities pa
		JOIN amenity_types at ON at.id = pa.amenity_type_id
		WHERE pa.place_id = ANY($1::uuid[])
		ORDER BY at.sort_order`

	rows, err := r.pool.Query(ctx, q, placeIDs)
	if err != nil {
		return nil, err
	}
	defer rows.Close()
	return scanAmenities(rows, false)
}

func (r *Repository) amenitiesForPlacesWithUser(ctx context.Context, placeIDs []string, userID string) (map[string][]models.PlaceAmenity, error) {
	const q = `
		SELECT pa.id, pa.place_id, pa.amenity_type_id, at.label, at.icon, at.sort_order,
		       pa.confirmation_count, pa.is_verified, pa.first_reported_at, pa.verified_at,
		       EXISTS (
		         SELECT 1 FROM amenity_confirmations ac
		         WHERE ac.place_amenity_id = pa.id AND ac.user_id = $2
		       )
		FROM place_amenities pa
		JOIN amenity_types at ON at.id = pa.amenity_type_id
		WHERE pa.place_id = ANY($1::uuid[])
		ORDER BY at.sort_order`

	rows, err := r.pool.Query(ctx, q, placeIDs, userID)
	if err != nil {
		return nil, err
	}
	defer rows.Close()
	return scanAmenities(rows, true)
}

func scanAmenities(rows pgx.Rows, withUser bool) (map[string][]models.PlaceAmenity, error) {
	out := make(map[string][]models.PlaceAmenity)
	for rows.Next() {
		var pa models.PlaceAmenity
		var placeID string
		var at models.AmenityType
		var userConfirmed bool
		dest := []any{
			&pa.ID, &placeID, &pa.AmenityTypeID, &at.Label, &at.Icon, &at.SortOrder,
			&pa.ConfirmationCount, &pa.IsVerified, &pa.FirstReportedAt, &pa.VerifiedAt,
		}
		if withUser {
			dest = append(dest, &userConfirmed)
		}
		if err := rows.Scan(dest...); err != nil {
			return nil, err
		}
		at.ID = pa.AmenityTypeID
		pa.AmenityType = &at
		if withUser {
			pa.UserHasConfirmed = &userConfirmed
		}
		out[placeID] = append(out[placeID], pa)
	}
	return out, rows.Err()
}

func (r *Repository) AmenityTypeExists(ctx context.Context, id string) (bool, error) {
	var exists bool
	err := r.pool.QueryRow(ctx, `SELECT EXISTS(SELECT 1 FROM amenity_types WHERE id = $1 AND is_active)`, id).Scan(&exists)
	return exists, err
}

func (r *Repository) CreatePlaceAmenity(ctx context.Context, placeID, amenityTypeID, userID string) (models.PlaceAmenity, error) {
	const q = `
		INSERT INTO place_amenities (place_id, amenity_type_id, first_reported_by)
		VALUES ($1, $2, $3)
		RETURNING id, amenity_type_id, confirmation_count, is_verified, first_reported_at, verified_at`

	var pa models.PlaceAmenity
	err := r.pool.QueryRow(ctx, q, placeID, amenityTypeID, userID).Scan(
		&pa.ID, &pa.AmenityTypeID, &pa.ConfirmationCount, &pa.IsVerified,
		&pa.FirstReportedAt, &pa.VerifiedAt,
	)
	if isUniqueViolation(err) {
		existing, gerr := r.GetPlaceAmenityByPlaceAndType(ctx, placeID, amenityTypeID, userID)
		if gerr != nil {
			return pa, gerr
		}
		return existing, apperrors.ErrAmenityAlreadyExists
	}
	if err != nil {
		return pa, err
	}
	return r.enrichPlaceAmenity(ctx, pa, userID)
}

func (r *Repository) GetPlaceAmenityByPlaceAndType(ctx context.Context, placeID, amenityTypeID, userID string) (models.PlaceAmenity, error) {
	const q = `
		SELECT pa.id, pa.amenity_type_id, at.label, at.icon, at.sort_order,
		       pa.confirmation_count, pa.is_verified, pa.first_reported_at, pa.verified_at,
		       pa.first_reported_by
		FROM place_amenities pa
		JOIN amenity_types at ON at.id = pa.amenity_type_id
		WHERE pa.place_id = $1 AND pa.amenity_type_id = $2`

	var pa models.PlaceAmenity
	var at models.AmenityType
	var reporter string
	err := r.pool.QueryRow(ctx, q, placeID, amenityTypeID).Scan(
		&pa.ID, &pa.AmenityTypeID, &at.Label, &at.Icon, &at.SortOrder,
		&pa.ConfirmationCount, &pa.IsVerified, &pa.FirstReportedAt, &pa.VerifiedAt,
		&reporter,
	)
	if errors.Is(err, pgx.ErrNoRows) {
		return pa, apperrors.ErrNotFound
	}
	if err != nil {
		return pa, err
	}
	at.ID = pa.AmenityTypeID
	pa.AmenityType = &at
	if userID != "" {
		confirmed, _ := r.userHasConfirmed(ctx, pa.ID, userID)
		pa.UserHasConfirmed = &confirmed
	}
	_ = reporter
	return pa, nil
}

func (r *Repository) enrichPlaceAmenity(ctx context.Context, pa models.PlaceAmenity, userID string) (models.PlaceAmenity, error) {
	types, err := r.ListAmenityTypes(ctx)
	if err != nil {
		return pa, err
	}
	for _, t := range types {
		if t.ID == pa.AmenityTypeID {
			at := t
			pa.AmenityType = &at
			break
		}
	}
	if userID != "" {
		confirmed, _ := r.userHasConfirmed(ctx, pa.ID, userID)
		pa.UserHasConfirmed = &confirmed
	}
	return pa, nil
}

func (r *Repository) ConfirmPlaceAmenity(ctx context.Context, placeID, amenityTypeID, userID string) (models.PlaceAmenity, error) {
	pa, err := r.GetPlaceAmenityByPlaceAndType(ctx, placeID, amenityTypeID, "")
	if err != nil {
		return pa, err
	}

	const reporterQ = `SELECT first_reported_by FROM place_amenities WHERE id = $1`
	var reporter string
	if err := r.pool.QueryRow(ctx, reporterQ, pa.ID).Scan(&reporter); err != nil {
		return pa, err
	}
	if reporter == userID {
		return pa, apperrors.ErrCannotConfirmOwnReport
	}

	confirmed, err := r.userHasConfirmed(ctx, pa.ID, userID)
	if err != nil {
		return pa, err
	}
	if confirmed {
		return pa, apperrors.ErrAlreadyConfirmed
	}

	_, err = r.pool.Exec(ctx,
		`INSERT INTO amenity_confirmations (place_amenity_id, user_id) VALUES ($1, $2)`,
		pa.ID, userID,
	)
	if isUniqueViolation(err) {
		return pa, apperrors.ErrAlreadyConfirmed
	}
	if err != nil {
		if isRaiseException(err, "Cannot confirm your own report") {
			return pa, apperrors.ErrCannotConfirmOwnReport
		}
		return pa, err
	}

	return r.GetPlaceAmenityByPlaceAndType(ctx, placeID, amenityTypeID, userID)
}

func (r *Repository) userHasConfirmed(ctx context.Context, placeAmenityID, userID string) (bool, error) {
	var exists bool
	err := r.pool.QueryRow(ctx,
		`SELECT EXISTS(SELECT 1 FROM amenity_confirmations WHERE place_amenity_id = $1 AND user_id = $2)`,
		placeAmenityID, userID,
	).Scan(&exists)
	return exists, err
}

func (r *Repository) PlaceExists(ctx context.Context, placeID string) (bool, error) {
	var exists bool
	err := r.pool.QueryRow(ctx,
		`SELECT EXISTS(SELECT 1 FROM places WHERE id = $1 AND is_active)`,
		placeID,
	).Scan(&exists)
	return exists, err
}

func (r *Repository) UpsertImportedPlaces(ctx context.Context, cityID string, places []models.ImportedPlace) (int, error) {
	if len(places) == 0 {
		return 0, nil
	}
	tx, err := r.pool.Begin(ctx)
	if err != nil {
		return 0, err
	}
	defer tx.Rollback(ctx)

	imported := 0
	const q = `
		INSERT INTO places (city_id, google_place_id, name, address, location, phone, website, category, source)
		VALUES ($1, $2, $3, $4, ST_SetSRID(ST_MakePoint($5, $6), 4326)::geography, $7, $8, $9, 'google')
		ON CONFLICT (google_place_id) DO NOTHING`

	for _, p := range places {
		ct, err := tx.Exec(ctx, q, cityID, p.GooglePlaceID, p.Name, p.Address, p.Lng, p.Lat, p.Phone, p.Website, p.Category)
		if err != nil {
			return 0, err
		}
		imported += int(ct.RowsAffected())
	}
	return imported, tx.Commit(ctx)
}

func (r *Repository) GetProfile(ctx context.Context, userID string) (models.Profile, error) {
	const q = `SELECT id, display_name, avatar_url, created_at, updated_at FROM profiles WHERE id = $1`
	var p models.Profile
	err := r.pool.QueryRow(ctx, q, userID).Scan(&p.ID, &p.DisplayName, &p.AvatarURL, &p.CreatedAt, &p.UpdatedAt)
	if errors.Is(err, pgx.ErrNoRows) {
		return p, apperrors.ErrNotFound
	}
	return p, err
}

func (r *Repository) UpdateProfile(ctx context.Context, userID string, displayName *string) (models.Profile, error) {
	const q = `
		UPDATE profiles SET display_name = $2, updated_at = now()
		WHERE id = $1
		RETURNING id, display_name, avatar_url, created_at, updated_at`
	var p models.Profile
	err := r.pool.QueryRow(ctx, q, userID, displayName).Scan(
		&p.ID, &p.DisplayName, &p.AvatarURL, &p.CreatedAt, &p.UpdatedAt,
	)
	if errors.Is(err, pgx.ErrNoRows) {
		return p, apperrors.ErrNotFound
	}
	return p, err
}

func (r *Repository) UpsertDevUser(ctx context.Context, email, displayName string) (string, error) {
	var id string
	err := r.pool.QueryRow(ctx, `SELECT id FROM auth.users WHERE email = $1`, email).Scan(&id)
	if errors.Is(err, pgx.ErrNoRows) {
		err = r.pool.QueryRow(ctx, `
			INSERT INTO auth.users (email, raw_user_meta_data)
			VALUES ($1, jsonb_build_object('display_name', $2::text))
			RETURNING id`, email, displayName).Scan(&id)
		return id, err
	}
	if err != nil {
		return "", err
	}

	_, err = r.pool.Exec(ctx, `
		UPDATE profiles SET display_name = $2, updated_at = now() WHERE id = $1`, id, displayName)
	return id, err
}

func isUniqueViolation(err error) bool {
	var pgErr *pgconn.PgError
	return errors.As(err, &pgErr) && pgErr.Code == "23505"
}

func isRaiseException(err error, msg string) bool {
	var pgErr *pgconn.PgError
	return errors.As(err, &pgErr) && pgErr.Code == "P0001" && strings.Contains(pgErr.Message, msg)
}
