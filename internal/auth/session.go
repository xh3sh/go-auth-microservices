package auth

import (
	"context"
	"errors"
	"fmt"
	"time"

	"github.com/xh3sh/go-auth-microservices/internal/models"
	"github.com/xh3sh/go-auth-microservices/internal/repository"
	"github.com/xh3sh/go-auth-microservices/internal/utils"
)

// SessionService СѓРїСЂР°РІР»СЏРµС‚ Р¶РёР·РЅРµРЅРЅС‹Рј С†РёРєР»РѕРј СЃРµСЃСЃРёР№ РїРѕР»СЊР·РѕРІР°С‚РµР»РµР№
type SessionService struct {
	sessionRepo repository.SessionRepository
}

func NewSessionService(sessionRepo repository.SessionRepository) *SessionService {
	return &SessionService{
		sessionRepo: sessionRepo,
	}
}

// CreateSession СЃРѕР·РґР°РµС‚ РЅРѕРІСѓСЋ СЃРµСЃСЃРёСЋ РґР»СЏ РїРѕР»СЊР·РѕРІР°С‚РµР»СЏ
func (s *SessionService) CreateSession(ctx context.Context, userID string, ttl time.Duration) (*models.Session, error) {
	if userID == "" {
		return nil, errors.New("invalid user data")
	}

	sessionID, err := utils.GenerateRandomString(32)
	if err != nil {
		return nil, fmt.Errorf("failed to generate session ID: %w", err)
	}

	now := time.Now()
	session := &models.Session{
		SessionID: sessionID,
		UserID:    userID,
		CreatedAt: now,
		ExpiresAt: now.Add(ttl),
		IsActive:  true,
	}

	err = s.sessionRepo.SetSession(ctx, sessionID, session, ttl)
	if err != nil {
		return nil, fmt.Errorf("failed to save session: %w", err)
	}

	return session, nil
}

// GetSession РёР·РІР»РµРєР°РµС‚ РёРЅС„РѕСЂРјР°С†РёСЋ Рѕ СЃРµСЃСЃРёРё РїРѕ РµС‘ ID
func (s *SessionService) GetSession(ctx context.Context, sessionID string) (*models.Session, error) {
	if len(sessionID) == 0 {
		return nil, errors.New("invalid session ID")
	}

	var session models.Session
	err := s.sessionRepo.GetSession(ctx, sessionID, &session)
	if err != nil {
		return nil, fmt.Errorf("failed to get session: %w", err)
	}

	if !session.IsActive || time.Now().After(session.ExpiresAt) {
		return nil, errors.New("session is expired or inactive")
	}

	return &session, nil
}

// ValidateSession РїСЂРѕРІРµСЂСЏРµС‚ РІР°Р»РёРґРЅРѕСЃС‚СЊ СЃРµСЃСЃРёРё
func (s *SessionService) ValidateSession(ctx context.Context, sessionID string) (*models.Session, error) {
	session, err := s.GetSession(ctx, sessionID)
	if err != nil {
		return nil, err
	}

	if session == nil || !session.IsActive || time.Now().After(session.ExpiresAt) {
		return nil, errors.New("invalid session")
	}

	return session, nil
}

// RevokeSession СѓРґР°Р»СЏРµС‚ СЃРµСЃСЃРёСЋ (РѕС‚Р·С‹РІР°РµС‚ РµС‘)
func (s *SessionService) RevokeSession(ctx context.Context, sessionID string) error {
	if len(sessionID) == 0 {
		return errors.New("invalid session ID")
	}

	err := s.sessionRepo.DeleteSession(ctx, sessionID)
	if err != nil {
		return fmt.Errorf("failed to revoke session: %w", err)
	}

	return nil
}

// ExtendSession РїСЂРѕРґР»РµРІР°РµС‚ СЃСЂРѕРє РґРµР№СЃС‚РІРёСЏ СЃРµСЃСЃРёРё
func (s *SessionService) ExtendSession(ctx context.Context, sessionID string, newTTL time.Duration) error {
	session, err := s.GetSession(ctx, sessionID)
	if err != nil {
		return err
	}

	session.ExpiresAt = time.Now().Add(newTTL)

	err = s.sessionRepo.SetSession(ctx, sessionID, session, newTTL)
	if err != nil {
		return fmt.Errorf("failed to extend session: %w", err)
	}

	return nil
}
