package api

import (
	"encoding/json"
	"math/big"
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

type apiAccountStruct struct {
	ID       uint64
	Workflow string
	Createdate,
	Updatedate time.Time
	Path uint32

	Label, Address, Publickey,
	Privatekey string

	Balance  *big.Int
	WalletID uint64

	Wallet string
}

func apiHandlerAccounts(middlewares alice.Chain, router *Router) {
	router.Get("/api/accounts", middlewares.ThenFunc(apiAccountGet))
	router.Post("/api/accounts", middlewares.ThenFunc(apiAccountPost))
	router.Post("/api/accounts/search", middlewares.ThenFunc(apiAccountSearch))
}

func apiAccountSearch(httpRes http.ResponseWriter, httpReq *http.Request) {

	httpRes.Header().Set("Content-Type", "application/json")

	var statusBody interface{}
	statusCode := http.StatusInternalServerError
	statusMessage := ""

	if claims := utils.VerifyJWT(httpRes, httpReq); claims == nil {
		statusBody = map[string]string{"Redirect": "/"}
	} else {

		searchResults := []buckets.Accounts{}
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

			searchList := make([]apiAccountStruct, len(searchResults))
			for pos, result := range searchResults {
				searchList[pos].ID = result.ID
				searchList[pos].Workflow = result.Workflow
				searchList[pos].Createdate = result.Createdate
				searchList[pos].Updatedate = result.Updatedate

				searchList[pos].Path = result.Path
				searchList[pos].Label = result.Label
				searchList[pos].Address = result.Address
				searchList[pos].Balance = result.Balance

				searchList[pos].WalletID = result.WalletID
				if result.WalletID != 0 {
					listRes, _ := buckets.Wallets{}.GetFieldValue("ID", result.WalletID)
					if len(listRes) == 1 {
						searchList[pos].Wallet = listRes[0].Label
					}
				}
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

func apiAccountGet(httpRes http.ResponseWriter, httpReq *http.Request) {

	httpRes.Header().Set("Content-Type", "application/json")

	var statusBody interface{}
	statusCode := http.StatusOK
	statusMessage := ""

	if claims := utils.VerifyJWT(httpRes, httpReq); claims == nil {
		statusBody = map[string]string{"Redirect": "/"}
	} else {

		sAccountID := strings.TrimSpace(httpReq.FormValue("id"))
		if sAccountID == "" {
			statusCode = http.StatusInternalServerError
			statusMessage = "Error Account ID is required to load form"
		} else {
			AccountID, _ := strconv.ParseUint(sAccountID, 0, 64)
			accountsList, err := buckets.Accounts{}.GetFieldValue("ID", AccountID)
			if err != nil {
				statusMessage = err.Error()
			} else {
				if len(accountsList) > 0 {

					Wallet := ""
					if accountsList[0].WalletID != 0 {
						listRes, _ := buckets.Wallets{}.GetFieldValue("ID", accountsList[0].WalletID)
						if len(listRes) == 1 {
							Wallet = listRes[0].Label
						}
					}

					statusBody = apiAccountStruct{
						Wallet: Wallet,

						ID:       accountsList[0].ID,
						WalletID: accountsList[0].WalletID,

						Workflow:   accountsList[0].Workflow,
						Createdate: accountsList[0].Createdate,
						Updatedate: accountsList[0].Updatedate,

						Path:       accountsList[0].Path,
						Label:      accountsList[0].Label,
						Address:    accountsList[0].Address,
						Balance:    accountsList[0].Balance,
						Publickey:  accountsList[0].Publickey,
						Privatekey: accountsList[0].Privatekey,
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

func apiAccountPost(httpRes http.ResponseWriter, httpReq *http.Request) {
	httpRes.Header().Set("Content-Type", "application/json")

	var statusBody interface{}
	statusCode := http.StatusInternalServerError
	statusMessage := ""

	if claims := utils.VerifyJWT(httpRes, httpReq); claims == nil {
		statusBody = map[string]string{"Redirect": "/"}
	} else {

		formStruct := buckets.Accounts{}
		err := json.NewDecoder(httpReq.Body).Decode(&formStruct)
		if err != nil {
			statusMessage = "Error Decoding Form Values: " + err.Error()
		} else {

			bucketAccount := buckets.Accounts{}

			if formStruct.ID != 0 {
				bucketAccountList, _ := buckets.Accounts{}.GetFieldValue("ID", formStruct.ID)
				if len(bucketAccountList) != 1 {
					statusMessage = "Error Decoding Form Values: " + err.Error()
				} else {
					bucketAccount = bucketAccountList[0]
				}
			} else {
				formStruct.Workflow = Enabled
			}

			bucketAccount.Workflow = formStruct.Workflow
			bucketAccount.Label = formStruct.Label

			if statusMessage == "" {
				if bucketAccount.Label == "" {
					statusMessage += "Label is Required \n"
				}

				if bucketAccount.Workflow == "" {
					statusMessage += WorkflowRequired
				}

				if strings.HasSuffix(statusMessage, "\n") {
					statusMessage = statusMessage[:len(statusMessage)-2]
				}
			}

			if statusMessage == "" {
				err = bucketAccount.Create(&bucketAccount)
				if err != nil {
					statusMessage = "Error Saving Record: " + err.Error()
				} else {
					statusCode = http.StatusOK
					statusMessage = RecordSaved
					statusBody = bucketAccount.ID
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
