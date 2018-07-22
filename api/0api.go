package api

import (
	"fmt"
	"time"

	"github.com/justinas/alice"
)

const (
	IsRequired = " is Required \n" //isRequired is required

	RecordSaved      = "Record Saved" //recordSaved is saved
	Enabled          = "enabled"
	TitleRequired    = "Title is Required \n"
	WorkflowRequired = "Workflow is Required \n"
)

func apiHandler(middlewares alice.Chain, router *Router) {

	//Website API Handler
	apiHandlerWebsite(middlewares, router)

	//Other Functionalities
	apiHandlerAccounts(middlewares, router)
	apiHandlerRecipients(middlewares, router)
	apiHandlerTransactions(middlewares, router)
	apiHandlerWallets(middlewares, router)

}

type JSONTime time.Time

func (t JSONTime) MarshalJSON() ([]byte, error) {
	//do your serializing here
	stamp := fmt.Sprintf("\"%s\"", time.Time(t).Format("02/01/2006 03:04:05 PM"))
	return []byte(stamp), nil
}

type apiSearch struct {
	Field, Text, Code,
	Workflow string

	ID, Skip, Limit,
	DocumentID, CID int
}

type apiSearchResult struct {
	ID uint64
	Code, Details, Username,
	SubDetails, Workflow string
	Date JSONTime
	Row  interface{}
}

/* This function reverses the order of a slice*/
func apiReverseSlice(result []apiSearchResult) {
	for i := len(result)/2 - 1; i >= 0; i-- {
		opp := len(result) - 1 - i
		result[i], result[opp] = result[opp], result[i]
	}
}
