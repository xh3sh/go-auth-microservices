package handler

import (
	"context"
	"crypto/rand"
	"encoding/base64"
	"encoding/json"
	"github.com/xh3sh/go-auth-microservices/internal/models"
	"github.com/xh3sh/go-auth-microservices/internal/utils"
	"net/http"
	"net/url"
	"strings"
	"time"

	"github.com/gin-gonic/gin"
)

func generateState() string {
	b := make([]byte, 32)
	rand.Read(b)
	return base64.URLEncoding.EncodeToString(b)
}

// HandleOAuthToken РѕР±СЂР°Р±Р°С‚С‹РІР°РµС‚ Р·Р°РїСЂРѕСЃ РЅР° РїРѕР»СѓС‡РµРЅРёРµ С‚РѕРєРµРЅР° OAuth
func (h *Handler) HandleOAuthToken(c *gin.Context) {
	var req models.OAuthTokenRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid request"})
		return
	}

	switch req.GrantType {
	case "client_credentials":
		token, err := h.oauthService.GenerateOAuthToken(req.ClientID, req.Scope)
		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to generate token"})
			return
		}
		c.JSON(http.StatusOK, token)

	case "authorization_code":
		token, err := h.oauthService.ExchangeAuthorizationCode(req.ClientID, req.ClientSecret, req.Code, req.RedirectURI)
		if err != nil {
			c.JSON(http.StatusUnauthorized, gin.H{"error": "Failed to exchange code"})
			return
		}
		c.JSON(http.StatusOK, token)

	default:
		c.JSON(http.StatusBadRequest, gin.H{"error": "Unsupported grant type"})
	}
}

// HandleOAuthValidate РїСЂРѕРІРµСЂСЏРµС‚ РІР°Р»РёРґРЅРѕСЃС‚СЊ OAuth СЃРµСЃСЃРёРё
func (h *Handler) HandleOAuthValidate(c *gin.Context) {
	token, err := c.Cookie("oauth_access_token")
	if err != nil {
		token = c.GetHeader("Authorization")
		if len(token) > 7 && strings.HasPrefix(strings.ToUpper(token), "BEARER ") {
			token = token[7:]
		}
	}

	if token == "" {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "OAuth session not found"})
		return
	}

	claims, err := h.jwtService.ExtractClaims(token)
	if err != nil {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "Invalid or expired OAuth token"})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"valid":    true,
		"user":     claims.Username,
		"user_id":  claims.UserID,
		"provider": "google",
	})
}

// HandleOAuthLogout Р·Р°РІРµСЂС€Р°РµС‚ OAuth СЃРµСЃСЃРёСЋ
func (h *Handler) HandleOAuthLogout(c *gin.Context) {
	if token, err := c.Cookie("oauth_access_token"); err == nil && token != "" {
		if claims, err := h.jwtService.ExtractClaims(token); err == nil {
			c.Set("user_id", claims.UserID)
		}
	}

	c.Set("auth_method", "oauth")
	c.SetCookie("oauth_access_token", "", -1, "/", "", false, true)
	c.SetCookie("oauth_refresh_token", "", -1, "/", "", false, true)

	c.JSON(http.StatusOK, gin.H{"message": "OAuth session cleared"})
}

