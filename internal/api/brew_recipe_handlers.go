package api

import (
	"fmt"
	"net/http"
	"strconv"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/jackc/pgx/v5/pgtype"

	"github.com/thrgamon/coffeeroasters/internal/auth"
	"github.com/thrgamon/coffeeroasters/internal/db"
	"github.com/thrgamon/coffeeroasters/internal/domain"
)

func (h *Handler) CreateBrewRecipe(c *gin.Context) {
	userID, ok := auth.GetUserID(c)
	if !ok {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "authentication required"})
		return
	}

	var req domain.BrewRecipeRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	isPublic := true
	if req.IsPublic != nil {
		isPublic = *req.IsPublic
	}

	recipe, err := h.queries.CreateBrewRecipe(c.Request.Context(), db.CreateBrewRecipeParams{
		CoffeeID:        req.CoffeeID,
		UserID:          userID,
		Title:           req.Title,
		BrewMethod:      req.BrewMethod,
		DoseGrams:       numericVal(req.DoseGrams),
		WaterMl:         int4PtrVal(req.WaterMl),
		WaterTempC:      int4PtrVal(req.WaterTempC),
		GrindSize:       textPtrVal(req.GrindSize),
		BrewTimeSeconds: int4PtrVal(req.BrewTimeSeconds),
		Notes:           textPtrVal(req.Notes),
		IsPublic:        isPublic,
	})
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "failed to create recipe"})
		return
	}

	c.JSON(http.StatusCreated, brewRecipeToResponse(recipe))
}

func (h *Handler) UpdateBrewRecipe(c *gin.Context) {
	userID, ok := auth.GetUserID(c)
	if !ok {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "authentication required"})
		return
	}

	id, err := strconv.ParseInt(c.Param("id"), 10, 32)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid recipe ID"})
		return
	}

	var req domain.BrewRecipeRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	isPublic := true
	if req.IsPublic != nil {
		isPublic = *req.IsPublic
	}

	recipe, err := h.queries.UpdateBrewRecipe(c.Request.Context(), db.UpdateBrewRecipeParams{
		ID:              int32(id),
		UserID:          userID,
		Title:           req.Title,
		BrewMethod:      req.BrewMethod,
		DoseGrams:       numericVal(req.DoseGrams),
		WaterMl:         int4PtrVal(req.WaterMl),
		WaterTempC:      int4PtrVal(req.WaterTempC),
		GrindSize:       textPtrVal(req.GrindSize),
		BrewTimeSeconds: int4PtrVal(req.BrewTimeSeconds),
		Notes:           textPtrVal(req.Notes),
		IsPublic:        isPublic,
	})
	if err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "recipe not found or not yours"})
		return
	}

	c.JSON(http.StatusOK, brewRecipeToResponse(recipe))
}

func (h *Handler) DeleteBrewRecipe(c *gin.Context) {
	userID, ok := auth.GetUserID(c)
	if !ok {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "authentication required"})
		return
	}

	id, err := strconv.ParseInt(c.Param("id"), 10, 32)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid recipe ID"})
		return
	}

	if err := h.queries.DeleteBrewRecipe(c.Request.Context(), db.DeleteBrewRecipeParams{
		ID:     int32(id),
		UserID: userID,
	}); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "failed to delete recipe"})
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "deleted"})
}

func (h *Handler) ListBrewRecipesByCoffee(c *gin.Context) {
	coffeeID, err := strconv.ParseInt(c.Param("id"), 10, 64)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid coffee id"})
		return
	}

	// Try to identify the viewer for private recipe visibility (public endpoint)
	var viewerUserID pgtype.Int4
	if token, cookieErr := c.Cookie("session_token"); cookieErr == nil && token != "" {
		if session, sessionErr := h.auth.ValidateSession(c.Request.Context(), token); sessionErr == nil {
			viewerUserID = pgtype.Int4{Int32: session.UserID, Valid: true}
		}
	}

	rows, err := h.queries.ListBrewRecipesByCoffee(c.Request.Context(), db.ListBrewRecipesByCoffeeParams{
		CoffeeID:     coffeeID,
		ViewerUserID: viewerUserID,
	})
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "failed to list recipes"})
		return
	}

	recipes := make([]domain.BrewRecipeResponse, 0, len(rows))
	for _, row := range rows {
		recipes = append(recipes, brewRecipeByCoffeeRowToResponse(row))
	}

	c.JSON(http.StatusOK, domain.BrewRecipeListResponse{Recipes: recipes})
}

func (h *Handler) ListBrewRecipesByUser(c *gin.Context) {
	userID, ok := auth.GetUserID(c)
	if !ok {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "authentication required"})
		return
	}

	rows, err := h.queries.ListBrewRecipesByUser(c.Request.Context(), userID)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "failed to list recipes"})
		return
	}

	recipes := make([]domain.BrewRecipeResponse, 0, len(rows))
	for _, row := range rows {
		recipes = append(recipes, brewRecipeByUserRowToResponse(row))
	}

	c.JSON(http.StatusOK, domain.BrewRecipeListResponse{Recipes: recipes})
}

