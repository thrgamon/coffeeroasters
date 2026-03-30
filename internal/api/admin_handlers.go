package api

import (
	"net/http"
	"strconv"

	"github.com/gin-gonic/gin"
	"github.com/jackc/pgx/v5/pgtype"

	"github.com/thrgamon/coffeeroasters/internal/db"
	"github.com/thrgamon/coffeeroasters/internal/domain"
)

// ---- Roasters ----

func (h *Handler) AdminListRoasters(c *gin.Context) {
	rows, err := h.queries.AdminListRoasters(c.Request.Context())
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "failed to list roasters"})
		return
	}

	roasters := make([]domain.AdminRoasterResponse, 0, len(rows))
	for _, r := range rows {
		roasters = append(roasters, domain.AdminRoasterResponse{
			ID:          r.ID,
			Slug:        r.Slug,
			Name:        r.Name,
			Website:     r.Website,
			State:       r.State.String,
			Description: r.Description.String,
			LogoURL:     r.LogoUrl.String,
			Active:      r.Active,
			OptedOut:    r.OptedOut,
			CreatedAt:   r.CreatedAt.Format("2006-01-02T15:04:05Z"),
			UpdatedAt:   r.UpdatedAt.Format("2006-01-02T15:04:05Z"),
		})
	}

	c.JSON(http.StatusOK, domain.AdminRoasterListResponse{Roasters: roasters})
}

func (h *Handler) AdminGetRoaster(c *gin.Context) {
	id, err := strconv.ParseInt(c.Param("id"), 10, 32)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid roaster ID"})
		return
	}

	r, err := h.queries.AdminGetRoaster(c.Request.Context(), int32(id))
	if err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "roaster not found"})
		return
	}

	c.JSON(http.StatusOK, domain.AdminRoasterResponse{
		ID:          r.ID,
		Slug:        r.Slug,
		Name:        r.Name,
		Website:     r.Website,
		State:       r.State.String,
		Description: r.Description.String,
		LogoURL:     r.LogoUrl.String,
		Active:      r.Active,
		OptedOut:    r.OptedOut,
		CreatedAt:   r.CreatedAt.Format("2006-01-02T15:04:05Z"),
		UpdatedAt:   r.UpdatedAt.Format("2006-01-02T15:04:05Z"),
	})
}

func (h *Handler) AdminCreateRoaster(c *gin.Context) {
	var req domain.AdminRoasterRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	id, err := h.queries.AdminCreateRoaster(c.Request.Context(), db.AdminCreateRoasterParams{
		Slug:        req.Slug,
		Name:        req.Name,
		Website:     req.Website,
		State:       textPtr(req.State),
		Description: textPtr(req.Description),
		LogoUrl:     textPtr(req.LogoURL),
		Active:      req.Active,
	})
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "failed to create roaster"})
		return
	}

	c.JSON(http.StatusCreated, gin.H{"id": id})
}

func (h *Handler) AdminUpdateRoaster(c *gin.Context) {
	id, err := strconv.ParseInt(c.Param("id"), 10, 32)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid roaster ID"})
		return
	}

	var req domain.AdminRoasterRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	if err := h.queries.AdminUpdateRoaster(c.Request.Context(), db.AdminUpdateRoasterParams{
		ID:          int32(id),
		Slug:        req.Slug,
		Name:        req.Name,
		Website:     req.Website,
		State:       textPtr(req.State),
		Description: textPtr(req.Description),
		LogoUrl:     textPtr(req.LogoURL),
		Active:      req.Active,
		OptedOut:    req.OptedOut,
	}); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "failed to update roaster"})
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "updated"})
}

// ---- Coffees ----

