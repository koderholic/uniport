package api

import (
	"net/http"
	"strings"
	"time"

	"github.com/julienschmidt/httprouter"
	"github.com/justinas/alice"

	"uniport/config"
	"uniport/utils"

	jwt "github.com/dgrijalva/jwt-go"
)

func verifyID(httpRes http.ResponseWriter, httpReq *http.Request, claims jwt.MapClaims) {
	if claims == nil {
		http.Redirect(httpRes, httpReq, "/", http.StatusTemporaryRedirect)
		return
	}

	if claims["ID"] == nil {
		http.Redirect(httpRes, httpReq, "/", http.StatusTemporaryRedirect)
		return
	}
}

func pageHandler(middlewares alice.Chain, router *Router) {

	router.GET("/login", func(httpRes http.ResponseWriter, httpReq *http.Request, httpParams httprouter.Params) {
		if claims := utils.VerifyJWT(httpRes, httpReq); claims != nil {
			if uint64(claims["ID"].(float64)) > 0 {
				switch {
				case claims["IsAdmin"] != nil && claims["IsAdmin"].(bool):
					http.Redirect(httpRes, httpReq, "/admin", http.StatusTemporaryRedirect)
					break

				default:
					http.Redirect(httpRes, httpReq, "/dashboard", http.StatusTemporaryRedirect)
					break
				}
			}
		}
		fileServe(httpRes, httpReq)
	})

	//Authenticated Pages --> Below
	router.GET("/admin/*page", func(httpRes http.ResponseWriter, httpReq *http.Request, httpParams httprouter.Params) {
		claims := utils.VerifyJWT(httpRes, httpReq)
		verifyID(httpRes, httpReq, claims)
		if claims["IsAdmin"] == nil || !claims["IsAdmin"].(bool) {
			http.Redirect(httpRes, httpReq, "/", http.StatusTemporaryRedirect)
			return
		}
		fileServe(httpRes, httpReq)
	})

	router.GET("/dashboard/*page", func(httpRes http.ResponseWriter, httpReq *http.Request, httpParams httprouter.Params) {
		claims := utils.VerifyJWT(httpRes, httpReq)
		verifyID(httpRes, httpReq, claims)
		fileServe(httpRes, httpReq)
	})
	//Authenticated Pages --> Below

	router.NotFound = middlewares.ThenFunc(func(httpRes http.ResponseWriter, httpReq *http.Request) {
		frontend := strings.Split(httpReq.URL.Path[1:], "/")
		switch frontend[0] {
		case "logout":
			cookieMonster := &http.Cookie{
				Name: config.Get().COOKIE, Value: "deleted", Path: "/",
				Expires: time.Now().Add(-(time.Hour * 24 * 30 * 12)), // set the expire time
			}
			http.SetCookie(httpRes, cookieMonster)
			httpReq.AddCookie(cookieMonster)
			http.Redirect(httpRes, httpReq, "/",
				http.StatusTemporaryRedirect)

		default:
			fileServe(httpRes, httpReq)
		}
	})

}

type pageStruct struct{ Bundle, Version, Page, SeoTitle, SeoContent string }
