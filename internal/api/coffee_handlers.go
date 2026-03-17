package api

import (
	"log/slog"
	"net/http"
	"strconv"

	"github.com/gin-gonic/gin"
	"github.com/jackc/pgx/v5/pgtype"

	"github.com/thrgamon/coffeeroasters/internal/db"
	"github.com/thrgamon/coffeeroasters/internal/domain"
	"github.com/thrgamon/coffeeroasters/internal/similarity"
)

// ListRoasters godoc
// @Summary List all active roasters
// @Tags roasters
// @Produce json
// @Param state query string false "Filter by AU state code (e.g. VIC, NSW)"
// @Success 200 {object} domain.RoasterListResponse
// @Router /api/roasters [get]
func (h *Handler) ListRoasters(c *gin.Context) {
	ctx := c.Request.Context()
	state := c.Query("state")

	var resp domain.RoasterListResponse

	if state != "" {
		rows, err := h.queries.ListRoastersByState(ctx, pgtype.Text{String: state, Valid: true})
		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": "failed to list roasters"})
			return
		}
		for _, r := range rows {
			resp.Roasters = append(resp.Roasters, domain.RoasterResponse{
				ID:      r.ID,
				Slug:    r.Slug,
				Name:    r.Name,
				Website: r.Website,
				State:   r.State.String,
			})
		}
	} else {
		rows, err := h.queries.ListRoasters(ctx)
		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": "failed to list roasters"})
			return
		}
		for _, r := range rows {
			resp.Roasters = append(resp.Roasters, domain.RoasterResponse{
				ID:      r.ID,
				Slug:    r.Slug,
				Name:    r.Name,
				Website: r.Website,
				State:   r.State.String,
			})
		}
	}

	if resp.Roasters == nil {
		resp.Roasters = []domain.RoasterResponse{}
	}

	c.JSON(http.StatusOK, resp)
}

// GetRoaster godoc
// @Summary Get a roaster by slug with their coffees
// @Tags roasters
// @Produce json
// @Param slug path string true "Roaster slug"
// @Success 200 {object} domain.RoasterDetailResponse
// @Failure 404 {object} map[string]string
// @Router /api/roasters/{slug} [get]
func (h *Handler) GetRoaster(c *gin.Context) {
	ctx := c.Request.Context()
	slug := c.Param("slug")

	roaster, err := h.queries.GetRoasterBySlug(ctx, slug)
	if err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "roaster not found"})
		return
	}

	coffeeRows, err := h.queries.ListCoffeesByRoaster(ctx, slug)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "failed to list coffees"})
		return
	}

	coffees := make([]domain.CoffeeResponse, 0, len(coffeeRows))
	for _, row := range coffeeRows {
		coffees = append(coffees, coffeeRowToResponse(row))
	}

	c.JSON(http.StatusOK, domain.RoasterDetailResponse{
		Roaster: domain.RoasterResponse{
			ID:      roaster.ID,
			Slug:    roaster.Slug,
			Name:    roaster.Name,
			Website: roaster.Website,
			State:   roaster.State.String,
		},
		Coffees: coffees,
	})
}

