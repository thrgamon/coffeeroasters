package api

import (
	"net/http"
	"strconv"

	"github.com/gin-gonic/gin"
	"github.com/jackc/pgx/v5/pgtype"

	"github.com/thrgamon/coffeeroasters/internal/db"
	"github.com/thrgamon/coffeeroasters/internal/domain"
)

// ListCountries godoc
// @Summary List countries with coffee counts
// @Tags countries
// @Produce json
// @Success 200 {object} domain.CountryListResponse
// @Router /api/countries [get]
func (h *Handler) ListCountries(c *gin.Context) {
	rows, err := h.queries.ListCountriesWithCoffeeCount(c.Request.Context())
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "failed to list countries"})
		return
	}

	resp := domain.CountryListResponse{
		Countries: make([]domain.CountryResponse, 0, len(rows)),
	}
	for _, r := range rows {
		resp.Countries = append(resp.Countries, domain.CountryResponse{
			Code:        r.Code,
			Name:        r.Name,
			CoffeeCount: r.CoffeeCount,
		})
	}

	c.JSON(http.StatusOK, resp)
}

// GetCountry godoc
// @Summary Get a country with its regions and coffees
// @Tags countries
// @Produce json
// @Param code path string true "ISO 3166-1 alpha-2 country code"
// @Success 200 {object} domain.CountryDetailResponse
// @Failure 404 {object} map[string]string
// @Router /api/countries/{code} [get]
func (h *Handler) GetCountry(c *gin.Context) {
	ctx := c.Request.Context()
	code := c.Param("code")

	country, err := h.queries.GetCountryByCode(ctx, code)
	if err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "country not found"})
		return
	}

	regionRows, err := h.queries.ListRegionsByCountry(ctx, code)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "failed to list regions"})
		return
	}

	regions := make([]domain.RegionSummary, 0, len(regionRows))
	for _, r := range regionRows {
		regions = append(regions, domain.RegionSummary{
			ID:          r.ID,
			Name:        r.Name,
			CoffeeCount: r.CoffeeCount,
		})
	}

	coffeeRows, err := h.queries.ListCoffeesByCountry(ctx, pgtype.Text{String: code, Valid: true})
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "failed to list coffees"})
		return
	}

	coffees := make([]domain.CoffeeResponse, 0, len(coffeeRows))
	for _, row := range coffeeRows {
		coffees = append(coffees, countryListRowToResponse(row))
	}

	c.JSON(http.StatusOK, domain.CountryDetailResponse{
		Code:    country.Code,
		Name:    country.Name,
		Regions: regions,
		Coffees: coffees,
	})
}

// GetRegion godoc
// @Summary Get a region with its coffees
// @Tags regions
// @Produce json
// @Param id path int true "Region ID"
// @Success 200 {object} domain.RegionDetailResponse
// @Failure 404 {object} map[string]string
// @Router /api/regions/{id} [get]
func (h *Handler) GetRegion(c *gin.Context) {
	id, err := strconv.ParseInt(c.Param("id"), 10, 32)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid region ID"})
		return
	}

	ctx := c.Request.Context()

	region, err := h.queries.GetRegionByID(ctx, int32(id))
	if err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "region not found"})
		return
	}

	coffeeRows, err := h.queries.ListCoffeesByRegion(ctx, pgtype.Int4{Int32: int32(id), Valid: true})
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "failed to list coffees"})
		return
	}

	coffees := make([]domain.CoffeeResponse, 0, len(coffeeRows))
	for _, row := range coffeeRows {
		coffees = append(coffees, regionListRowToResponse(row))
	}

	resp := domain.RegionDetailResponse{
		ID:          region.ID,
		Name:        region.Name,
		CountryCode: region.CountryCode,
		CountryName: region.CountryName,
		Coffees:     coffees,
	}

	if region.Latitude.Valid {
		resp.Latitude = &region.Latitude.Float64
	}
	if region.Longitude.Valid {
		resp.Longitude = &region.Longitude.Float64
	}

	if region.Latitude.Valid && region.Longitude.Valid {
		nearbyRows, err := h.queries.ListNearbyRegions(ctx, db.ListNearbyRegionsParams{
			SourceLat:       region.Latitude.Float64,
			SourceLon:       region.Longitude.Float64,
			ExcludeRegionID: region.ID,
			MaxDistanceKm:   500,
		})
		if err == nil {
			nearby := make([]domain.NearbyRegion, 0, len(nearbyRows))
			for _, nr := range nearbyRows {
				nearby = append(nearby, domain.NearbyRegion{
					ID:          nr.ID,
					Name:        nr.Name,
					CountryCode: nr.CountryCode,
					CountryName: nr.CountryName,
					DistanceKm:  nr.DistanceKm,
					CoffeeCount: nr.CoffeeCount,
				})
			}
			resp.NearbyRegions = nearby
		}
	}

	c.JSON(http.StatusOK, resp)
}

// GetProducer godoc
// @Summary Get a producer with their coffees
// @Tags producers
// @Produce json
// @Param id path int true "Producer ID"
// @Success 200 {object} domain.ProducerDetailResponse
// @Failure 404 {object} map[string]string
// @Router /api/producers/{id} [get]
func (h *Handler) GetProducer(c *gin.Context) {
	id, err := strconv.ParseInt(c.Param("id"), 10, 32)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid producer ID"})
		return
	}

	ctx := c.Request.Context()

	producer, err := h.queries.GetProducerByID(ctx, int32(id))
	if err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "producer not found"})
		return
	}

	coffeeRows, err := h.queries.ListCoffeesByProducer(ctx, pgtype.Int4{Int32: int32(id), Valid: true})
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "failed to list coffees"})
		return
	}

	coffees := make([]domain.CoffeeResponse, 0, len(coffeeRows))
	for _, row := range coffeeRows {
		coffees = append(coffees, producerListRowToResponse(row))
	}

	c.JSON(http.StatusOK, domain.ProducerDetailResponse{
		ID:          producer.ID,
		Name:        producer.Name,
		CountryCode: producer.CountryCode.String,
		CountryName: producer.CountryName.String,
		RegionName:  producer.RegionName.String,
		Coffees:     coffees,
	})
}
