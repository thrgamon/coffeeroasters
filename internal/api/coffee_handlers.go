package api

import (
	"log/slog"
	"net/http"
	"strconv"
	"strings"

	"github.com/gin-gonic/gin"
	"github.com/jackc/pgx/v5/pgtype"

	"github.com/thrgamon/coffeeroasters/internal/db"
	"github.com/thrgamon/coffeeroasters/internal/domain"
	"github.com/thrgamon/coffeeroasters/internal/similarity"
)

func textPtr(s string) pgtype.Text {
	s = strings.TrimSpace(s)
	if s == "" {
		return pgtype.Text{}
	}
	return pgtype.Text{String: s, Valid: true}
}

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
				ID:          r.ID,
				Slug:        r.Slug,
				Name:        r.Name,
				Website:     r.Website,
				State:       r.State.String,
				CoffeeCount: r.CoffeeCount,
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
				ID:          r.ID,
				Slug:        r.Slug,
				Name:        r.Name,
				Website:     r.Website,
				State:       r.State.String,
				CoffeeCount: r.CoffeeCount,
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
// @Param roaster_state query string false "Filter by roaster AU state code (e.g. VIC, NSW)"
// @Param similar_to query int false "Coffee ID to find similar coffees for"
// @Param liked query string false "Comma-separated liked coffee IDs for recommendations"
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
	roasterState := c.Query("roaster_state")
	similarToStr := c.Query("similar_to")

	var resp domain.CoffeeListResponse
	resp.Page = int32(page)
	resp.PageSize = int32(pageSize)

	if similarToStr != "" {
		h.listSimilarCoffees(c, similarToStr, pageSize, offset)
		return
	}

	likedStr := c.Query("liked")
	if likedStr != "" {
		h.listRecommendedCoffees(c, likedStr, pageSize, offset)
		return
	}

	rows, err := h.queries.ListCoffeesFiltered(ctx, db.ListCoffeesFilteredParams{
		Query:        textPtr(q),
		Origin:       textPtr(origin),
		Process:      textPtr(process),
		Roast:        textPtr(roast),
		Variety:      textPtr(varietyFilter),
		RoasterState: textPtr(roasterState),
		Lim:          int32(pageSize),
		Off:          int32(offset),
	})
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "failed to list coffees"})
		return
	}
	for _, row := range rows {
		resp.Coffees = append(resp.Coffees, filteredRowToResponse(row))
	}
	if len(rows) > 0 {
		resp.TotalCount = rows[0].TotalCount
	}

	if resp.Coffees == nil {
		resp.Coffees = []domain.CoffeeResponse{}
	}

	c.JSON(http.StatusOK, resp)
}

// listSimilarCoffees returns coffees ranked by similarity to a source coffee.
func (h *Handler) listSimilarCoffees(c *gin.Context, similarToStr string, pageSize, offset int) {
	ctx := c.Request.Context()

	sourceID, err := strconv.ParseInt(similarToStr, 10, 64)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid similar_to ID"})
		return
	}

	sourceRow, err := h.queries.GetCoffeeByID(ctx, sourceID)
	if err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "source coffee not found"})
		return
	}

	simRows, err := h.queries.ListCoffeesForSimilarity(ctx)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "failed to load coffees"})
		return
	}

	source := similarity.CoffeeAttrs{
		CoffeeID:     sourceRow.ID,
		TastingNotes: sourceRow.TastingNotes,
		Process:      sourceRow.Process.String,
		RoastLevel:   sourceRow.RoastLevel.String,
		Variety:      sourceRow.Variety.String,
	}
	if sourceRow.RegionID.Valid {
		source.RegionID = &sourceRow.RegionID.Int32
	}

	candidates := make([]similarity.CoffeeAttrs, 0, len(simRows))
	for _, sr := range simRows {
		attrs := simRowToAttrs(sr)
		if sr.ID == sourceID {
			source.Embedding = attrs.Embedding
		}
		candidates = append(candidates, attrs)
	}

	ranked := similarity.Rank(source, candidates, pageSize+offset)

	// Paginate the ranked results
	if offset > len(ranked) {
		ranked = nil
	} else if offset+pageSize > len(ranked) {
		ranked = ranked[offset:]
	} else {
		ranked = ranked[offset : offset+pageSize]
	}

	var resp domain.CoffeeListResponse
	resp.Page = int32(offset/pageSize + 1)
	resp.PageSize = int32(pageSize)

	coffees := make([]domain.CoffeeResponse, 0, len(ranked))
	for _, sc := range ranked {
		fullRow, err := h.queries.GetCoffeeByID(ctx, sc.CoffeeID)
		if err != nil {
			continue
		}
		coffee := coffeeDetailRowToResponse(fullRow)
		coffee.SimilarityScore = sc.Score
		coffees = append(coffees, coffee)
	}

	resp.Coffees = coffees
	resp.TotalCount = int64(len(candidates) - 1) // exclude source

	c.JSON(http.StatusOK, resp)
}

