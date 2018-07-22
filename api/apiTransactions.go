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

type apiTransactionStruct struct {
	ID       uint64
	Workflow string
	Createdate,
	Updatedate time.Time

	Block uint32

	Reference, FromAddress,
	ToAddress string

	Amount                   *big.Int
	WalletID, AccountID      uint64
	Wallet, Account, Address string
}

func apiHandlerTransactions(middlewares alice.Chain, router *Router) {
	router.Get("/api/transactions", middlewares.ThenFunc(apiTransactionGet))
	router.Post("/api/transactions", middlewares.ThenFunc(apiTransactionPost))
	router.Post("/api/transactions/search", middlewares.ThenFunc(apiTransactionSearch))
}

func apiTransactionSearch(httpRes http.ResponseWriter, httpReq *http.Request) {

	httpRes.Header().Set("Content-Type", "application/json")

	var statusBody interface{}
	statusCode := http.StatusInternalServerError
	statusMessage := ""

	if claims := utils.VerifyJWT(httpRes, httpReq); claims == nil {
		statusBody = map[string]string{"Redirect": "/"}
	} else {

		searchResults := []buckets.Transactions{}
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
			formSearch.Field = "Reference"
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

			searchList := make([]apiTransactionStruct, len(searchResults))
			for pos, result := range searchResults {

				searchList[pos].ID = result.ID
				searchList[pos].Workflow = result.Workflow
				searchList[pos].Createdate = result.Createdate
				searchList[pos].Updatedate = result.Updatedate

				searchList[pos].Block = result.Block
				searchList[pos].Reference = result.Reference
				searchList[pos].ToAddress = result.ToAddress
				searchList[pos].FromAddress = result.FromAddress
				searchList[pos].Amount = result.Amount

				searchList[pos].WalletID = result.WalletID
				if result.WalletID != 0 {
					listRes, _ := buckets.Wallets{}.GetFieldValue("ID", result.WalletID)
					if len(listRes) == 1 {
						searchList[pos].Wallet = listRes[0].Label
					}
				}

				searchList[pos].AccountID = result.AccountID
				if result.AccountID != 0 {
					listRes, _ := buckets.Accounts{}.GetFieldValue("ID", result.AccountID)
					if len(listRes) == 1 {
						searchList[pos].Account = listRes[0].Label
						searchList[pos].Address = listRes[0].Address
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

func apiTransactionGet(httpRes http.ResponseWriter, httpReq *http.Request) {

	httpRes.Header().Set("Content-Type", "application/json")

	var statusBody interface{}
	statusCode := http.StatusOK
	statusMessage := ""

	if claims := utils.VerifyJWT(httpRes, httpReq); claims == nil {
		statusBody = map[string]string{"Redirect": "/"}
	} else {

		sTransactionID := strings.TrimSpace(httpReq.FormValue("id"))
		if sTransactionID == "" {
			statusCode = http.StatusInternalServerError
			statusMessage = "Error Transaction ID is required to load form"
		} else {
			TransactionID, _ := strconv.ParseUint(sTransactionID, 0, 64)
			transactionsList, err := buckets.Transactions{}.GetFieldValue("ID", TransactionID)
			if err != nil {
				statusMessage = err.Error()
			} else {
				if len(transactionsList) > 0 {

					Wallet := ""
					if transactionsList[0].WalletID != 0 {
						listRes, _ := buckets.Wallets{}.GetFieldValue("ID", transactionsList[0].WalletID)
						if len(listRes) == 1 {
							Wallet = listRes[0].Label
						}
					}

					Account := ""
					Address := ""
					if transactionsList[0].AccountID != 0 {
						listRes, _ := buckets.Accounts{}.GetFieldValue("ID", transactionsList[0].AccountID)
						if len(listRes) == 1 {
							Account = listRes[0].Label
							Address = listRes[0].Address
						}
					}

					statusBody = apiTransactionStruct{
						Wallet:  Wallet,
						Account: Account,
						Address: Address,

						ID:        transactionsList[0].ID,
						WalletID:  transactionsList[0].WalletID,
						AccountID: transactionsList[0].AccountID,

						Workflow:   transactionsList[0].Workflow,
						Createdate: transactionsList[0].Createdate,
						Updatedate: transactionsList[0].Updatedate,

						Block:       transactionsList[0].Block,
						Reference:   transactionsList[0].Reference,
						ToAddress:   transactionsList[0].ToAddress,
						FromAddress: transactionsList[0].FromAddress,
						Amount:      transactionsList[0].Amount,
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

func apiTransactionPost(httpRes http.ResponseWriter, httpReq *http.Request) {
	httpRes.Header().Set("Content-Type", "application/json")

	var statusBody interface{}
	statusCode := http.StatusInternalServerError
	statusMessage := ""

	if claims := utils.VerifyJWT(httpRes, httpReq); claims == nil {
		statusBody = map[string]string{"Redirect": "/"}
	} else {

		formStruct := buckets.Transactions{}
		err := json.NewDecoder(httpReq.Body).Decode(&formStruct)
		if err != nil {
			statusMessage = "Error Decoding Form Values: " + err.Error()
		} else {

			bucketTransaction := buckets.Transactions{}

			if formStruct.ID != 0 {
				bucketTransactionList, _ := buckets.Transactions{}.GetFieldValue("ID", formStruct.ID)
				if len(bucketTransactionList) != 1 {
					statusMessage = "Error Decoding Form Values: " + err.Error()
				} else {
					bucketTransaction = bucketTransactionList[0]
				}
			} else {
				formStruct.Workflow = Enabled
			}

			bucketTransaction.WalletID = formStruct.WalletID
			bucketTransaction.AccountID = formStruct.AccountID

			bucketTransaction.Workflow = formStruct.Workflow

			bucketTransaction.Block = formStruct.Block
			bucketTransaction.Reference = formStruct.Reference
			bucketTransaction.ToAddress = formStruct.ToAddress
			bucketTransaction.FromAddress = formStruct.FromAddress
			bucketTransaction.Amount = formStruct.Amount

			if statusMessage == "" {

				// if bucketTransaction.Reference == "" {
				// 	statusMessage += "Reference is Required \n"
				// }
				if bucketTransaction.ToAddress == "" {
					statusMessage += "ToAddress is Required \n"
				}
				if bucketTransaction.FromAddress == "" {
					statusMessage += "FromAddress is Required \n"
				}

				if bucketTransaction.Workflow == "" {
					statusMessage += WorkflowRequired
				}

				if bucketTransaction.Amount.Cmp(big.NewInt(0)) <= 0 {
					statusMessage += "Pls provide Amount \n"
				}

				if strings.HasSuffix(statusMessage, "\n") {
					statusMessage = statusMessage[:len(statusMessage)-2]
				}
			}

			if statusMessage == "" {
				err = bucketTransaction.Create(&bucketTransaction)
				if err != nil {
					statusMessage = "Error Saving Record: " + err.Error()
				} else {
					statusCode = http.StatusOK
					statusMessage = RecordSaved
					statusBody = bucketTransaction.ID
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
