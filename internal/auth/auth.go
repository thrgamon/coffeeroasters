package auth

import (
	"context"
	"crypto/rand"
	"encoding/hex"
	"errors"
	"fmt"
	"log/slog"
	"time"

	"github.com/thrgamon/coffeeroasters/internal/config"
	"github.com/thrgamon/coffeeroasters/internal/db"
	"github.com/thrgamon/coffeeroasters/internal/domain"
)

type Service struct {
	queries *db.Queries
	cfg     config.Config
}

func NewService(queries *db.Queries, cfg config.Config) *Service {
	return &Service{queries: queries, cfg: cfg}
}

// SendMagicLink creates a magic link token for the given email.
// In production this would send an email; in development it returns the token.
func (s *Service) SendMagicLink(ctx context.Context, email string) (string, error) {
	token, err := generateToken()
	if err != nil {
		return "", fmt.Errorf("generating token: %w", err)
	}

	_, err = s.queries.CreateMagicLink(ctx, db.CreateMagicLinkParams{
		Email:     email,
		Token:     token,
		ExpiresAt: time.Now().Add(15 * time.Minute),
	})
	if err != nil {
		return "", fmt.Errorf("creating magic link: %w", err)
	}

	if s.cfg.Environment == "production" {
		// In production, send the email here
		slog.Info("magic link created", "email", email)
		return "", nil
	}

	// In development, return the token directly for testing
	return token, nil
}

// VerifyMagicLink verifies a magic link token and creates a session.
// If the user doesn't exist, they are created automatically.
func (s *Service) VerifyMagicLink(ctx context.Context, token string) (*domain.AuthResponse, string, error) {
	link, err := s.queries.GetMagicLinkByToken(ctx, token)
	if err != nil {
		return nil, "", errors.New("invalid or expired link")
	}

	// Mark the link as used
	if err := s.queries.MarkMagicLinkUsed(ctx, token); err != nil {
		return nil, "", fmt.Errorf("marking link used: %w", err)
	}

	// Find or create user
	user, err := s.queries.GetUserByEmail(ctx, link.Email)
	if err != nil {
		// User doesn't exist, create them
		user, err = s.queries.CreateUserPasswordless(ctx, link.Email)
		if err != nil {
			return nil, "", fmt.Errorf("creating user: %w", err)
		}
	}

	sessionToken, err := s.createSession(ctx, user.ID)
	if err != nil {
		return nil, "", fmt.Errorf("creating session: %w", err)
	}

	return &domain.AuthResponse{
		User: domain.UserResponse{ID: user.ID, Email: user.Email, IsAdmin: user.IsAdmin},
	}, sessionToken, nil
}

func (s *Service) Logout(ctx context.Context, token string) error {
	return s.queries.DeleteSessionByToken(ctx, token)
}

func (s *Service) ValidateSession(ctx context.Context, token string) (*db.GetSessionByTokenRow, error) {
	session, err := s.queries.GetSessionByToken(ctx, token)
	if err != nil {
		return nil, errors.New("invalid session")
	}
	return &session, nil
}

func (s *Service) DeleteExpiredSessions(ctx context.Context) error {
	return s.queries.DeleteExpiredSessions(ctx)
}

func (s *Service) DeleteExpiredMagicLinks(ctx context.Context) error {
	return s.queries.DeleteExpiredMagicLinks(ctx)
}

func (s *Service) createSession(ctx context.Context, userID int32) (string, error) {
	token, err := generateToken()
	if err != nil {
		return "", fmt.Errorf("generating token: %w", err)
	}

	_, err = s.queries.CreateSession(ctx, db.CreateSessionParams{
		UserID:    userID,
		Token:     token,
		ExpiresAt: time.Now().Add(s.cfg.SessionMaxAge),
	})
	if err != nil {
		return "", fmt.Errorf("storing session: %w", err)
	}

	return token, nil
}

func generateToken() (string, error) {
	b := make([]byte, 32)
	if _, err := rand.Read(b); err != nil {
		return "", err
	}
	return hex.EncodeToString(b), nil
}
