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

type apiRecipientsStruct struct {
	ID       uint64
	Workflow string
	Createdate,
	Updatedate time.Time

	Label, Email, Mobile,
	Address string
}

func apiHandlerRecipients(middlewares alice.Chain, router *Router) {
	router.Get("/api/recipients", middlewares.ThenFunc(apiRecipientGet))
	router.Post("/api/recipients", middlewares.ThenFunc(apiRecipientPost))
	router.Get("/api/recipients/search", middlewares.ThenFunc(apiRecipientsSearch))
}

func apiRecipientsSearch(httpRes http.ResponseWriter, httpReq *http.Request) {

	httpRes.Header().Set("Content-Type", "application/json")

	var statusBody interface{}
	statusCode := http.StatusInternalServerError
	statusMessage := ""

	if claims := utils.VerifyJWT(httpRes, httpReq); claims == nil {
		statusBody = map[string]string{"Redirect": "/"}
	} else {

		searchResults := []buckets.Recipients{}
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

func apiRecipientGet(httpRes http.ResponseWriter, httpReq *http.Request) {

	httpRes.Header().Set("Content-Type", "application/json")

	var statusBody interface{}
	statusCode := http.StatusOK
	statusMessage := ""

	if claims := utils.VerifyJWT(httpRes, httpReq); claims == nil {
		statusBody = map[string]string{"Redirect": "/"}
	} else {

		sRecipientID := strings.TrimSpace(httpReq.FormValue("id"))
		if sRecipientID == "" {
			statusCode = http.StatusInternalServerError
			statusMessage = "Error Recipient ID is required to load form"
		} else {
			RecipientID, _ := strconv.ParseUint(sRecipientID, 0, 64)
			recipientsList, err := buckets.Recipients{}.GetFieldValue("ID", RecipientID)
			if err != nil {
				statusMessage = err.Error()
			} else {
				if len(recipientsList) > 0 {
					statusBody = apiRecipientsStruct{
						ID:         recipientsList[0].ID,
						Createdate: recipientsList[0].Createdate,
						Updatedate: recipientsList[0].Updatedate,

						Label:   recipientsList[0].Label,
						Email:   recipientsList[0].Email,
						Mobile:  recipientsList[0].Mobile,
						Address: recipientsList[0].Address,
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

func apiRecipientPost(httpRes http.ResponseWriter, httpReq *http.Request) {
	httpRes.Header().Set("Content-Type", "application/json")

	var statusBody interface{}
	statusCode := http.StatusInternalServerError
	statusMessage := ""

	if claims := utils.VerifyJWT(httpRes, httpReq); claims == nil {
		statusBody = map[string]string{"Redirect": "/"}
	} else {

		formStruct := buckets.Recipients{}
		err := json.NewDecoder(httpReq.Body).Decode(&formStruct)
		if err != nil {
			statusMessage = "Error Decoding Form Values: " + err.Error()
		} else {

			bucketRecipient := buckets.Recipients{}

			if formStruct.ID != 0 {
				bucketRecipientList, _ := buckets.Recipients{}.GetFieldValue("ID", formStruct.ID)
				if len(bucketRecipientList) != 1 {
					statusMessage = "Error Decoding Form Values: " + err.Error()
				} else {
					bucketRecipient = bucketRecipientList[0]
				}
			} else {
				formStruct.Workflow = Enabled
			}

			bucketRecipient.Label = formStruct.Label
			bucketRecipient.Email = formStruct.Email
			bucketRecipient.Mobile = formStruct.Mobile
			bucketRecipient.Address = formStruct.Address

			if statusMessage == "" {
				if bucketRecipient.Label == "" {
					statusMessage += "Label is Required \n"
				}

				if bucketRecipient.Address == "" {
					statusMessage += "Address is Required \n"
				}

				if strings.HasSuffix(statusMessage, "\n") {
					statusMessage = statusMessage[:len(statusMessage)-2]
				}
			}

			if statusMessage == "" {
				err = bucketRecipient.Create(&bucketRecipient)
				if err != nil {
					statusMessage = "Error Saving Record: " + err.Error()
				} else {
					statusCode = http.StatusOK
					statusMessage = RecordSaved
					statusBody = bucketRecipient.ID
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