// ListCoffees godoc
// @Summary List or search coffees with optional filters
// @Tags coffees
// @Produce json
// @Param q query string false "Full-text search query"
// @Param origin query string false "Filter by origin"
// @Param process query string false "Filter by process"
// @Param roast query string false "Filter by roast level"
// @Param variety query string false "Filter by variety"
// @Param in_stock query boolean false "Filter by stock status"
// @Param page query int false "Page number (default 1)"
// @Param page_size query int false "Page size (default 20, max 100)"
// @Success 200 {object} domain.CoffeeListResponse
// @Router /api/coffees [get]
func (h *Handler) ListCoffees(c *gin.Context) {
	ctx := c.Request.Context()

	page := intQuery(c, "page", 1)
	pageSize := intQuery(c, "page_size", 20)
	if pageSize > 100 {
		pageSize = 100
	}
	offset := (page - 1) * pageSize

	q := c.Query("q")
	origin := c.Query("origin")
	process := c.Query("process")
	roast := c.Query("roast")
	varietyFilter := c.Query("variety")
	inStockStr := c.Query("in_stock")

	var resp domain.CoffeeListResponse
	resp.Page = int32(page)
	resp.PageSize = int32(pageSize)

	if q != "" {
		// Full-text search
		rows, err := h.queries.SearchCoffees(ctx, db.SearchCoffeesParams{
			PlaintoTsquery: q,
			Limit:          int32(pageSize),
			Offset:         int32(offset),
		})
		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": "search failed"})
			return
		}
		for _, row := range rows {
			resp.Coffees = append(resp.Coffees, searchRowToResponse(row))
		}
		resp.TotalCount = int64(len(rows)) // approximate for search
	} else if origin != "" || process != "" || roast != "" || varietyFilter != "" || inStockStr != "" {
		// Filter
		var inStock pgtype.Bool
		if inStockStr != "" {
			v, _ := strconv.ParseBool(inStockStr)
			inStock = pgtype.Bool{Bool: v, Valid: true}
		}

		rows, err := h.queries.FilterCoffees(ctx, db.FilterCoffeesParams{
			Column1: origin,
			Column2: process,
			Column3: roast,
			Column4: inStock.Bool,
			Column5: varietyFilter,
			Limit:   int32(pageSize),
			Offset:  int32(offset),
		})
		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": "filter failed"})
			return
		}
		for _, row := range rows {
			resp.Coffees = append(resp.Coffees, filterRowToResponse(row))
		}
		resp.TotalCount = int64(len(rows))
	} else {
		// List all
		rows, err := h.queries.ListCoffees(ctx, db.ListCoffeesParams{
			Limit:  int32(pageSize),
			Offset: int32(offset),
		})
		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": "failed to list coffees"})
			return
		}
		for _, row := range rows {
			resp.Coffees = append(resp.Coffees, listRowToResponse(row))
		}
		total, _ := h.queries.CountCoffees(ctx)
		resp.TotalCount = total
	}

	if resp.Coffees == nil {
		resp.Coffees = []domain.CoffeeResponse{}
	}

	c.JSON(http.StatusOK, resp)
}

