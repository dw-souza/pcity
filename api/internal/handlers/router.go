package handlers

import (
	"net/http"

	"github.com/dw-souza/pcity/api/internal/middleware"
	"github.com/go-chi/chi/v5"
	chimw "github.com/go-chi/chi/v5/middleware"
)

type API struct {
	jwtSecret string
	devAuth   *DevAuthHandler

	health    *HealthHandler
	cities    *CitiesHandler
	places    *PlacesHandler
	amenities *AmenitiesHandler
	profile   *ProfileHandler
	admin     *AdminHandler
}

func NewAPI(
	jwtSecret string,
	devAuth *DevAuthHandler,
	health *HealthHandler,
	cities *CitiesHandler,
	places *PlacesHandler,
	amenities *AmenitiesHandler,
	profile *ProfileHandler,
	admin *AdminHandler,
) *API {
	return &API{
		jwtSecret: jwtSecret,
		devAuth:   devAuth,
		health:    health,
		cities:    cities,
		places:    places,
		amenities: amenities,
		profile:   profile,
		admin:     admin,
	}
}

func (a *API) Router(adminKey string, devAuthEnabled bool) http.Handler {
	r := chi.NewRouter()
	r.Use(chimw.Recoverer)
	r.Use(middleware.RequestLogger)

	r.Route("/v1", func(r chi.Router) {
		r.Get("/health", a.health.Get)
		r.Get("/amenity-types", a.amenities.ListTypes)

		if devAuthEnabled && a.devAuth != nil {
			r.Post("/dev/auth", a.devAuth.Register)
		}

		r.Get("/cities/{slug}", a.cities.Get)
		r.With(middleware.OptionalAuth(a.jwtSecret)).Get("/cities/{slug}/places", a.cities.ListPlaces)

		r.Route("/places/{placeId}", func(r chi.Router) {
			r.With(middleware.OptionalAuth(a.jwtSecret)).Get("/", a.places.Get)

			r.Group(func(r chi.Router) {
				r.Use(middleware.RequireAuth(a.jwtSecret))
				r.Post("/amenities", a.amenities.Report)
				r.Post("/amenities/{amenityTypeId}/confirm", a.amenities.Confirm)
			})
		})

		r.Route("/me", func(r chi.Router) {
			r.Use(middleware.RequireAuth(a.jwtSecret))
			r.Get("/", a.profile.Get)
			r.Patch("/", a.profile.Patch)
		})

		r.Route("/admin", func(r chi.Router) {
			r.Use(middleware.RequireAdmin(adminKey))
			r.Post("/cities/{slug}/index", a.admin.IndexCity)
		})
	})

	return r
}