// listRecommendedCoffees returns coffees ranked by similarity to multiple liked coffees.
func (h *Handler) listRecommendedCoffees(c *gin.Context, likedStr string, pageSize, offset int) {
	ctx := c.Request.Context()

	likedIDs, err := parseIDList(likedStr, 50)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid liked IDs"})
		return
	}
	if len(likedIDs) == 0 {
		c.JSON(http.StatusOK, domain.CoffeeListResponse{
			Coffees:  []domain.CoffeeResponse{},
			Page:     1,
			PageSize: int32(pageSize),
		})
		return
	}

	simRows, err := h.queries.ListCoffeesForSimilarity(ctx)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "failed to load coffees"})
		return
	}

	// Build attrs index and candidate list
	attrsByID := make(map[int64]similarity.CoffeeAttrs, len(simRows))
	candidates := make([]similarity.CoffeeAttrs, 0, len(simRows))
	for _, sr := range simRows {
		attrs := simRowToAttrs(sr)
		attrsByID[sr.ID] = attrs
		candidates = append(candidates, attrs)
	}

	// Extract source attrs for liked IDs
	var sources []similarity.CoffeeAttrs
	for _, id := range likedIDs {
		if attrs, ok := attrsByID[id]; ok {
			sources = append(sources, attrs)
		}
	}

	if len(sources) == 0 {
		c.JSON(http.StatusOK, domain.CoffeeListResponse{
			Coffees:  []domain.CoffeeResponse{},
			Page:     1,
			PageSize: int32(pageSize),
		})
		return
	}

	ranked := similarity.RankFromMultiple(sources, candidates, pageSize+offset)

	// Paginate
	if offset > len(ranked) {
		ranked = nil
	} else if offset+pageSize > len(ranked) {
		ranked = ranked[offset:]
	} else {
		ranked = ranked[offset : offset+pageSize]
	}

	var resp domain.CoffeeListResponse
	resp.Page = int32(offset/pageSize + 1)
	resp.PageSize = int32(pageSize)

	coffees := make([]domain.CoffeeResponse, 0, len(ranked))
	for _, sc := range ranked {
		fullRow, err := h.queries.GetCoffeeByID(ctx, sc.CoffeeID)
		if err != nil {
			continue
		}
		coffee := coffeeDetailRowToResponse(fullRow)
		coffee.SimilarityScore = sc.Score
		coffees = append(coffees, coffee)
	}

	resp.Coffees = coffees
	resp.TotalCount = int64(len(candidates) - len(sources))

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

	// Fetch blend components if this is a blend
	if row.IsBlend {
		bcRows, err := h.queries.ListBlendComponents(ctx, int32(id))
		if err != nil {
			slog.Warn("list blend components", "error", err)
		} else {
			components := make([]domain.BlendComponentResponse, 0, len(bcRows))
			for _, bc := range bcRows {
				components = append(components, domain.BlendComponentResponse{
					CountryCode: bc.CountryCode.String,
					CountryName: bc.CountryName.String,
					RegionID:    bc.RegionID.Int32,
					RegionName:  bc.RegionName.String,
					Variety:     bc.Variety.String,
					Percentage:  bc.Percentage.Int32,
				})
			}
			resp.BlendComponents = components
		}
	}

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
		attrs := simRowToAttrs(sr)
		if sr.ID == id {
			source.Embedding = attrs.Embedding
		}
		candidates = append(candidates, attrs)
	}

	ranked := similarity.Rank(source, candidates, 3)

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
				Reasons:      sc.Reasons,
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
		ID:              row.ID,
		RoasterID:       row.RoasterID,
		RoasterName:     row.RoasterName,
		RoasterSlug:     row.RoasterSlug,
		Name:            row.Name,
		ProductURL:      row.ProductUrl.String,
		ImageURL:        row.ImageUrl.String,
		CountryCode:     row.CountryCode.String,
		CountryName:     row.CountryName.String,
		RegionID:        row.RegionID.Int32,
		RegionName:      row.RegionName.String,
		ProducerID:      row.ProducerID.Int32,
		ProducerName:    row.ProducerName.String,
		Process:         row.Process.String,
		RoastLevel:      row.RoastLevel.String,
		TastingNotes:    row.TastingNotes,
		Variety:         row.Variety.String,
		Species:         row.Species.String,
		PriceCents:      row.PriceCents.Int32,
		WeightGrams:     row.WeightGrams.Int32,
		PricePer100gMin: row.PricePer100gMin.Int32,
		PricePer100gMax: row.PricePer100gMax.Int32,
		IsBlend:         row.IsBlend,
		InStock:         row.InStock,
		Description:     row.Description.String,
	}
}