func (h *Handler) AdminListCoffees(c *gin.Context) {
	page := intQuery(c, "page", 1)
	pageSize := intQuery(c, "page_size", 50)
	if pageSize > 100 {
		pageSize = 100
	}

	rows, err := h.queries.AdminListCoffees(c.Request.Context(), db.AdminListCoffeesParams{
		Lim: int32(pageSize),
		Off: int32((page - 1) * pageSize),
	})
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "failed to list coffees"})
		return
	}

	coffees := make([]domain.AdminCoffeeResponse, 0, len(rows))
	for _, r := range rows {
		coffees = append(coffees, domain.AdminCoffeeResponse{
			ID:              r.ID,
			RoasterID:       r.RoasterID,
			RoasterName:     r.RoasterName,
			RoasterSlug:     r.RoasterSlug,
			Name:            r.Name,
			ProductURL:      r.ProductUrl.String,
			ImageURL:        r.ImageUrl.String,
			CountryCode:     r.CountryCode.String,
			CountryName:     r.CountryName.String,
			Process:         r.Process.String,
			RoastLevel:      r.RoastLevel.String,
			TastingNotes:    r.TastingNotes,
			PriceCents:      r.PriceCents.Int32,
			WeightGrams:     r.WeightGrams.Int32,
			PricePer100gMin: r.PricePer100gMin.Int32,
			PricePer100gMax: r.PricePer100gMax.Int32,
			Variety:         r.Variety.String,
			Species:         r.Species.String,
			IsBlend:         r.IsBlend,
			IsDecaf:         r.IsDecaf,
			InStock:         r.InStock,
			Description:     r.Description.String,
			FirstSeenAt:     r.FirstSeenAt.Format("2006-01-02T15:04:05Z"),
			LastSeenAt:      r.LastSeenAt.Format("2006-01-02T15:04:05Z"),
		})
	}

	var totalCount int64
	if len(rows) > 0 {
		totalCount = rows[0].TotalCount
	}

	c.JSON(http.StatusOK, domain.AdminCoffeeListResponse{
		Coffees:    coffees,
		TotalCount: totalCount,
		Page:       int32(page),
		PageSize:   int32(pageSize),
	})
}

func (h *Handler) AdminGetCoffee(c *gin.Context) {
	id, err := strconv.ParseInt(c.Param("id"), 10, 64)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid coffee ID"})
		return
	}

	row, err := h.queries.GetCoffeeByID(c.Request.Context(), id)
	if err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "coffee not found"})
		return
	}

	c.JSON(http.StatusOK, coffeeDetailRowToResponse(row))
}

func (h *Handler) AdminCreateCoffee(c *gin.Context) {
	var req domain.AdminCoffeeRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	id, err := h.queries.AdminCreateCoffee(c.Request.Context(), db.AdminCreateCoffeeParams{
		RoasterID:       req.RoasterID,
		Name:            req.Name,
		ProductUrl:      textPtr(req.ProductURL),
		ImageUrl:        textPtr(req.ImageURL),
		CountryCode:     textPtr(req.CountryCode),
		RegionID:        int4Ptr(req.RegionID),
		ProducerID:      int4Ptr(req.ProducerID),
		Process:         textPtr(req.Process),
		RoastLevel:      textPtr(req.RoastLevel),
		TastingNotes:    req.TastingNotes,
		PriceCents:      int4Ptr(req.PriceCents),
		WeightGrams:     int4Ptr(req.WeightGrams),
		PricePer100gMin: int4Ptr(req.PricePer100gMin),
		PricePer100gMax: int4Ptr(req.PricePer100gMax),
		Variety:         textPtr(req.Variety),
		Species:         textPtr(req.Species),
		IsBlend:         req.IsBlend,
		IsDecaf:         req.IsDecaf,
		InStock:         req.InStock,
		Description:     textPtr(req.Description),
	})
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "failed to create coffee"})
		return
	}

	c.JSON(http.StatusCreated, gin.H{"id": id})
}

