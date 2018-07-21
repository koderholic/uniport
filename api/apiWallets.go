package api

import (
	"encoding/json"
	"fmt"
	"net/http"
	"regexp"
	"strconv"
	"strings"
	"time"

	"github.com/coreos/bbolt"
	"github.com/justinas/alice"
	"github.com/timshannon/bolthold"

	"uniport/buckets"
	"uniport/config"
	"uniport/utils"
)

type apiWalletsStruct struct {
	ID       uint64
	Workflow string
	Createdate,
	Updatedate time.Time
	Path uint32

	Label, Mnemonic,
	Password, Description string
}

func apiHandlerWallets(middlewares alice.Chain, router *Router) {
	router.Get("/api/wallets", middlewares.ThenFunc(apiWalletGet))
	router.Post("/api/wallets", middlewares.ThenFunc(apiWalletPost))
	router.Get("/api/wallets/search", middlewares.ThenFunc(apiWalletsSearch))
}

func apiWalletsSearch(httpRes http.ResponseWriter, httpReq *http.Request) {

	httpRes.Header().Set("Content-Type", "application/json")

	var statusBody interface{}
	statusCode := http.StatusInternalServerError
	statusMessage := ""

	if claims := utils.VerifyJWT(httpRes, httpReq); claims == nil {
		statusBody = map[string]string{"Redirect": "/"}
	} else {

		searchResults := []buckets.Wallets{}
		formSearch := new(apiSearch)
		formSearch.Skip, _ = strconv.Atoi(httpReq.FormValue("skip"))
		formSearch.Limit, _ = strconv.Atoi(httpReq.FormValue("limit"))

		formSearch.Text = strings.TrimSpace(httpReq.FormValue("search"))
		formSearch.Field = strings.TrimSpace(httpReq.FormValue("field"))
		if formSearch.Text == "" {
			formSearch.Text = "."
		} else {
			formSearch.Text = regexp.QuoteMeta(formSearch.Text)
		}

		switch formSearch.Field {
		default:
			formSearch.Field = strings.Title(strings.ToLower(formSearch.Field))
		case "":
			formSearch.Field = "Label"
		}

		if err := config.Get().BoltHold.Bolt().View(func(tx *bolt.Tx) error {
			err := config.Get().BoltHold.Find(&searchResults,
				bolthold.Where(formSearch.Field).RegExp(
					regexp.MustCompile(`(?im)`+formSearch.Text)).SortBy("ID").Reverse().Limit(formSearch.Limit).Skip(formSearch.Skip),
			)
			return err
		}); err != nil {
			statusMessage = err.Error()
		} else {

			searchList := make([]apiSearchResult, len(searchResults))
			for pos, result := range searchResults {
				searchList[pos].ID = result.ID
				searchList[pos].Date = JSONTime(result.Updatedate)
				searchList[pos].Details = fmt.Sprintf("%v", result.Label)

			}

			statusCode = http.StatusOK
			statusBody = searchList
		}
	}

	json.NewEncoder(httpRes).Encode(Message{
		Code:    statusCode,
		Body:    statusBody,
		Message: statusMessage,
	})
}

func apiWalletGet(httpRes http.ResponseWriter, httpReq *http.Request) {

	httpRes.Header().Set("Content-Type", "application/json")

	var statusBody interface{}
	statusCode := http.StatusOK
	statusMessage := ""

	if claims := utils.VerifyJWT(httpRes, httpReq); claims == nil {
		statusBody = map[string]string{"Redirect": "/"}
	} else {

		sWalletID := strings.TrimSpace(httpReq.FormValue("id"))
		if sWalletID == "" {
			statusCode = http.StatusInternalServerError
			statusMessage = "Error Wallet ID is required to load form"
		} else {
			WalletID, _ := strconv.ParseUint(sWalletID, 0, 64)
			walletsList, err := buckets.Wallets{}.GetFieldValue("ID", WalletID)
			if err != nil {
				statusMessage = err.Error()
			} else {
				if len(walletsList) > 0 {

					statusBody = apiWalletsStruct{
						ID:         walletsList[0].ID,
						Createdate: walletsList[0].Createdate,
						Updatedate: walletsList[0].Updatedate,

						Path:        walletsList[0].Path,
						Label:       walletsList[0].Label,
						Mnemonic:    walletsList[0].Mnemonic,
						Description: walletsList[0].Description,
					}
				}
			}
		}
	}

	json.NewEncoder(httpRes).Encode(Message{
		Code:    statusCode,
		Body:    statusBody,
		Message: statusMessage,
	})
}

func apiWalletPost(httpRes http.ResponseWriter, httpReq *http.Request) {
	httpRes.Header().Set("Content-Type", "application/json")

	var statusBody interface{}
	statusCode := http.StatusInternalServerError
	statusMessage := ""

	if claims := utils.VerifyJWT(httpRes, httpReq); claims == nil {
		statusBody = map[string]string{"Redirect": "/"}
	} else {

		formStruct := buckets.Wallets{}
		err := json.NewDecoder(httpReq.Body).Decode(&formStruct)
		if err != nil {
			statusMessage = "Error Decoding Form Values: " + err.Error()
		} else {

			bucketWallet := buckets.Wallets{}

			if formStruct.ID != 0 {
				bucketWalletList, _ := buckets.Wallets{}.GetFieldValue("ID", formStruct.ID)
				if len(bucketWalletList) != 1 {
					statusMessage = "Error Decoding Form Values: " + err.Error()
				} else {
					bucketWallet = bucketWalletList[0]
				}
			} else {
				formStruct.Workflow = Enabled
			}

			bucketWallet.Path = formStruct.Path
			bucketWallet.Label = formStruct.Label
			bucketWallet.Description = formStruct.Description

			if statusMessage == "" {
				if bucketWallet.Label == "" {
					statusMessage += "Label is Required \n"
				}

				if strings.HasSuffix(statusMessage, "\n") {
					statusMessage = statusMessage[:len(statusMessage)-2]
				}
			}

			if statusMessage == "" {
				err = bucketWallet.Create(&bucketWallet)
				if err != nil {
					statusMessage = "Error Saving Record: " + err.Error()
				} else {
					statusCode = http.StatusOK
					statusMessage = RecordSaved
					statusBody = bucketWallet.ID
				}
			}
		}
	}

	json.NewEncoder(httpRes).Encode(Message{
		Code:    statusCode,
		Body:    statusBody,
		Message: statusMessage,
	})
}