// GetCoffee godoc
// @Summary Get a single coffee by ID with similar coffees
// @Tags coffees
// @Produce json
// @Param id path int true "Coffee ID"
// @Success 200 {object} domain.CoffeeDetailResponse
// @Failure 404 {object} map[string]string
// @Router /api/coffees/{id} [get]
func (h *Handler) GetCoffee(c *gin.Context) {
	ctx := c.Request.Context()

	id, err := strconv.ParseInt(c.Param("id"), 10, 64)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid coffee ID"})
		return
	}

	row, err := h.queries.GetCoffeeByID(ctx, id)
	if err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "coffee not found"})
		return
	}

	coffeeResp := coffeeDetailRowToResponse(row)
	resp := domain.CoffeeDetailResponse{CoffeeResponse: coffeeResp}

	// Build similar coffees
	simRows, err := h.queries.ListCoffeesForSimilarity(ctx)
	if err != nil {
		slog.Warn("list coffees for similarity", "error", err)
		c.JSON(http.StatusOK, resp)
		return
	}

	source := similarity.CoffeeAttrs{
		CoffeeID:     row.ID,
		TastingNotes: row.TastingNotes,
		Process:      row.Process.String,
		RoastLevel:   row.RoastLevel.String,
		Variety:      row.Variety.String,
	}
	if row.RegionID.Valid {
		source.RegionID = &row.RegionID.Int32
	}

	// Build candidate index and attrs
	candidateIndex := make(map[int64]db.ListCoffeesForSimilarityRow, len(simRows))
	candidates := make([]similarity.CoffeeAttrs, 0, len(simRows))
	for _, sr := range simRows {
		candidateIndex[sr.ID] = sr
		attrs := similarity.CoffeeAttrs{
			CoffeeID:     sr.ID,
			TastingNotes: sr.TastingNotes,
			Process:      sr.Process.String,
			RoastLevel:   sr.RoastLevel.String,
			Variety:      sr.Variety.String,
		}
		if sr.RegionID.Valid {
			attrs.RegionID = &sr.RegionID.Int32
		}
		if sr.Latitude.Valid {
			attrs.Latitude = &sr.Latitude.Float64
		}
		if sr.Longitude.Valid {
			attrs.Longitude = &sr.Longitude.Float64
		}
		candidates = append(candidates, attrs)
	}

	ranked := similarity.Rank(source, candidates, 6)

	// Resolve ranked IDs to response objects
	if len(ranked) > 0 {
		// Fetch full details for similar coffees
		similar := make([]domain.SimilarCoffee, 0, len(ranked))
		for _, sc := range ranked {
			sr, ok := candidateIndex[sc.CoffeeID]
			if !ok {
				continue
			}
			// Get full coffee details for the similar coffee
			fullRow, err := h.queries.GetCoffeeByID(ctx, sc.CoffeeID)
			if err != nil {
				continue
			}
			similar = append(similar, domain.SimilarCoffee{
				ID:           fullRow.ID,
				Name:         fullRow.Name,
				RoasterName:  fullRow.RoasterName,
				RoasterSlug:  fullRow.RoasterSlug,
				ImageURL:     fullRow.ImageUrl.String,
				Process:      sr.Process.String,
				RoastLevel:   sr.RoastLevel.String,
				TastingNotes: sr.TastingNotes,
				Variety:      sr.Variety.String,
				Score:        sc.Score,
			})
		}
		resp.SimilarCoffees = similar
	}

	c.JSON(http.StatusOK, resp)
}

// GetStats godoc
// @Summary Get aggregate statistics
// @Tags stats
// @Produce json
// @Success 200 {object} domain.StatsResponse
// @Router /api/stats [get]
func (h *Handler) GetStats(c *gin.Context) {
	ctx := c.Request.Context()

	roasterCount, _ := h.queries.CountRoasters(ctx)
	coffeeCount, _ := h.queries.CountCoffees(ctx)
	originRows, _ := h.queries.ListDistinctOrigins(ctx)

	origins := make([]string, 0, len(originRows))
	for _, o := range originRows {
		if o.Valid {
			origins = append(origins, o.String)
		}
	}

	c.JSON(http.StatusOK, domain.StatsResponse{
		RoasterCount: roasterCount,
		CoffeeCount:  coffeeCount,
		Origins:      origins,
	})
}

func intQuery(c *gin.Context, key string, defaultVal int) int {
	v := c.Query(key)
	if v == "" {
		return defaultVal
	}
	n, err := strconv.Atoi(v)
	if err != nil || n < 1 {
		return defaultVal
	}
	return n
}

func coffeeDetailRowToResponse(row db.GetCoffeeByIDRow) domain.CoffeeResponse {
	return domain.CoffeeResponse{
		ID:           row.ID,
		RoasterID:    row.RoasterID,
		RoasterName:  row.RoasterName,
		RoasterSlug:  row.RoasterSlug,
		Name:         row.Name,
		ProductURL:   row.ProductUrl.String,
		ImageURL:     row.ImageUrl.String,
		CountryCode:  row.CountryCode.String,
		CountryName:  row.CountryName.String,
		RegionID:     row.RegionID.Int32,
		RegionName:   row.RegionName.String,
		ProducerID:   row.ProducerID.Int32,
		ProducerName: row.ProducerName.String,
		Process:      row.Process.String,
		RoastLevel:   row.RoastLevel.String,
		TastingNotes: row.TastingNotes,
		Variety:      row.Variety.String,
		Species:      row.Species.String,
		PriceCents:   row.PriceCents.Int32,
		WeightGrams:  row.WeightGrams.Int32,
		InStock:      row.InStock,
	}
}

