package buckets

import (
	"fmt"
	"log"
	"time"

	"github.com/coreos/bbolt"
	"github.com/timshannon/bolthold"

	"uniport/config"
)

//Suppliers ...
type Suppliers struct {
	ID, Createdby,
	Updatedby uint64

	Createdate, Updatedate time.Time

	Workflow, Code, Email, Mobile,
	Company, Title, Firstname,
	Lastname, Website,

	Image, Phone,
	Description string
}

func (supplier Suppliers) bucketName() string {
	return "Suppliers"
}

//Create ...
func (supplier Suppliers) Create(bucketType *Suppliers) (err error) {

	if err = config.Get().BoltHold.Bolt().Update(func(tx *bolt.Tx) error {

		if bucketType.Createdate.IsZero() {
			bucketType.Createdate = time.Now()
			bucketType.Updatedate = bucketType.Createdate
		}

		if bucketType.ID == 0 {
			bucket := tx.Bucket([]byte(supplier.bucketName()))
			bucketType.ID, _ = bucket.NextSequence()
			bucketType.Createdate = time.Now()
			bucketType.Createdby = bucketType.Updatedby
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
func (supplier Suppliers) List() (resultsALL []string) {
	var results []Suppliers

	if err := config.Get().BoltHold.Bolt().View(func(tx *bolt.Tx) error {
		err := config.Get().BoltHold.Find(&results, bolthold.Where("ID").Gt(uint64(0)))
		return err
	}); err != nil {
		log.Printf(err.Error())
	} else {
		for _, row := range results {
			resultsALL = append(resultsALL, fmt.Sprintf("%+v", row))
		}
	}
	return
}

//GetFieldValue ...
func (supplier Suppliers) GetFieldValue(Field string, Value interface{}) (results []Suppliers, err error) {

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
