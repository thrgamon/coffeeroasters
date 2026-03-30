package api

import (
	"net/http"

	"github.com/gin-gonic/gin"

	"github.com/thrgamon/coffeeroasters/internal/auth"
	"github.com/thrgamon/coffeeroasters/internal/config"
	"github.com/thrgamon/coffeeroasters/internal/db"
	"github.com/thrgamon/coffeeroasters/internal/domain"
)

type HandlerConfig struct {
	Auth    *auth.Service
	Cfg     config.Config
	Queries *db.Queries
}

type Handler struct {
	auth    *auth.Service
	cfg     config.Config
	queries *db.Queries
}

func NewHandler(cfg HandlerConfig) *Handler {
	return &Handler{auth: cfg.Auth, cfg: cfg.Cfg, queries: cfg.Queries}
}

// Routes registers all HTTP routes on the given router group.
func (h *Handler) Routes(rg *gin.RouterGroup) {
	rg.GET("/health", h.Health)

	authGroup := rg.Group("/auth")
	{
		authGroup.POST("/magic-link", h.SendMagicLink)
		authGroup.POST("/verify", h.VerifyMagicLink)
		authGroup.POST("/logout", h.Logout)
		authGroup.GET("/me", h.Me)
	}

	// Public read-only endpoints
	rg.GET("/roasters", h.ListRoasters)
	rg.GET("/roasters/:slug", h.GetRoaster)
	rg.GET("/coffees", h.ListCoffees)
	rg.GET("/coffees/find", h.FindCoffees)
	rg.GET("/coffees/:id", h.GetCoffee)
	rg.GET("/cafes", h.ListCafes)
	rg.GET("/cafes/:slug", h.GetCafe)
	rg.GET("/stats", h.GetStats)
	rg.GET("/countries", h.ListCountries)
	rg.GET("/countries/:code", h.GetCountry)
	rg.GET("/regions/:id", h.GetRegion)
	rg.GET("/producers/:id", h.GetProducer)

	protected := rg.Group("")
	protected.Use(auth.RequireAuth(h.auth))
	{
		protected.GET("/dashboard", h.Dashboard)
		protected.POST("/user/coffees", h.UpsertUserCoffee)
		protected.DELETE("/user/coffees/:coffee_id", h.DeleteUserCoffee)
		protected.GET("/user/coffees", h.ListUserCoffees)
		protected.GET("/user/coffee-ids", h.ListUserCoffeeIDs)
	}

	admin := rg.Group("/admin")
	admin.Use(auth.RequireAuth(h.auth), auth.RequireAdmin())
	{
		admin.GET("", h.AdminDashboard)
	}
}

// Health godoc
// @Summary Health check
// @Tags health
// @Produce json
// @Success 200 {object} map[string]string
// @Router /api/health [get]
func (h *Handler) Health(c *gin.Context) {
	c.JSON(http.StatusOK, gin.H{"status": "ok"})
}

// SendMagicLink godoc
// @Summary Request a magic link for passwordless login
// @Tags auth
// @Accept json
// @Produce json
// @Param body body domain.MagicLinkRequest true "Email address"
// @Success 200 {object} domain.MagicLinkResponse
// @Failure 400 {object} map[string]string
// @Router /api/auth/magic-link [post]
func (h *Handler) SendMagicLink(c *gin.Context) {
	var req domain.MagicLinkRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	token, err := h.auth.SendMagicLink(c.Request.Context(), req.Email)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "failed to send magic link"})
		return
	}

	resp := domain.MagicLinkResponse{
		Message: "Check your email for a login link",
	}
	if token != "" {
		// Development mode: return token directly
		resp.Token = token
	}
	c.JSON(http.StatusOK, resp)
}

// VerifyMagicLink godoc
// @Summary Verify a magic link token and create a session
// @Tags auth
// @Accept json
// @Produce json
// @Param body body domain.VerifyMagicLinkRequest true "Magic link token"
// @Success 200 {object} domain.AuthResponse
// @Failure 401 {object} map[string]string
// @Router /api/auth/verify [post]
func (h *Handler) VerifyMagicLink(c *gin.Context) {
	var req domain.VerifyMagicLinkRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	resp, token, err := h.auth.VerifyMagicLink(c.Request.Context(), req.Token)
	if err != nil {
		c.JSON(http.StatusUnauthorized, gin.H{"error": err.Error()})
		return
	}

	h.setSessionCookie(c, token)
	c.JSON(http.StatusOK, resp)
}

// Logout godoc
// @Summary Logout current session
// @Tags auth
// @Produce json
// @Success 200 {object} map[string]string
// @Router /api/auth/logout [post]
func (h *Handler) Logout(c *gin.Context) {
	token, err := c.Cookie("session_token")
	if err == nil && token != "" {
		_ = h.auth.Logout(c.Request.Context(), token)
	}

	h.clearSessionCookie(c)
	c.JSON(http.StatusOK, gin.H{"message": "logged out"})
}

// Me godoc
// @Summary Get current user (returns null user when unauthenticated)
// @Tags auth
// @Produce json
// @Success 200 {object} domain.MeResponse
// @Router /api/auth/me [get]
func (h *Handler) Me(c *gin.Context) {
	token, err := c.Cookie("session_token")
	if err != nil || token == "" {
		c.JSON(http.StatusOK, domain.MeResponse{})
		return
	}

	session, err := h.auth.ValidateSession(c.Request.Context(), token)
	if err != nil {
		c.JSON(http.StatusOK, domain.MeResponse{})
		return
	}

	user := &domain.UserResponse{
		ID:      session.UserID,
		Email:   session.UserEmail,
		IsAdmin: session.UserIsAdmin,
	}
	c.JSON(http.StatusOK, domain.MeResponse{User: user})
}

// Dashboard godoc
// @Summary Example protected endpoint
// @Tags dashboard
// @Produce json
// @Success 200 {object} map[string]string
// @Failure 401 {object} map[string]string
// @Router /api/dashboard [get]
func (h *Handler) Dashboard(c *gin.Context) {
	email, _ := c.Get("user_email")
	c.JSON(http.StatusOK, gin.H{
		"message": "welcome to the dashboard",
		"email":   email,
	})
}

// AdminDashboard godoc
// @Summary Admin dashboard (admin only)
// @Tags admin
// @Produce json
// @Success 200 {object} map[string]string
// @Failure 403 {object} map[string]string
// @Router /api/admin [get]
func (h *Handler) AdminDashboard(c *gin.Context) {
	email, _ := c.Get("user_email")
	c.JSON(http.StatusOK, gin.H{
		"message": "admin dashboard",
		"email":   email,
	})
}

func (h *Handler) setSessionCookie(c *gin.Context, token string) {
	maxAge := int(h.cfg.SessionMaxAge.Seconds())
	c.SetSameSite(http.SameSiteLaxMode)
	c.SetCookie("session_token", token, maxAge, "/", h.cfg.CookieDomain, h.cfg.CookieSecure, true)
}

func (h *Handler) clearSessionCookie(c *gin.Context) {
	c.SetSameSite(http.SameSiteLaxMode)
	c.SetCookie("session_token", "", -1, "/", h.cfg.CookieDomain, h.cfg.CookieSecure, true)
}
