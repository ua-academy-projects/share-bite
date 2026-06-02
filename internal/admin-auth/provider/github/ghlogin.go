package github

import (
	"context"
	"crypto/hmac"
	"crypto/rand"
	"crypto/sha256"
	"encoding/base64"
	"encoding/json"
	"errors"
	"html/template"
	"net/http"
	"net/url"
	"strings"
	"time"

	"github.com/ua-academy-projects/share-bite/pkg/database"
)

type Config struct {
	ClientID           string
	ClientSecret       string
	RedirectURL        string
	SuccessRedirectURL string
}

type Handler struct {
	cfg     Config
	gh      *githubClient
	users   UserRepo
	session SessionStore
	txm    database.TxManager
	secret []byte
}

type statePayload struct {
	Nonce string `json:"nonce"`
	Exp   int64  `json:"exp"`
}

func NewHandler(cfg Config, users UserRepo, session SessionStore, txm database.TxManager) *Handler {
	return &Handler{
		cfg: cfg,
		gh: &githubClient{
			httpClient:   &http.Client{Timeout: 10 * time.Second},
			clientID:     cfg.ClientID,
			clientSecret: cfg.ClientSecret,
		},
		users:   users,
		session: session,
		txm:     txm,
	}
}

func (h *Handler) Login(w http.ResponseWriter, r *http.Request) {
	state, err := h.generateState()
	if err != nil {
		http.Error(w, "internal error while Login", http.StatusInternalServerError)
		return
	}

	http.SetCookie(w, &http.Cookie{
		Name:     "oauth_state",
		Value:    state,
		HttpOnly: true,
		Secure:   false,
		SameSite: http.SameSiteLaxMode,
		MaxAge:   300, // 5 min
    Path: "/",
	})

  params := url.Values{
    "client_id": {h.cfg.ClientID},
    "redirect_uri": {h.cfg.RedirectURL},
    "state": {state},
    "scope": {"user:email"},
  }

	authURL := "https://github.com/login/oauth/authorize?" + params.Encode()
	http.Redirect(w, r, authURL, http.StatusTemporaryRedirect)
}

func (h *Handler) Callback(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()
	code := r.URL.Query().Get("code")
	state := r.URL.Query().Get("state")

	cookie, err := r.Cookie("oauth_state")
	if err != nil || cookie.Value != r.URL.Query().Get("state") {
		http.Error(w, "invalid CSRF state", http.StatusBadRequest)
		return
	}

  if err := h.validateState(state); err != nil {
		http.Error(w, "invalid state", http.StatusBadRequest)
		return
	}

	http.SetCookie(w, &http.Cookie{
	  Name:     "oauth_state",
		MaxAge:   -1,
		HttpOnly: true,
		Secure:   true,
		SameSite: http.SameSiteLaxMode,
		Path:     "/",
	})

	if code == "" {
		http.Error(w, "missing code", http.StatusBadRequest)
		return
	}

	accessToken, err := h.gh.exchangeCode(ctx, code)
	if err != nil {
		http.Error(w, "github auth failed in exchange code on token", http.StatusBadGateway)
		return
	}

	ghUser, err := h.gh.getUser(ctx, accessToken)
	if err != nil {
		http.Error(w, "failed to get github user", http.StatusBadGateway)
		return
	}
	if ghUser.Email == "" {
		primaryEmail, err := h.gh.getPrimaryEmail(ctx, accessToken)
		if err != nil {
			http.Error(w, "failed to get github user email", http.StatusBadGateway)
			return
		}
		ghUser.Email = primaryEmail
	}

	role := "user"
	var userID string
	if err := h.txm.ReadCommitted(ctx, func(txCtx context.Context) error {
		user, err := h.users.UpsertByGitHubID(txCtx, *ghUser)
		if err != nil {
			return err
		}
		userID = user.ID

		userWithRole, err := h.users.FindByID(txCtx, userID)
		if err != nil {
			return err
		}
		if userWithRole != nil && userWithRole.RoleSlug != "" {
			role = userWithRole.RoleSlug
		}
		return nil
	}); err != nil {
		http.Error(w, "db error", http.StatusInternalServerError)
		return
	}

	token, err := h.session.Create(ctx, userID, role)
	if err != nil {
		http.Error(w, "session error", http.StatusInternalServerError)
		return
	}

	http.SetCookie(w, &http.Cookie{
		Name:     "session",
		Value:    token,
		HttpOnly: false,
		Secure:   false,
		SameSite: http.SameSiteLaxMode,
		Path:     "/",
	})

	redirectURL := h.cfg.SuccessRedirectURL
	if redirectURL == "" {
		redirectURL = "/auth/github/success"
	}

	http.Redirect(w, r, redirectURL, http.StatusSeeOther)
}

func (h *Handler) Success(w http.ResponseWriter, r *http.Request) {
	t, err := template.ParseFiles("internal/admin-auth/provider/github/templates/success-tmp.html")
  if err != nil{
    http.Error(w, "Template error", http.StatusInternalServerError)
  }
  w.Header().Set("Context-Type", "text/html; charser=utf-8")
  t.Execute(w, nil)
}

func (h *Handler) generateState() (string, error) {

	nonceBytes := make([]byte, 16)
	if _, err := rand.Read(nonceBytes); err != nil {
		return "", err
	}
	payload := statePayload{
		Nonce: base64.URLEncoding.EncodeToString(nonceBytes),
		Exp:   time.Now().Add(5 * time.Minute).Unix(),
	}
	data, err := json.Marshal(payload)
	if err != nil {
		return "", err
	}
	dataB64 := base64.URLEncoding.EncodeToString(data)
	mac := hmac.New(sha256.New, h.secret)
	mac.Write([]byte(dataB64))
	signature := mac.Sum(nil)
	sigB64 := base64.URLEncoding.EncodeToString(signature)
	return dataB64 + "." + sigB64, nil

}

func (h *Handler) validateState(state string) error {

	parts := strings.Split(state, ".")
	if len(parts) != 2 {
		return errors.New("invalid state format")
	}
	dataB64 := parts[0]
	sigB64 := parts[1]

	mac := hmac.New(sha256.New, h.secret)
	mac.Write([]byte(dataB64))
	expectedSig := mac.Sum(nil)
	actualSig, err := base64.URLEncoding.DecodeString(sigB64)
	if err != nil {
		return err
	}
	if !hmac.Equal(expectedSig, actualSig) {
		return errors.New("invalid signature")
	}

	data, err := base64.URLEncoding.DecodeString(dataB64)
	if err != nil {
		return err
	}
	var payload statePayload
	if err := json.Unmarshal(data, &payload); err != nil {
		return err
	}

	if time.Now().Unix() > payload.Exp {
		return errors.New("state expired")
	}
	return nil

}