func listRowToResponse(row db.ListCoffeesRow) domain.CoffeeResponse {
	return domain.CoffeeResponse{
		ID:           row.ID,
		RoasterID:    row.RoasterID,
		RoasterName:  row.RoasterName,
		RoasterSlug:  row.RoasterSlug,
		Name:         row.Name,
		ProductURL:   row.ProductUrl.String,
		ImageURL:     row.ImageUrl.String,
		CountryCode:  row.CountryCode.String,
		CountryName:  row.CountryName.String,
		RegionID:     row.CoffeeRegionID.Int32,
		RegionName:   row.RegionName.String,
		ProducerID:   row.CoffeeProducerID.Int32,
		ProducerName: row.ProducerName.String,
		Process:      row.Process.String,
		RoastLevel:   row.RoastLevel.String,
		TastingNotes: row.TastingNotes,
		Variety:      row.Variety.String,
		Species:      row.Species.String,
		PriceCents:   row.PriceCents.Int32,
		WeightGrams:  row.WeightGrams.Int32,
		InStock:      row.InStock,
	}
}

func coffeeRowToResponse(row db.ListCoffeesByRoasterRow) domain.CoffeeResponse {
	return domain.CoffeeResponse{
		ID:           row.ID,
		RoasterID:    row.RoasterID,
		RoasterName:  row.RoasterName,
		RoasterSlug:  row.RoasterSlug,
		Name:         row.Name,
		ProductURL:   row.ProductUrl.String,
		ImageURL:     row.ImageUrl.String,
		CountryCode:  row.CountryCode.String,
		CountryName:  row.CountryName.String,
		RegionID:     row.CoffeeRegionID.Int32,
		RegionName:   row.RegionName.String,
		ProducerID:   row.CoffeeProducerID.Int32,
		ProducerName: row.ProducerName.String,
		Process:      row.Process.String,
		RoastLevel:   row.RoastLevel.String,
		TastingNotes: row.TastingNotes,
		Variety:      row.Variety.String,
		Species:      row.Species.String,
		PriceCents:   row.PriceCents.Int32,
		WeightGrams:  row.WeightGrams.Int32,
		InStock:      row.InStock,
	}
}

func searchRowToResponse(row db.SearchCoffeesRow) domain.CoffeeResponse {
	return domain.CoffeeResponse{
		ID:           row.ID,
		RoasterID:    row.RoasterID,
		RoasterName:  row.RoasterName,
		RoasterSlug:  row.RoasterSlug,
		Name:         row.Name,
		ProductURL:   row.ProductUrl.String,
		ImageURL:     row.ImageUrl.String,
		CountryCode:  row.CountryCode.String,
		CountryName:  row.CountryName.String,
		RegionID:     row.CoffeeRegionID.Int32,
		RegionName:   row.RegionName.String,
		ProducerID:   row.CoffeeProducerID.Int32,
		ProducerName: row.ProducerName.String,
		Process:      row.Process.String,
		RoastLevel:   row.RoastLevel.String,
		TastingNotes: row.TastingNotes,
		Variety:      row.Variety.String,
		Species:      row.Species.String,
		PriceCents:   row.PriceCents.Int32,
		WeightGrams:  row.WeightGrams.Int32,
		InStock:      row.InStock,
	}
}