func (h *Handler) AdminUpdateCoffee(c *gin.Context) {
	id, err := strconv.ParseInt(c.Param("id"), 10, 64)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid coffee ID"})
		return
	}

	var req domain.AdminCoffeeRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	if err := h.queries.AdminUpdateCoffee(c.Request.Context(), db.AdminUpdateCoffeeParams{
		ID:              id,
		RoasterID:       req.RoasterID,
		Name:            req.Name,
		ProductUrl:      textPtr(req.ProductURL),
		ImageUrl:        textPtr(req.ImageURL),
		CountryCode:     textPtr(req.CountryCode),
		RegionID:        int4Ptr(req.RegionID),
		ProducerID:      int4Ptr(req.ProducerID),
		Process:         textPtr(req.Process),
		RoastLevel:      textPtr(req.RoastLevel),
		TastingNotes:    req.TastingNotes,
		PriceCents:      int4Ptr(req.PriceCents),
		WeightGrams:     int4Ptr(req.WeightGrams),
		PricePer100gMin: int4Ptr(req.PricePer100gMin),
		PricePer100gMax: int4Ptr(req.PricePer100gMax),
		Variety:         textPtr(req.Variety),
		Species:         textPtr(req.Species),
		IsBlend:         req.IsBlend,
		IsDecaf:         req.IsDecaf,
		InStock:         req.InStock,
		Description:     textPtr(req.Description),
	}); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "failed to update coffee"})
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "updated"})
}

// ---- Cafes ----

func (h *Handler) AdminListCafes(c *gin.Context) {
	rows, err := h.queries.AdminListCafes(c.Request.Context())
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "failed to list cafes"})
		return
	}

	cafes := make([]domain.AdminCafeResponse, 0, len(rows))
	for _, r := range rows {
		cafe := domain.AdminCafeResponse{
			ID:          r.ID,
			RoasterID:   r.RoasterID,
			RoasterName: r.RoasterName,
			RoasterSlug: r.RoasterSlug,
			Slug:        r.Slug,
			Name:        r.Name,
			Type:        r.Type,
			Address:     r.Address.String,
			Suburb:      r.Suburb.String,
			State:       r.State.String,
			Postcode:    r.Postcode.String,
			Phone:       r.Phone.String,
			Instagram:   r.Instagram.String,
			WebsiteURL:  r.WebsiteUrl.String,
			ImageURL:    r.ImageUrl.String,
			Active:      r.Active.Bool,
		}
		if r.Latitude.Valid {
			cafe.Latitude = &r.Latitude.Float64
		}
		if r.Longitude.Valid {
			cafe.Longitude = &r.Longitude.Float64
		}
		cafes = append(cafes, cafe)
	}

	c.JSON(http.StatusOK, domain.AdminCafeListResponse{Cafes: cafes})
}

func (h *Handler) AdminGetCafe(c *gin.Context) {
	id, err := strconv.ParseInt(c.Param("id"), 10, 32)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid cafe ID"})
		return
	}

	r, err := h.queries.AdminGetCafe(c.Request.Context(), int32(id))
	if err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "cafe not found"})
		return
	}

	cafe := domain.AdminCafeResponse{
		ID:          r.ID,
		RoasterID:   r.RoasterID,
		RoasterName: r.RoasterName,
		RoasterSlug: r.RoasterSlug,
		Slug:        r.Slug,
		Name:        r.Name,
		Type:        r.Type,
		Address:     r.Address.String,
		Suburb:      r.Suburb.String,
		State:       r.State.String,
		Postcode:    r.Postcode.String,
		Phone:       r.Phone.String,
		Instagram:   r.Instagram.String,
		WebsiteURL:  r.WebsiteUrl.String,
		ImageURL:    r.ImageUrl.String,
		Active:      r.Active.Bool,
	}
	if r.Latitude.Valid {
		cafe.Latitude = &r.Latitude.Float64
	}
	if r.Longitude.Valid {
		cafe.Longitude = &r.Longitude.Float64
	}

	c.JSON(http.StatusOK, cafe)
}

