package api

import (
	"net/http"
	"strconv"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/jackc/pgx/v5/pgtype"

	"github.com/thrgamon/coffeeroasters/internal/auth"
	"github.com/thrgamon/coffeeroasters/internal/db"
	"github.com/thrgamon/coffeeroasters/internal/domain"
)

func (h *Handler) UpsertUserCoffee(c *gin.Context) {
	userID, ok := auth.GetUserID(c)
	if !ok {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "authentication required"})
		return
	}

	var req domain.UserCoffeeRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	if req.Rating != nil && (*req.Rating < 1 || *req.Rating > 5) {
		c.JSON(http.StatusBadRequest, gin.H{"error": "rating must be between 1 and 5"})
		return
	}

	var liked pgtype.Bool
	if req.Liked != nil {
		liked = pgtype.Bool{Bool: *req.Liked, Valid: true}
	}

	var rating pgtype.Int2
	if req.Rating != nil {
		rating = pgtype.Int2{Int16: *req.Rating, Valid: true}
	}

	var review pgtype.Text
	if req.Review != nil {
		review = pgtype.Text{String: *req.Review, Valid: true}
	}

	var drunkAt pgtype.Date
	if req.DrunkAt != nil {
		t, err := time.Parse("2006-01-02", *req.DrunkAt)
		if err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": "drunk_at must be YYYY-MM-DD"})
			return
		}
		drunkAt = pgtype.Date{Time: t, Valid: true}
	}

	uc, err := h.queries.UpsertUserCoffee(c.Request.Context(), db.UpsertUserCoffeeParams{
		UserID:   userID,
		CoffeeID: req.CoffeeID,
		Status:   req.Status,
		Liked:    liked,
		Rating:   rating,
		Review:   review,
		DrunkAt:  drunkAt,
	})
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "failed to save coffee"})
		return
	}

	c.JSON(http.StatusOK, userCoffeeToResponse(uc))
}

func (h *Handler) DeleteUserCoffee(c *gin.Context) {
	userID, ok := auth.GetUserID(c)
	if !ok {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "authentication required"})
		return
	}

	coffeeID, err := strconv.ParseInt(c.Param("coffee_id"), 10, 64)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid coffee_id"})
		return
	}

	if err := h.queries.DeleteUserCoffee(c.Request.Context(), db.DeleteUserCoffeeParams{
		UserID:   userID,
		CoffeeID: coffeeID,
	}); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "failed to remove coffee"})
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "removed"})
}

func (h *Handler) ListUserCoffees(c *gin.Context) {
	userID, ok := auth.GetUserID(c)
	if !ok {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "authentication required"})
		return
	}

	rows, err := h.queries.ListUserCoffees(c.Request.Context(), userID)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "failed to list coffees"})
		return
	}

	coffees := make([]domain.UserCoffeeDetailResponse, 0, len(rows))
	for _, row := range rows {
		detail := domain.UserCoffeeDetailResponse{
			CoffeeResponse: domain.CoffeeResponse{
				ID:              row.CoffeeID,
				Name:            row.CoffeeName,
				RoasterName:     row.RoasterName,
				RoasterSlug:     row.RoasterSlug,
				RoasterLogoURL:  row.RoasterLogoUrl.String,
				ImageURL:        row.CoffeeImageUrl.String,
				Process:         row.Process.String,
				RoastLevel:      row.RoastLevel.String,
				TastingNotes:    row.TastingNotes,
				Variety:         row.Variety.String,
				PriceCents:      row.PriceCents.Int32,
				WeightGrams:     row.WeightGrams.Int32,
				PricePer100gMin: row.PricePer100gMin.Int32,
				PricePer100gMax: row.PricePer100gMax.Int32,
				InStock:         row.InStock,
				IsBlend:         row.IsBlend,
				IsDecaf:         row.IsDecaf,
				ProductURL:      row.ProductUrl.String,
				CountryCode:     row.CountryCode.String,
				CountryName:     row.CountryName.String,
				RegionID:        row.RegionID.Int32,
				RegionName:      row.RegionName.String,
				Species:         row.Species.String,
				Description:     row.Description.String,
			},
			Status: row.Status,
		}
		if row.Liked.Valid {
			detail.Liked = &row.Liked.Bool
		}
		if row.Rating.Valid {
			detail.Rating = &row.Rating.Int16
		}
		if row.Review.Valid {
			detail.Review = &row.Review.String
		}
		if row.DrunkAt.Valid {
			s := row.DrunkAt.Time.Format("2006-01-02")
			detail.DrunkAt = &s
		}
		coffees = append(coffees, detail)
	}

	c.JSON(http.StatusOK, domain.UserCoffeeListResponse{Coffees: coffees})
}

func (h *Handler) ListUserCoffeeIDs(c *gin.Context) {
	userID, ok := auth.GetUserID(c)
	if !ok {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "authentication required"})
		return
	}

	rows, err := h.queries.ListUserCoffeeIDs(c.Request.Context(), userID)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "failed to list coffee IDs"})
		return
	}

	items := make([]domain.UserCoffeeResponse, 0, len(rows))
	for _, row := range rows {
		item := domain.UserCoffeeResponse{
			CoffeeID: row.CoffeeID,
			Status:   row.Status,
		}
		if row.Liked.Valid {
			item.Liked = &row.Liked.Bool
		}
		if row.Rating.Valid {
			item.Rating = &row.Rating.Int16
		}
		items = append(items, item)
	}

	c.JSON(http.StatusOK, items)
}

func userCoffeeToResponse(uc db.UserCoffee) domain.UserCoffeeResponse {
	resp := domain.UserCoffeeResponse{
		CoffeeID: uc.CoffeeID,
		Status:   uc.Status,
	}
	if uc.Liked.Valid {
		resp.Liked = &uc.Liked.Bool
	}
	if uc.Rating.Valid {
		resp.Rating = &uc.Rating.Int16
	}
	if uc.Review.Valid {
		resp.Review = &uc.Review.String
	}
	if uc.DrunkAt.Valid {
		s := uc.DrunkAt.Time.Format("2006-01-02")
		resp.DrunkAt = &s
	}
	return resp
}