func filteredRowToResponse(row db.ListCoffeesFilteredRow) domain.CoffeeResponse {
	return domain.CoffeeResponse{
		ID:              row.ID,
		RoasterID:       row.RoasterID,
		RoasterName:     row.RoasterName,
		RoasterSlug:     row.RoasterSlug,
		Name:            row.Name,
		ProductURL:      row.ProductUrl.String,
		ImageURL:        row.ImageUrl.String,
		CountryCode:     row.CountryCode.String,
		CountryName:     row.CountryName.String,
		RegionID:        row.CoffeeRegionID.Int32,
		RegionName:      row.RegionName.String,
		ProducerID:      row.CoffeeProducerID.Int32,
		ProducerName:    row.ProducerName.String,
		Process:         row.Process.String,
		RoastLevel:      row.RoastLevel.String,
		TastingNotes:    row.TastingNotes,
		Variety:         row.Variety.String,
		Species:         row.Species.String,
		PriceCents:      row.PriceCents.Int32,
		WeightGrams:     row.WeightGrams.Int32,
		PricePer100gMin: row.PricePer100gMin.Int32,
		PricePer100gMax: row.PricePer100gMax.Int32,
		IsBlend:         row.IsBlend,
		InStock:         row.InStock,
	}
}

func coffeeRowToResponse(row db.ListCoffeesByRoasterRow) domain.CoffeeResponse {
	return domain.CoffeeResponse{
		ID:              row.ID,
		RoasterID:       row.RoasterID,
		RoasterName:     row.RoasterName,
		RoasterSlug:     row.RoasterSlug,
		Name:            row.Name,
		ProductURL:      row.ProductUrl.String,
		ImageURL:        row.ImageUrl.String,
		CountryCode:     row.CountryCode.String,
		CountryName:     row.CountryName.String,
		RegionID:        row.CoffeeRegionID.Int32,
		RegionName:      row.RegionName.String,
		ProducerID:      row.CoffeeProducerID.Int32,
		ProducerName:    row.ProducerName.String,
		Process:         row.Process.String,
		RoastLevel:      row.RoastLevel.String,
		TastingNotes:    row.TastingNotes,
		Variety:         row.Variety.String,
		Species:         row.Species.String,
		PriceCents:      row.PriceCents.Int32,
		WeightGrams:     row.WeightGrams.Int32,
		PricePer100gMin: row.PricePer100gMin.Int32,
		PricePer100gMax: row.PricePer100gMax.Int32,
		IsBlend:         row.IsBlend,
		InStock:         row.InStock,
	}
}


// Helper for country/region/producer list row responses
func countryListRowToResponse(row db.ListCoffeesByCountryRow) domain.CoffeeResponse {
	return domain.CoffeeResponse{
		ID:              row.ID,
		RoasterID:       row.RoasterID,
		RoasterName:     row.RoasterName,
		RoasterSlug:     row.RoasterSlug,
		Name:            row.Name,
		ProductURL:      row.ProductUrl.String,
		ImageURL:        row.ImageUrl.String,
		CountryCode:     row.CountryCode.String,
		CountryName:     row.CountryName.String,
		RegionID:        row.CoffeeRegionID.Int32,
		RegionName:      row.RegionName.String,
		ProducerID:      row.CoffeeProducerID.Int32,
		ProducerName:    row.ProducerName.String,
		Process:         row.Process.String,
		RoastLevel:      row.RoastLevel.String,
		TastingNotes:    row.TastingNotes,
		Variety:         row.Variety.String,
		Species:         row.Species.String,
		PriceCents:      row.PriceCents.Int32,
		WeightGrams:     row.WeightGrams.Int32,
		PricePer100gMin: row.PricePer100gMin.Int32,
		PricePer100gMax: row.PricePer100gMax.Int32,
		IsBlend:         row.IsBlend,
		InStock:         row.InStock,
	}
}