func brewRecipeToResponse(r db.BrewRecipe) domain.BrewRecipeResponse {
	resp := domain.BrewRecipeResponse{
		ID:         r.ID,
		CoffeeID:   r.CoffeeID,
		UserID:     r.UserID,
		Title:      r.Title,
		BrewMethod: r.BrewMethod,
		IsPublic:   r.IsPublic,
		CreatedAt:  r.CreatedAt.Format(time.RFC3339),
		UpdatedAt:  r.UpdatedAt.Format(time.RFC3339),
	}
	if r.DoseGrams.Valid {
		resp.DoseGrams = numericToFloat64Ptr(r.DoseGrams)
	}
	if r.WaterMl.Valid {
		v := r.WaterMl.Int32
		resp.WaterMl = &v
	}
	if r.WaterTempC.Valid {
		v := r.WaterTempC.Int32
		resp.WaterTempC = &v
	}
	if r.GrindSize.Valid {
		resp.GrindSize = &r.GrindSize.String
	}
	if r.BrewTimeSeconds.Valid {
		v := r.BrewTimeSeconds.Int32
		resp.BrewTimeSeconds = &v
	}
	if r.Notes.Valid {
		resp.Notes = &r.Notes.String
	}
	return resp
}

func brewRecipeByCoffeeRowToResponse(row db.ListBrewRecipesByCoffeeRow) domain.BrewRecipeResponse {
	resp := domain.BrewRecipeResponse{
		ID:         row.ID,
		CoffeeID:   row.CoffeeID,
		UserID:     row.UserID,
		UserEmail:  row.UserEmail,
		Title:      row.Title,
		BrewMethod: row.BrewMethod,
		IsPublic:   row.IsPublic,
		CreatedAt:  row.CreatedAt.Format(time.RFC3339),
		UpdatedAt:  row.UpdatedAt.Format(time.RFC3339),
	}
	if row.DoseGrams.Valid {
		resp.DoseGrams = numericToFloat64Ptr(row.DoseGrams)
	}
	if row.WaterMl.Valid {
		v := row.WaterMl.Int32
		resp.WaterMl = &v
	}
	if row.WaterTempC.Valid {
		v := row.WaterTempC.Int32
		resp.WaterTempC = &v
	}
	if row.GrindSize.Valid {
		resp.GrindSize = &row.GrindSize.String
	}
	if row.BrewTimeSeconds.Valid {
		v := row.BrewTimeSeconds.Int32
		resp.BrewTimeSeconds = &v
	}
	if row.Notes.Valid {
		resp.Notes = &row.Notes.String
	}
	return resp
}

func brewRecipeByUserRowToResponse(row db.ListBrewRecipesByUserRow) domain.BrewRecipeResponse {
	resp := domain.BrewRecipeResponse{
		ID:          row.ID,
		CoffeeID:    row.CoffeeID,
		UserID:      row.UserID,
		Title:       row.Title,
		BrewMethod:  row.BrewMethod,
		IsPublic:    row.IsPublic,
		CoffeeName:  row.CoffeeName,
		RoasterName: row.RoasterName,
		RoasterSlug: row.RoasterSlug,
		CreatedAt:   row.CreatedAt.Format(time.RFC3339),
		UpdatedAt:   row.UpdatedAt.Format(time.RFC3339),
	}
	if row.DoseGrams.Valid {
		resp.DoseGrams = numericToFloat64Ptr(row.DoseGrams)
	}
	if row.WaterMl.Valid {
		v := row.WaterMl.Int32
		resp.WaterMl = &v
	}
	if row.WaterTempC.Valid {
		v := row.WaterTempC.Int32
		resp.WaterTempC = &v
	}
	if row.GrindSize.Valid {
		resp.GrindSize = &row.GrindSize.String
	}
	if row.BrewTimeSeconds.Valid {
		v := row.BrewTimeSeconds.Int32
		resp.BrewTimeSeconds = &v
	}
	if row.Notes.Valid {
		resp.Notes = &row.Notes.String
	}
	return resp
}

func numericVal(v *float64) pgtype.Numeric {
	if v == nil {
		return pgtype.Numeric{}
	}
	var n pgtype.Numeric
	_ = n.Scan(fmt.Sprintf("%g", *v))
	return n
}

func numericToFloat64Ptr(n pgtype.Numeric) *float64 {
	if !n.Valid {
		return nil
	}
	f8, err := n.Float64Value()
	if err != nil || !f8.Valid {
		return nil
	}
	return &f8.Float64
}

func int4PtrVal(v *int32) pgtype.Int4 {
	if v == nil {
		return pgtype.Int4{}
	}
	return pgtype.Int4{Int32: *v, Valid: true}
}

func textPtrVal(v *string) pgtype.Text {
	if v == nil {
		return pgtype.Text{}
	}
	return pgtype.Text{String: *v, Valid: true}
}
