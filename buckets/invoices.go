package buckets

import (
	"fmt"
	"log"
	"time"

	"github.com/coreos/bbolt"
	"github.com/timshannon/bolthold"

	"uniport/config"
)

//ProformaInvoice
//AgentInvoice
//CommercialInvoice

//Invoices ...
type Invoices struct {
	ID uint64
	Title, Description,
	Workflow string

	Createdate, Updatedate,
	ShipmentDate time.Time

	AgentID, SupplierID uint64

	Type, Agent, VatNo, Tin, OurRef,
	YourRef, Date, Attention,
	Supplier, Country, Products,
	Quantity, DutyRate,
	DisbursementRate,
	MaintanceFeeRate,
	Currency, VatRate,
	DeliveryDate, DateDue,
	Paymentterms, Shipmentterms,
	BankOne, BankTwo, AccountOne,
	AccountTwo, WayBillNo,
	Origin string

	ExWorks, FOBcharges,
	Freight, CandF, SubTotal,
	Disbursement, MaintanceFee,
	TotalExclVat, Vat,
	TotalInclVat, AmountDue float64
}

func (invoice Invoices) bucketName() string {
	return "Invoices"
}

//Create ...
func (invoice Invoices) Create(bucketType *Invoices) (err error) {

	if err = config.Get().BoltHold.Bolt().Update(func(tx *bolt.Tx) error {

		if bucketType.Createdate.IsZero() {
			bucketType.Createdate = time.Now()
			bucketType.Updatedate = bucketType.Createdate
		}

		if bucketType.ID == 0 {
			bucket := tx.Bucket([]byte(invoice.bucketName()))
			bucketType.ID, _ = bucket.NextSequence()
			bucketType.Createdate = time.Now()
		} else {
			bucketType.Updatedate = time.Now()
		}

		err = config.Get().BoltHold.TxUpsert(tx, bucketType.ID, bucketType)
		return err
	}); err != nil {
		log.Printf(err.Error())
	}
	return
}

//List ...
func (invoice Invoices) List() (resultsALL []string) {
	var results []Invoices

	if err := config.Get().BoltHold.Bolt().View(func(tx *bolt.Tx) error {
		err := config.Get().BoltHold.Find(&results, bolthold.Where("ID").Gt(uint64(0)))
		return err
	}); err != nil {
		log.Printf(err.Error())
	} else {
		for _, record := range results {
			resultsALL = append(resultsALL, fmt.Sprintf("%+v", record))
		}
	}
	return
}

//GetFieldValue ...
func (invoice Invoices) GetFieldValue(Field string, Value interface{}) (results []Invoices, err error) {

	if len(Field) > 0 {
		if err = config.Get().BoltHold.Bolt().View(func(tx *bolt.Tx) error {
			err = config.Get().BoltHold.Find(&results, bolthold.Where(Field).Eq(Value).SortBy("ID").Reverse())
			return err
		}); err != nil {
			log.Printf(err.Error())
		}
	}
	return
}
