package api

import (
	"net/http"
	"strconv"
	"strings"

	"github.com/gin-gonic/gin"

	"github.com/thrgamon/coffeeroasters/internal/auth"
	"github.com/thrgamon/coffeeroasters/internal/db"
	"github.com/thrgamon/coffeeroasters/internal/domain"
)

// AddTastingNoteVote adds the current user's vote for a tasting note on a coffee.
func (h *Handler) AddTastingNoteVote(c *gin.Context) {
	userID, ok := auth.GetUserID(c)
	if !ok {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "authentication required"})
		return
	}

	var req domain.TastingNoteVoteRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	note := strings.TrimSpace(strings.ToLower(req.TastingNote))
	if note == "" || len(note) > 100 {
		c.JSON(http.StatusBadRequest, gin.H{"error": "tasting note must be 1-100 characters"})
		return
	}

	_, err := h.queries.AddUserTastingNote(c.Request.Context(), db.AddUserTastingNoteParams{
		UserID:      userID,
		CoffeeID:    req.CoffeeID,
		TastingNote: note,
	})
	if err != nil {
		// ON CONFLICT DO NOTHING returns no rows — treat as success (already voted)
		if err.Error() == "no rows in result set" {
			c.JSON(http.StatusOK, gin.H{"message": "already voted"})
			return
		}
		c.JSON(http.StatusInternalServerError, gin.H{"error": "failed to add vote"})
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "vote added"})
}

// RemoveTastingNoteVote removes the current user's vote for a tasting note on a coffee.
func (h *Handler) RemoveTastingNoteVote(c *gin.Context) {
	userID, ok := auth.GetUserID(c)
	if !ok {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "authentication required"})
		return
	}

	var req domain.TastingNoteVoteRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	note := strings.TrimSpace(strings.ToLower(req.TastingNote))

	err := h.queries.RemoveUserTastingNote(c.Request.Context(), db.RemoveUserTastingNoteParams{
		UserID:      userID,
		CoffeeID:    req.CoffeeID,
		TastingNote: note,
	})
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "failed to remove vote"})
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "vote removed"})
}

// GetTastingNoteVotes returns crowdsourced tasting notes for a coffee, plus the
// current user's votes if authenticated.
func (h *Handler) GetTastingNoteVotes(c *gin.Context) {
	coffeeID, err := strconv.ParseInt(c.Param("id"), 10, 64)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid coffee_id"})
		return
	}

	ctx := c.Request.Context()

	crowdsourced, err := h.queries.ListCrowdsourcedTastingNotes(ctx, coffeeID)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "failed to list tasting notes"})
		return
	}

	notes := make([]domain.CrowdsourcedTastingNote, 0, len(crowdsourced))
	for _, row := range crowdsourced {
		notes = append(notes, domain.CrowdsourcedTastingNote{
			Note:      row.TastingNote,
			VoteCount: row.VoteCount,
		})
	}

	var userVotes []string

	// If user is authenticated, get their votes too
	userID, ok := auth.GetUserID(c)
	if ok {
		userVotes, err = h.queries.ListUserTastingNotesForCoffee(ctx, db.ListUserTastingNotesForCoffeeParams{
			UserID:   userID,
			CoffeeID: coffeeID,
		})
		if err != nil {
			userVotes = []string{}
		}
	}

	if userVotes == nil {
		userVotes = []string{}
	}

	c.JSON(http.StatusOK, domain.TastingNoteVotesResponse{
		CrowdsourcedNotes: notes,
		UserVotes:         userVotes,
	})
}
