package api

import (
	"net/http"
	"strconv"

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

	var enjoyed pgtype.Bool
	if req.Enjoyed != nil {
		enjoyed = pgtype.Bool{Bool: *req.Enjoyed, Valid: true}
	}

	uc, err := h.queries.UpsertUserCoffee(c.Request.Context(), db.UpsertUserCoffeeParams{
		UserID:   userID,
		CoffeeID: req.CoffeeID,
		Status:   req.Status,
		Enjoyed:  enjoyed,
	})
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "failed to save coffee"})
		return
	}

	resp := domain.UserCoffeeResponse{
		CoffeeID: uc.CoffeeID,
		Status:   uc.Status,
	}
	if uc.Enjoyed.Valid {
		resp.Enjoyed = &uc.Enjoyed.Bool
	}

	c.JSON(http.StatusOK, resp)
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
		if row.Enjoyed.Valid {
			detail.Enjoyed = &row.Enjoyed.Bool
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
		if row.Enjoyed.Valid {
			item.Enjoyed = &row.Enjoyed.Bool
		}
		items = append(items, item)
	}

	c.JSON(http.StatusOK, items)
}
