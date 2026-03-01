package models

import "time"

// JWT СЃС‚СЂСѓРєС‚СѓСЂС‹
type JWTClaims struct {
	UserID   string `json:"user_id"`
	Username string `json:"username"`
	Exp      int64  `json:"exp"`
	Iat      int64  `json:"iat"`
}

type JWTTokenPair struct {
	AccessToken  string `json:"access_token"`
	RefreshToken string `json:"refresh_token"`
	ExpiresIn    int64  `json:"expires_in"`
	TokenType    string `json:"token_type"`
}

// API Key СЃС‚СЂСѓРєС‚СѓСЂС‹
type APIKey struct {
	ID        string    `json:"id"`
	Key       string    `json:"key"`
	UserID    string    `json:"user_id"`
	Name      string    `json:"name"`
	CreatedAt time.Time `json:"created_at"`
	ExpiresAt time.Time `json:"expires_at"`
	IsActive  bool      `json:"is_active"`
}

type APIKeyResponse struct {
	ID        string    `json:"id"`
	Key       string    `json:"key"`
	UserID    string    `json:"user_id"`
	Name      string    `json:"name"`
	CreatedAt time.Time `json:"created_at"`
	ExpiresAt time.Time `json:"expires_at"`
}

// OAuth 2.0 СЃС‚СЂСѓРєС‚СѓСЂС‹
type OAuthClient struct {
	ClientID     string    `json:"client_id"`
	ClientSecret string    `json:"client_secret"`
	RedirectURI  string    `json:"redirect_uri"`
	Scope        []string  `json:"scope"`
	UserID       string    `json:"user_id"`
	IsActive     bool      `json:"is_active"`
	CreatedAt    time.Time `json:"created_at"`
}

type OAuthToken struct {
	AccessToken  string `json:"access_token"`
	TokenType    string `json:"token_type"`
	ExpiresIn    int64  `json:"expires_in"`
	RefreshToken string `json:"refresh_token,omitempty"`
	Scope        string `json:"scope"`
	IssuedAt     int64  `json:"issued_at"`
}

// Session СЃС‚СЂСѓРєС‚СѓСЂР°
type Session struct {
	SessionID string    `json:"session_id"`
	UserID    string    `json:"user_id"`
	CreatedAt time.Time `json:"created_at"`
	ExpiresAt time.Time `json:"expires_at"`
	IsActive  bool      `json:"is_active"`
}

// РЎС‚СЂСѓРєС‚СѓСЂС‹ Р·Р°РїСЂРѕСЃРѕРІ
type LoginRequest struct {
	Username string `json:"username" form:"username"`
	Password string `json:"password" form:"password"`
}

type RegisterRequest struct {
	Username string `json:"username" form:"username"`
	Password string `json:"password" form:"password"`
}

type RefreshTokenRequest struct {
	RefreshToken string `json:"refresh_token"`
}

type OAuthTokenRequest struct {
	ClientID     string `json:"client_id"`
	ClientSecret string `json:"client_secret"`
	GrantType    string `json:"grant_type"`
	Code         string `json:"code,omitempty"`
	RedirectURI  string `json:"redirect_uri,omitempty"`
	Scope        string `json:"scope,omitempty"`
}

// РЎС‚СЂСѓРєС‚СѓСЂР° РѕС€РёР±РєРё
type AuthError struct {
	Code    string `json:"code"`
	Message string `json:"message"`
	Details string `json:"details,omitempty"`
}

// UserCredentials С…СЂР°РЅРёС‚ РґР°РЅРЅС‹Рµ РїРѕР»СЊР·РѕРІР°С‚РµР»СЏ РґР»СЏ Р°СѓС‚РµРЅС‚РёС„РёРєР°С†РёРё
type UserCredentials struct {
	ID           string `json:"id"`
	Username     string `json:"username"`
	PasswordHash string `json:"password_hash"`
}