func regionListRowToResponse(row db.ListCoffeesByRegionRow) domain.CoffeeResponse {
	return domain.CoffeeResponse{
		ID:              row.ID,
		RoasterID:       row.RoasterID,
		RoasterName:     row.RoasterName,
		RoasterSlug:     row.RoasterSlug,
		Name:            row.Name,
		ProductURL:      row.ProductUrl.String,
		ImageURL:        row.ImageUrl.String,
		CountryCode:     row.CountryCode.String,
		CountryName:     row.CountryName.String,
		RegionID:        row.CoffeeRegionID.Int32,
		RegionName:      row.RegionName.String,
		ProducerID:      row.CoffeeProducerID.Int32,
		ProducerName:    row.ProducerName.String,
		Process:         row.Process.String,
		RoastLevel:      row.RoastLevel.String,
		TastingNotes:    row.TastingNotes,
		Variety:         row.Variety.String,
		Species:         row.Species.String,
		PriceCents:      row.PriceCents.Int32,
		WeightGrams:     row.WeightGrams.Int32,
		PricePer100gMin: row.PricePer100gMin.Int32,
		PricePer100gMax: row.PricePer100gMax.Int32,
		IsBlend:         row.IsBlend,
		InStock:         row.InStock,
	}
}

func simRowToAttrs(sr db.ListCoffeesForSimilarityRow) similarity.CoffeeAttrs {
	attrs := similarity.CoffeeAttrs{
		CoffeeID:     sr.ID,
		TastingNotes: sr.TastingNotes,
		Process:      sr.Process.String,
		RoastLevel:   sr.RoastLevel.String,
		Variety:      sr.Variety.String,
		Embedding:    sr.Embedding,
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
	return attrs
}

// parseIDList parses a comma-separated string of IDs, deduplicating and validating.
// Returns at most maxIDs entries.
func parseIDList(s string, maxIDs int) ([]int64, error) {
	parts := strings.Split(s, ",")
	seen := make(map[int64]bool, len(parts))
	var ids []int64
	for _, p := range parts {
		p = strings.TrimSpace(p)
		if p == "" {
			continue
		}
		id, err := strconv.ParseInt(p, 10, 64)
		if err != nil {
			return nil, err
		}
		if seen[id] {
			continue
		}
		seen[id] = true
		ids = append(ids, id)
		if len(ids) >= maxIDs {
			break
		}
	}
	return ids, nil
}

func producerListRowToResponse(row db.ListCoffeesByProducerRow) domain.CoffeeResponse {
	return domain.CoffeeResponse{
		ID:              row.ID,
		RoasterID:       row.RoasterID,
		RoasterName:     row.RoasterName,
		RoasterSlug:     row.RoasterSlug,
		Name:            row.Name,
		ProductURL:      row.ProductUrl.String,
		ImageURL:        row.ImageUrl.String,
		CountryCode:     row.CountryCode.String,
		CountryName:     row.CountryName.String,
		RegionID:        row.CoffeeRegionID.Int32,
		RegionName:      row.RegionName.String,
		ProducerID:      row.CoffeeProducerID.Int32,
		ProducerName:    row.ProducerName.String,
		Process:         row.Process.String,
		RoastLevel:      row.RoastLevel.String,
		TastingNotes:    row.TastingNotes,
		Variety:         row.Variety.String,
		Species:         row.Species.String,
		PriceCents:      row.PriceCents.Int32,
		WeightGrams:     row.WeightGrams.Int32,
		PricePer100gMin: row.PricePer100gMin.Int32,
		PricePer100gMax: row.PricePer100gMax.Int32,
		IsBlend:         row.IsBlend,
		InStock:         row.InStock,
	}
}
