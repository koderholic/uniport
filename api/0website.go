package api

import (
	"encoding/json"
	"net/http"
	"strings"
	"time"

	"fmt"

	"golang.org/x/crypto/bcrypt"

	jwt "github.com/dgrijalva/jwt-go"
	"github.com/justinas/alice"

	"uniport/buckets"
	"uniport/config"
	"uniport/utils"
)

func apiHandlerWebsite(middlewares alice.Chain, router *Router) {
	router.Post("/api/login", middlewares.ThenFunc(apiLogin))
	router.Post("/api/signup", middlewares.ThenFunc(apiSignup))
}

func apiLogin(httpRes http.ResponseWriter, httpReq *http.Request) {
	httpRes.Header().Set("Content-Type", "application/json")

	var formStruct struct {
		Password string
	}

	statusBody := make(map[string]interface{})
	statusCode := http.StatusInternalServerError
	statusMessage := "Invalid Password"

	err := json.NewDecoder(httpReq.Body).Decode(&formStruct)
	if err == nil {
		wallets, _ := buckets.Wallets{}.GetFieldValue("ID", uint64(1))

		if len(wallets) == 1 {
			lValid := true
			Wallet := wallets[0]

			if err = bcrypt.CompareHashAndPassword(Wallet.Password, []byte(formStruct.Password)); err != nil {
				lValid = false
			}

			if Wallet.Workflow != "enabled" && Wallet.Workflow != "active" {
				lValid = false
			}

			if lValid {
				jwtClaims := jwt.MapClaims{}
				jwtClaims["ID"] = Wallet.ID
				jwtClaims["Label"] = Wallet.Label
				jwtClaims["Path"] = Wallet.Path

				statusBody["Redirect"] = "/dashboard"

				cookieExpires := time.Now().Add(time.Hour * 24 * 14) // set the expire time
				jwtClaims["exp"] = cookieExpires.Unix()

				if jwtToken, errJwt := utils.GenerateJWT(jwtClaims); errJwt == nil {
					cookieMonster := &http.Cookie{
						Name: config.Get().COOKIE, Value: jwtToken, Expires: cookieExpires, Path: "/",
					}
					http.SetCookie(httpRes, cookieMonster)
					httpReq.AddCookie(cookieMonster)

					statusCode = http.StatusOK
					statusMessage = "Password Verified"
				}
			}

		}
	} else {
		println(err.Error())
	}

	json.NewEncoder(httpRes).Encode(Message{
		Code:    statusCode,
		Message: statusMessage,
		Body:    statusBody,
	})
}

func apiSignup(httpRes http.ResponseWriter, httpReq *http.Request) {

	httpRes.Header().Set("Content-Type", "application/json")

	statusMessage := ""
	statusBody := make(map[string]interface{})
	statusCode := http.StatusInternalServerError

	var formStruct struct {
		Password, Confirm, Mnemonic string
	}

	err := json.NewDecoder(httpReq.Body).Decode(&formStruct)
	if err != nil {
		statusMessage = "Error Decoding Form Values " + err.Error()
	} else {
		wallets, err := buckets.Wallets{}.GetFieldValue("ID", uint64(1))
		if err != nil {
			statusMessage = fmt.Sprintf("Error Validating Walletname %s", err.Error())
		} else if len(wallets) > 0 {

			statusMessage = fmt.Sprintf("A wallet already exists")
			// if formStruct.Mnemonic != "" && formStruct.Mnemonic != wallets[0].Mnemonic {
			// 	statusMessage = fmt.Sprintf("A wallet already exists")
			// }

		} else {

			//All Seems Clear, Create New Wallet Now Now

			if formStruct.Password == "" {
				statusMessage += "Password " + IsRequired
			}

			if formStruct.Confirm == "" {
				statusMessage += "Confirm Password " + IsRequired
			}

			if statusMessage == "" {
				if formStruct.Password != formStruct.Confirm {
					statusMessage += "Passwords do not match "
				}
			}

			if strings.HasSuffix(statusMessage, "\n") {
				statusMessage = statusMessage[:len(statusMessage)-2]
			}

			if statusMessage == "" {
				bucketWallet := buckets.Wallets{}
				bucketWallet.Workflow = Enabled
				bucketWallet.Label = "uniport"
				bucketWallet.Mnemonic = formStruct.Mnemonic

				hash, _ := bcrypt.GenerateFromPassword([]byte(formStruct.Password), bcrypt.DefaultCost)
				bucketWallet.Password = hash

				statusCode = http.StatusOK
				statusMessage = "Sign up successful, please login"
				bucketWallet.Create(&bucketWallet)

			}
			//All Seems Clear, Create New Wallet Now Now

		}
	}

	json.NewEncoder(httpRes).Encode(Message{
		Code:    statusCode,
		Message: statusMessage,
		Body:    statusBody,
	})
	// //Send E-Mail
}
