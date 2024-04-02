package webapi

import (
	"errors"
	"net/http"
	"time"

	"github.com/go-chi/render"

	"github.com/pubgolf/pubgolf/api/internal/lib/config"
)

const (
	authCookieName = "web_admin_user_token"
)

var errUserAuth = errors.New("missing or invalid auth credential")

type logInRequestBody struct {
	Password string `json:"password"`
}

type logInResponseBody struct {
	Success bool `json:"success"`
}

func logIn(cfg *config.App) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		ctx := r.Context()

		req := logInRequestBody{}
		if ok := guardParseJSONRequest(ctx, parseJSONRequest(w, r, &req), w, r); !ok {
			return
		}

		if req.Password != cfg.AdminAuth.Password {
			newErrorResponse(ctx, errorCodeNotAuthorized, "Incorrect password", errUserAuth).Render(w, r)

			return
		}

		http.SetCookie(w, &http.Cookie{
			Name:     authCookieName,
			Value:    cfg.AdminAuth.CookieToken,
			Domain:   cfg.HostOrigin,
			HttpOnly: true,
			Secure:   true,
			SameSite: http.SameSiteStrictMode,
			Expires:  time.Now().Add(30 * 24 * time.Hour),
		})

		render.JSON(w, r, logInResponseBody{true})
	}
}

type logOutResponseBody struct {
	Success bool `json:"success"`
}

func logOut(cfg *config.App) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		ctx := r.Context()

		c, err := r.Cookie(authCookieName)
		if err != nil {
			if errors.Is(err, http.ErrNoCookie) {
				newErrorResponse(ctx, errorCodeNotAuthorized, "Missing auth cookie", errUserAuth).Render(w, r)

				return
			}

			newErrorResponse(ctx, errorCodeGenericNonRetryable, "Could not read auth cookie", err).Render(w, r)

			return
		}

		if c.Value != cfg.AdminAuth.CookieToken {
			newErrorResponse(ctx, errorCodeNotAuthorized, "Invalid auth cookie", errUserAuth).Render(w, r)

			return
		}

		http.SetCookie(w, &http.Cookie{
			Name:     authCookieName,
			Value:    "",
			Domain:   cfg.HostOrigin,
			HttpOnly: true,
			Secure:   true,
			SameSite: http.SameSiteStrictMode,

			MaxAge: -1,
		})

		render.JSON(w, r, logInResponseBody{true})
	}
}

type tokenResponseBody struct {
	Token string `json:"token"`
}

func getAPIToken(cfg *config.App) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		ctx := r.Context()

		c, err := r.Cookie(authCookieName)
		if err != nil {
			if errors.Is(err, http.ErrNoCookie) {
				newErrorResponse(ctx, errorCodeNotAuthorized, "Missing auth cookie", errUserAuth).Render(w, r)

				return
			}

			newErrorResponse(ctx, errorCodeGenericNonRetryable, "Could not read auth cookie", err).Render(w, r)

			return
		}

		if c.Value != cfg.AdminAuth.CookieToken {
			newErrorResponse(ctx, errorCodeNotAuthorized, "Invalid auth cookie", errUserAuth).Render(w, r)

			return
		}

		render.JSON(w, r, tokenResponseBody{cfg.AdminAuth.AdminServiceToken})
	}
}
