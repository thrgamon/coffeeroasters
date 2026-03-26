package api

import (
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/jackc/pgx/v5/pgtype"

	"github.com/thrgamon/coffeeroasters/internal/db"
	"github.com/thrgamon/coffeeroasters/internal/domain"
)

// ListCafes godoc
// @Summary List all active cafes
// @Tags cafes
// @Produce json
// @Param state query string false "Filter by AU state code (e.g. VIC, NSW)"
// @Success 200 {object} domain.CafeListResponse
// @Router /api/cafes [get]
func (h *Handler) ListCafes(c *gin.Context) {
	ctx := c.Request.Context()
	state := c.Query("state")

	var resp domain.CafeListResponse

	if state != "" {
		rows, err := h.queries.ListCafesByState(ctx, pgtype.Text{String: state, Valid: true})
		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": "failed to list cafes"})
			return
		}
		for _, r := range rows {
			resp.Cafes = append(resp.Cafes, cafesByStateRowToResponse(r))
		}
	} else {
		rows, err := h.queries.ListCafes(ctx)
		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": "failed to list cafes"})
			return
		}
		for _, r := range rows {
			resp.Cafes = append(resp.Cafes, cafesRowToResponse(r))
		}
	}

	if resp.Cafes == nil {
		resp.Cafes = []domain.CafeResponse{}
	}

	c.JSON(http.StatusOK, resp)
}

// GetCafe godoc
// @Summary Get a cafe by slug
// @Tags cafes
// @Produce json
// @Param slug path string true "Cafe slug"
// @Success 200 {object} domain.CafeDetailResponse
// @Failure 404 {object} map[string]string
// @Router /api/cafes/{slug} [get]
func (h *Handler) GetCafe(c *gin.Context) {
	ctx := c.Request.Context()
	slug := c.Param("slug")

	row, err := h.queries.GetCafeBySlug(ctx, slug)
	if err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "cafe not found"})
		return
	}

	resp := domain.CafeDetailResponse{
		CafeResponse: cafeDetailRowToResponse(row),
	}

	c.JSON(http.StatusOK, resp)
}

func cafesRowToResponse(row db.ListCafesRow) domain.CafeResponse {
	resp := domain.CafeResponse{
		ID:          row.ID,
		RoasterID:   row.RoasterID,
		RoasterName: row.RoasterName,
		RoasterSlug: row.RoasterSlug,
		Slug:        row.Slug,
		Name:        row.Name,
		Type:        row.Type,
		Address:     row.Address.String,
		Suburb:      row.Suburb.String,
		State:       row.State.String,
		Postcode:    row.Postcode.String,
		Phone:       row.Phone.String,
		Instagram:   row.Instagram.String,
		WebsiteURL:  row.WebsiteUrl.String,
		ImageURL:    row.ImageUrl.String,
	}
	if row.Latitude.Valid {
		resp.Latitude = &row.Latitude.Float64
	}
	if row.Longitude.Valid {
		resp.Longitude = &row.Longitude.Float64
	}
	return resp
}

func cafesByStateRowToResponse(row db.ListCafesByStateRow) domain.CafeResponse {
	resp := domain.CafeResponse{
		ID:          row.ID,
		RoasterID:   row.RoasterID,
		RoasterName: row.RoasterName,
		RoasterSlug: row.RoasterSlug,
		Slug:        row.Slug,
		Name:        row.Name,
		Type:        row.Type,
		Address:     row.Address.String,
		Suburb:      row.Suburb.String,
		State:       row.State.String,
		Postcode:    row.Postcode.String,
		Phone:       row.Phone.String,
		Instagram:   row.Instagram.String,
		WebsiteURL:  row.WebsiteUrl.String,
		ImageURL:    row.ImageUrl.String,
	}
	if row.Latitude.Valid {
		resp.Latitude = &row.Latitude.Float64
	}
	if row.Longitude.Valid {
		resp.Longitude = &row.Longitude.Float64
	}
	return resp
}

func cafeDetailRowToResponse(row db.GetCafeBySlugRow) domain.CafeResponse {
	resp := domain.CafeResponse{
		ID:          row.ID,
		RoasterID:   row.RoasterID,
		RoasterName: row.RoasterName,
		RoasterSlug: row.RoasterSlug,
		Slug:        row.Slug,
		Name:        row.Name,
		Type:        row.Type,
		Address:     row.Address.String,
		Suburb:      row.Suburb.String,
		State:       row.State.String,
		Postcode:    row.Postcode.String,
		Phone:       row.Phone.String,
		Instagram:   row.Instagram.String,
		WebsiteURL:  row.WebsiteUrl.String,
		ImageURL:    row.ImageUrl.String,
	}
	if row.Latitude.Valid {
		resp.Latitude = &row.Latitude.Float64
	}
	if row.Longitude.Valid {
		resp.Longitude = &row.Longitude.Float64
	}
	return resp
}

func cafeByRoasterRowToResponse(row db.ListCafesByRoasterRow) domain.CafeResponse {
	resp := domain.CafeResponse{
		ID:       row.ID,
		Slug:     row.Slug,
		Name:     row.Name,
		Type:     row.Type,
		Address:  row.Address.String,
		Suburb:   row.Suburb.String,
		State:    row.State.String,
		Postcode: row.Postcode.String,
		Phone:    row.Phone.String,
		Instagram: row.Instagram.String,
		WebsiteURL: row.WebsiteUrl.String,
		ImageURL:   row.ImageUrl.String,
	}
	if row.Latitude.Valid {
		resp.Latitude = &row.Latitude.Float64
	}
	if row.Longitude.Valid {
		resp.Longitude = &row.Longitude.Float64
	}
	return resp
}
