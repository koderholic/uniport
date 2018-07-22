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

//Checklistdetails ...
type Checklistdetails struct {
	ID uint64
	Title, Workflow,
	Description string
	Createdate,
	Updatedate time.Time

	ChecklistID,
	DocID, RefID uint64
	Ref, Type string
}

func (Checklist Checklistdetails) bucketName() string {
	return "Checklistdetails"
}

//Create ...
func (Checklist Checklistdetails) Create(bucketType *Checklistdetails) (err error) {

	if err = config.Get().BoltHold.Bolt().Update(func(tx *bolt.Tx) error {

		if bucketType.Createdate.IsZero() {
			bucketType.Createdate = time.Now()
			bucketType.Updatedate = bucketType.Createdate
		}

		if bucketType.ID == 0 {
			bucket := tx.Bucket([]byte(Checklist.bucketName()))
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
func (Checklist Checklistdetails) List() (resultsALL []string) {
	var results []Checklistdetails

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
func (Checklist Checklistdetails) GetFieldValue(Field string, Value interface{}) (results []Checklistdetails, err error) {

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