// HandleOAuthAuthorize РёРЅРёС†РёРёСЂСѓРµС‚ РїСЂРѕС†РµСЃСЃ Р°РІС‚РѕСЂРёР·Р°С†РёРё С‡РµСЂРµР· РІРЅРµС€РЅРµРіРѕ РїСЂРѕРІР°Р№РґРµСЂР°
func (h *Handler) HandleOAuthAuthorize(c *gin.Context) {
	provider := c.Query("provider")
	if provider == "google" {
		if h.config.GoogleClientID == "" {
			c.JSON(http.StatusInternalServerError, gin.H{"error": "Google OAuth not configured"})
			return
		}

		state := generateState()
		c.SetCookie("oauth_state", state, 300, "/auth/oauth", "", false, true)

		u := "https://accounts.google.com/o/oauth2/v2/auth"
		q := url.Values{}
		q.Set("client_id", h.config.GoogleClientID)
		q.Set("redirect_uri", h.config.GoogleRedirectURL)
		q.Set("response_type", "code")
		q.Set("scope", "https://www.googleapis.com/auth/userinfo.email https://www.googleapis.com/auth/userinfo.profile")
		q.Set("state", state)

		c.Redirect(http.StatusFound, u+"?"+q.Encode())
		return
	}

	clientID := c.Query("client_id")
	redirectURI := c.Query("redirect_uri")
	responseType := c.Query("response_type")
	state := c.Query("state")

	if clientID == "" || redirectURI == "" || responseType != "code" {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid authorization request"})
		return
	}

	code, _ := utils.GenerateRandomString(16)
	target := redirectURI + "?code=" + code
	if state != "" {
		target += "&state=" + state
	}

	c.Redirect(http.StatusFound, target)
}

// HandleGoogleCallback РѕР±СЂР°Р±Р°С‚С‹РІР°РµС‚ РѕС‚РІРµС‚ РѕС‚ Google OAuth
func (h *Handler) HandleGoogleCallback(c *gin.Context) {
	state := c.Query("state")
	cookieState, err := c.Cookie("oauth_state")
	if err != nil || state != cookieState {
		c.JSON(http.StatusForbidden, gin.H{"error": "Invalid OAuth state (CSRF detected)"})
		return
	}
	c.SetCookie("oauth_state", "", -1, "/auth/oauth", "", false, true)

	code := c.Query("code")
	if code == "" {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Code missing"})
		return
	}

	client := &http.Client{Timeout: 10 * time.Second}
	ctx, cancel := context.WithTimeout(c.Request.Context(), 10*time.Second)
	defer cancel()

	resp, err := client.PostForm("https://oauth2.googleapis.com/token", url.Values{
		"code":          {code},
		"client_id":     {h.config.GoogleClientID},
		"client_secret": {h.config.GoogleClientSecret},
		"redirect_uri":  {h.config.GoogleRedirectURL},
		"grant_type":    {"authorization_code"},
	})
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to connect to Google"})
		return
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "Google rejected the authorization code"})
		return
	}

	var tokenData struct {
		AccessToken string `json:"access_token"`
	}
	if err := json.NewDecoder(resp.Body).Decode(&tokenData); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to parse Google response"})
		return
	}

	req, _ := http.NewRequestWithContext(ctx, "GET", "https://www.googleapis.com/oauth2/v2/userinfo", nil)
	req.Header.Set("Authorization", "Bearer "+tokenData.AccessToken)
	userResp, err := client.Do(req)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to get user info from Google"})
		return
	}
	defer userResp.Body.Close()

	if userResp.StatusCode != http.StatusOK {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "Failed to authorize user profile"})
		return
	}

	var googleUser struct {
		ID    string `json:"id"`
		Email string `json:"email"`
	}
	if err := json.NewDecoder(userResp.Body).Decode(&googleUser); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to parse profile data"})
		return
	}

	username := googleUser.Email
	creds, err := h.repo.GetUserByUsername(c.Request.Context(), username)
	if err != nil {
		id, _ := utils.GenerateRandomString(8)
		creds = models.UserCredentials{
			ID:           id,
			Username:     username,
			PasswordHash: "oauth_google_protected_" + googleUser.ID,
		}
		_ = h.repo.CreateUserCredentials(c.Request.Context(), creds)
	}

	tokenPair, err := h.jwtService.GenerateTokenPair(creds.ID, creds.Username)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to issue session"})
		return
	}

	c.SetCookie("oauth_access_token", tokenPair.AccessToken, int(h.config.JWTExpiration), "/", "", false, true)
	c.SetCookie("oauth_refresh_token", tokenPair.RefreshToken, int(h.config.JWTRefreshExpiration), "/", "", false, true)

	c.Redirect(http.StatusFound, "/")
}