func filterRowToResponse(row db.FilterCoffeesRow) domain.CoffeeResponse {
	return domain.CoffeeResponse{
		ID:           row.ID,
		RoasterID:    row.RoasterID,
		RoasterName:  row.RoasterName,
		RoasterSlug:  row.RoasterSlug,
		Name:         row.Name,
		ProductURL:   row.ProductUrl.String,
		ImageURL:     row.ImageUrl.String,
		CountryCode:  row.CountryCode.String,
		CountryName:  row.CountryName.String,
		RegionID:     row.CoffeeRegionID.Int32,
		RegionName:   row.RegionName.String,
		ProducerID:   row.CoffeeProducerID.Int32,
		ProducerName: row.ProducerName.String,
		Process:      row.Process.String,
		RoastLevel:   row.RoastLevel.String,
		TastingNotes: row.TastingNotes,
		Variety:      row.Variety.String,
		Species:      row.Species.String,
		PriceCents:   row.PriceCents.Int32,
		WeightGrams:  row.WeightGrams.Int32,
		InStock:      row.InStock,
	}
}

// Helper for country/region/producer list row responses
func countryListRowToResponse(row db.ListCoffeesByCountryRow) domain.CoffeeResponse {
	return domain.CoffeeResponse{
		ID:           row.ID,
		RoasterID:    row.RoasterID,
		RoasterName:  row.RoasterName,
		RoasterSlug:  row.RoasterSlug,
		Name:         row.Name,
		ProductURL:   row.ProductUrl.String,
		ImageURL:     row.ImageUrl.String,
		CountryCode:  row.CountryCode.String,
		CountryName:  row.CountryName.String,
		RegionID:     row.CoffeeRegionID.Int32,
		RegionName:   row.RegionName.String,
		ProducerID:   row.CoffeeProducerID.Int32,
		ProducerName: row.ProducerName.String,
		Process:      row.Process.String,
		RoastLevel:   row.RoastLevel.String,
		TastingNotes: row.TastingNotes,
		Variety:      row.Variety.String,
		Species:      row.Species.String,
		PriceCents:   row.PriceCents.Int32,
		WeightGrams:  row.WeightGrams.Int32,
		InStock:      row.InStock,
	}
}

func regionListRowToResponse(row db.ListCoffeesByRegionRow) domain.CoffeeResponse {
	return domain.CoffeeResponse{
		ID:           row.ID,
		RoasterID:    row.RoasterID,
		RoasterName:  row.RoasterName,
		RoasterSlug:  row.RoasterSlug,
		Name:         row.Name,
		ProductURL:   row.ProductUrl.String,
		ImageURL:     row.ImageUrl.String,
		CountryCode:  row.CountryCode.String,
		CountryName:  row.CountryName.String,
		RegionID:     row.CoffeeRegionID.Int32,
		RegionName:   row.RegionName.String,
		ProducerID:   row.CoffeeProducerID.Int32,
		ProducerName: row.ProducerName.String,
		Process:      row.Process.String,
		RoastLevel:   row.RoastLevel.String,
		TastingNotes: row.TastingNotes,
		Variety:      row.Variety.String,
		Species:      row.Species.String,
		PriceCents:   row.PriceCents.Int32,
		WeightGrams:  row.WeightGrams.Int32,
		InStock:      row.InStock,
	}
}

func producerListRowToResponse(row db.ListCoffeesByProducerRow) domain.CoffeeResponse {
	return domain.CoffeeResponse{
		ID:           row.ID,
		RoasterID:    row.RoasterID,
		RoasterName:  row.RoasterName,
		RoasterSlug:  row.RoasterSlug,
		Name:         row.Name,
		ProductURL:   row.ProductUrl.String,
		ImageURL:     row.ImageUrl.String,
		CountryCode:  row.CountryCode.String,
		CountryName:  row.CountryName.String,
		RegionID:     row.CoffeeRegionID.Int32,
		RegionName:   row.RegionName.String,
		ProducerID:   row.CoffeeProducerID.Int32,
		ProducerName: row.ProducerName.String,
		Process:      row.Process.String,
		RoastLevel:   row.RoastLevel.String,
		TastingNotes: row.TastingNotes,
		Variety:      row.Variety.String,
		Species:      row.Species.String,
		PriceCents:   row.PriceCents.Int32,
		WeightGrams:  row.WeightGrams.Int32,
		InStock:      row.InStock,
	}
}
