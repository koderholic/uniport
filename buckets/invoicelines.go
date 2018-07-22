package buckets

import (
	"fmt"
	"log"
	"time"

	"github.com/coreos/bbolt"
	"github.com/timshannon/bolthold"

	"uniport/config"
)

type Invoicelines struct {
	ID uint64
	Title, Description,
	Workflow string
	Createdate,
	Updatedate time.Time

	Lineno, Pallets uint64
	Quantity, Price,
	Total float64

	SalesOrder, DeliveryNote,
	PurchaseOrder string
}

func (this Invoicelines) bucketName() string {
	return "Invoicelines"
}

func (this Invoicelines) Create(bucketType *Invoicelines) (err error) {

	if err = config.Get().BoltHold.Bolt().Update(func(tx *bolt.Tx) error {

		if bucketType.Createdate.IsZero() {
			bucketType.Createdate = time.Now()
			bucketType.Updatedate = bucketType.Createdate
		}

		if bucketType.ID == 0 {
			bucket := tx.Bucket([]byte(this.bucketName()))
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

func (this Invoicelines) List() (resultsALL []string) {
	var results []Invoicelines

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

func (this Invoicelines) GetFieldValue(Field string, Value interface{}) (results []Invoicelines, err error) {

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
