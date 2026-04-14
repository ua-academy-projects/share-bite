package ghAuth

import (
	"crypto/rand"
	"encoding/hex"
	"errors"
	"fmt"
	"net/http"
	"time"
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
}

func NewHandler(cfg Config, users UserRepo, session SessionStore) *Handler {
	return &Handler{
		cfg: cfg,
		gh: &githubClient{
			httpClient:   &http.Client{Timeout: 10 * time.Second},
			clientID:     cfg.ClientID,
			clientSecret: cfg.ClientSecret,
		},
		users:   users,
		session: session,
	}
}

func (h *Handler) Login(w http.ResponseWriter, r *http.Request) {
	state, err := generateState()
	if err != nil {
		http.Error(w, "internal error while Login", http.StatusInternalServerError)
		return
	}

	http.SetCookie(w, &http.Cookie{
		Name:     "oauth_state",
		Value:    state,
		HttpOnly: true,
		Secure:   true,
		SameSite: http.SameSiteLaxMode,
		MaxAge:   300, // 5 min
	})

	authURL := fmt.Sprintf(
		"https://github.com/login/oauth/authorize?client_id=%s&redirect_uri=%s&state=%s&scope=user:email",
		h.cfg.ClientID, h.cfg.RedirectURL, state,
	)
	http.Redirect(w, r, authURL, http.StatusTemporaryRedirect)
}

func (h *Handler) Callback(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()

	cookie, err := r.Cookie("oauth_state")
	if err != nil || cookie.Value != r.URL.Query().Get("state") {
		http.Error(w, "invalid CSRF state", http.StatusBadRequest)
		return
	}

	http.SetCookie(w, &http.Cookie{Name: "oauth_state", MaxAge: -1})

	code := r.URL.Query().Get("code")
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

	user, err := h.users.UpsertByGitHubID(ctx, *ghUser)
	if err != nil {
		http.Error(w, "db error", http.StatusInternalServerError)
		return
	}

	token, err := h.session.Create(ctx, fmt.Sprintf("%d", user.ID))
	if err != nil {
		http.Error(w, "session error", http.StatusInternalServerError)
		return
	}

	http.SetCookie(w, &http.Cookie{
		Name:     "session",
		Value:    token,
		HttpOnly: true,
		Secure:   true,
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
	w.Header().Set("Content-Type", "text/html; charset=utf-8")
	w.WriteHeader(http.StatusOK)
	_, _ = w.Write([]byte(`<!DOCTYPE html>
<html lang="uk">
<head>
  <meta charset="UTF-8">
  <meta name="viewport" content="width=device-width, initial-scale=1.0">
  <title>Threads Clone</title>
  <style>
    body {
      margin: 0;
      font-family: system-ui, -apple-system, BlinkMacSystemFont, "Segoe UI", Roboto, sans-serif;
      background-color: #0d0d0d;
      color: #fff;
      display: flex;
      justify-content: center;
      height: 100vh;
    }

    .container {
      display: flex;
      width: 100%;
      max-width: 1100px;
      margin-top: 20px;
    }

    /* Бічна панель */
    .sidebar {
      width: 80px;
      display: flex;
      flex-direction: column;
      align-items: center;
      gap: 20px;
      padding: 20px 0;
      border-right: 1px solid #2a2a2a;
    }

    .sidebar div {
      width: 40px;
      height: 40px;
      background-color: #181818;
      border-radius: 10px;
      display: flex;
      align-items: center;
      justify-content: center;
      color: #aaa;
      cursor: pointer;
      transition: all 0.2s ease;
    }

    .sidebar div:hover {
      background-color: #2e2e2e;
      color: #fff;
    }

    /* Основна стрічка */
    .feed {
      flex: 1;
      padding: 20px;
    }

    .post {
      background-color: #121212;
      border-radius: 12px;
      padding: 16px;
      margin-bottom: 16px;
      border: 1px solid #1f1f1f;
    }

    .user {
      font-weight: 600;
      color: #ddd;
    }

    .time {
      color: #777;
      font-size: 14px;
      margin-left: 6px;
    }

    .content {
      margin-top: 8px;
      line-height: 1.5;
    }

    .actions {
      display: flex;
      gap: 20px;
      margin-top: 10px;
      color: #999;
      font-size: 14px;
    }

    .actions span {
      cursor: pointer;
    }

    .actions span:hover {
      color: #fff;
    }

    /* Права панель логіну */
    .login-box {
      width: 360px;
      background-color: #181818;
      border: 1px solid #2a2a2a;
      border-radius: 16px;
      padding: 20px;
      margin-top: 20px;
      height: fit-content;
    }

    .login-box h2 {
      text-align: center;
      font-size: 18px;
    }

    .button {
      display: block;
      background-color: #fff;
      color: #000;
      text-align: center;
      padding: 10px;
      border-radius: 10px;
      margin: 20px auto;
      width: 80%;
      font-weight: 600;
      text-decoration: none;
    }

    .button:hover {
      background-color: #e5e5e5;
    }

    footer {
      text-align: center;
      font-size: 12px;
      color: #777;
      margin-top: 20px;
    }
  </style>
</head>
<body>
  <div class="container">
    <div class="sidebar">
      <div>🏠</div>
      <div>🔍</div>
      <div>➕</div>
      <div>❤️</div>
      <div>👤</div>
      <div>⚙️</div>
    </div>

    <div class="feed">
      <div class="post">
        <div><span class="user">ryan_mavity</span><span class="time">17 год</span></div>
        <div class="content">Can’t think of a better way, if this is really it, for Ovi to go out in Pittsburgh. Team gets the W, the young fellas like Little Pro & Leno carry the load & Ovi gets one last goal into the empty net to say goodbye.</div>
        <div class="actions"><span>❤️ 7</span><span>💬</span><span>🔁</span><span>▶️</span></div>
      </div>

      <div class="post">
        <div><span class="user">j.o.s.e.l.y.n.r</span><span class="time">16 год</span></div>
        <div class="content">Outgrowing your past self is a sign you’re no longer surviving, you’re choosing yourself.</div>
        <div class="actions"><span>❤️ 32</span><span>💬 2</span><span>🔁</span><span>▶️</span></div>
      </div>

      <div class="post">
        <div><span class="user">alasdair_gold✔️</span><span class="time">3 год</span></div>
        <div class="content">What do you make of Roberto De Zerbi's first Tottenham starting XI?<br>
        <a href="#" style="color:#3b82f6;">football.london/totte...</a></div>
        <div class="actions"><span>❤️</span><span>💬</span><span>🔁</span><span>▶️</span></div>
      </div>
    </div>

    <div class="login-box">
      <h2>Welcome to Share-Bite</h2>
      <p style="text-align:center;color:#bbb;">Діліться хавкою, де нормально можна посидіти і почілити.</p>
      <a href="#" class="button">Продовжити з Share Bite</a>
      <p style="text-align:center;font-size:13px;color:#888;">Увійти натомість за допомогою імені користувача</p>
      <footer>
        © 2026 Share-bite<br>
        Політика конфіденційності · Умови · Cookies
      </footer>
    </div>
  </div>
</body>
</html>
`))
}

func generateState() (string, error) {
	b := make([]byte, 16)
	if _, err := rand.Read(b); err != nil {
		return "", errors.New("ghlogin: failed to generate state")
	}
	return hex.EncodeToString(b), nil
}
