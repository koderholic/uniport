package buckets

import (
	"fmt"
	"log"
	"time"

	"github.com/coreos/bbolt"
	"github.com/timshannon/bolthold"

	"uniport/config"
)

//Type, Status, DocID,

//Formms ...
type Formms struct {
	ID uint64
	Title, Workflow,
	Description string
	Createdate,
	Updatedate time.Time

	PackageQty, UomQty   float64
	PackageType, UomType string

	DocNo, Hscode,
	Currency, PFIno, PFIdate,
	Transport, CountryOrigin,
	CountrySupply string

	InvoiceID uint64
}

func (Formm Formms) bucketName() string {
	return "Formms"
}

//Create ...
func (Formm Formms) Create(bucketType *Formms) (err error) {

	if err = config.Get().BoltHold.Bolt().Update(func(tx *bolt.Tx) error {

		if bucketType.Createdate.IsZero() {
			bucketType.Createdate = time.Now()
			bucketType.Updatedate = bucketType.Createdate
		}

		if bucketType.ID == 0 {
			bucket := tx.Bucket([]byte(Formm.bucketName()))
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
func (Formm Formms) List() (resultsALL []string) {
	var results []Formms

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
func (Formm Formms) GetFieldValue(Field string, Value interface{}) (results []Formms, err error) {

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