func (h *Handler) AdminCreateCafe(c *gin.Context) {
	var req domain.AdminCafeRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	var lat pgtype.Float8
	if req.Latitude != nil {
		lat = pgtype.Float8{Float64: *req.Latitude, Valid: true}
	}
	var lon pgtype.Float8
	if req.Longitude != nil {
		lon = pgtype.Float8{Float64: *req.Longitude, Valid: true}
	}

	id, err := h.queries.AdminCreateCafe(c.Request.Context(), db.AdminCreateCafeParams{
		RoasterID:  req.RoasterID,
		Slug:       req.Slug,
		Name:       req.Name,
		Type:       req.Type,
		Address:    textPtr(req.Address),
		Suburb:     textPtr(req.Suburb),
		State:      textPtr(req.State),
		Postcode:   textPtr(req.Postcode),
		Latitude:   lat,
		Longitude:  lon,
		Phone:      textPtr(req.Phone),
		Instagram:  textPtr(req.Instagram),
		WebsiteUrl: textPtr(req.WebsiteURL),
		ImageUrl:   textPtr(req.ImageURL),
	})
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "failed to create cafe"})
		return
	}

	c.JSON(http.StatusCreated, gin.H{"id": id})
}

func (h *Handler) AdminUpdateCafe(c *gin.Context) {
	id, err := strconv.ParseInt(c.Param("id"), 10, 32)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid cafe ID"})
		return
	}

	var req domain.AdminCafeRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	var lat pgtype.Float8
	if req.Latitude != nil {
		lat = pgtype.Float8{Float64: *req.Latitude, Valid: true}
	}
	var lon pgtype.Float8
	if req.Longitude != nil {
		lon = pgtype.Float8{Float64: *req.Longitude, Valid: true}
	}

	if err := h.queries.AdminUpdateCafe(c.Request.Context(), db.AdminUpdateCafeParams{
		ID:         int32(id),
		RoasterID:  req.RoasterID,
		Slug:       req.Slug,
		Name:       req.Name,
		Type:       req.Type,
		Address:    textPtr(req.Address),
		Suburb:     textPtr(req.Suburb),
		State:      textPtr(req.State),
		Postcode:   textPtr(req.Postcode),
		Latitude:   lat,
		Longitude:  lon,
		Phone:      textPtr(req.Phone),
		Instagram:  textPtr(req.Instagram),
		WebsiteUrl: textPtr(req.WebsiteURL),
		ImageUrl:   textPtr(req.ImageURL),
		Active:     pgtype.Bool{Bool: req.Active, Valid: true},
	}); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "failed to update cafe"})
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "updated"})
}

// ---- Scrape Runs ----

func (h *Handler) AdminListScrapeRuns(c *gin.Context) {
	page := intQuery(c, "page", 1)
	pageSize := intQuery(c, "page_size", 50)
	if pageSize > 100 {
		pageSize = 100
	}

	rows, err := h.queries.AdminListScrapeRuns(c.Request.Context(), db.AdminListScrapeRunsParams{
		Lim: int32(pageSize),
		Off: int32((page - 1) * pageSize),
	})
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "failed to list scrape runs"})
		return
	}

	runs := make([]domain.AdminScrapeRunResponse, 0, len(rows))
	for _, r := range rows {
		run := domain.AdminScrapeRunResponse{
			ID:             r.ID,
			RoasterID:      r.RoasterID,
			RoasterName:    r.RoasterName,
			RoasterSlug:    r.RoasterSlug,
			StartedAt:      r.StartedAt.Format("2006-01-02T15:04:05Z"),
			Status:         r.Status,
			CoffeesFound:   r.CoffeesFound.Int32,
			CoffeesAdded:   r.CoffeesAdded.Int32,
			CoffeesUpdated: r.CoffeesUpdated.Int32,
			CoffeesRemoved: r.CoffeesRemoved.Int32,
			ErrorMessage:   r.ErrorMessage.String,
			DurationMs:     r.DurationMs.Int32,
		}
		if r.FinishedAt.Valid {
			run.FinishedAt = r.FinishedAt.Time.Format("2006-01-02T15:04:05Z")
		}
		runs = append(runs, run)
	}

	c.JSON(http.StatusOK, domain.AdminScrapeRunListResponse{Runs: runs})
}

// int4Ptr converts an int32 to pgtype.Int4, treating 0 as null.
func int4Ptr(v int32) pgtype.Int4 {
	if v == 0 {
		return pgtype.Int4{}
	}
	return pgtype.Int4{Int32: v, Valid: true}
}
