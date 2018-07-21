package buckets

import (
	"fmt"
	"log"
	"time"

	"golang.org/x/crypto/bcrypt"

	"github.com/coreos/bbolt"
	"github.com/timshannon/bolthold"

	"uniport/config"
)

//Users ...
type Users struct {
	ID, IDSync, Createdby, Updatedby,
	FailedMax, Failed uint64

	Createdate, Updatedate time.Time

	Password []byte

	IsStaff, IsAdmin bool

	Workflow, Code, Email, Mobile,
	Username, Title, Firstname,
	Lastname, Othername,

	Image, Phone,
	Description string

	PasswordString string `json:"Password"`
}

func (user Users) bucketName() string {
	return "Users"
}

//Setup ....
func (user Users) Setup() (err error) {

	hash, err := bcrypt.GenerateFromPassword([]byte("toor"), bcrypt.DefaultCost)
	if err != nil {
		log.Println(err)
		return
	}

	bucketType := Users{
		Email: "root@localhost",

		Workflow: "enabled",
		Fullname: "ROOT",
		Username: "root",
		Password: hash,
		IsAdmin:  true,
	}

	err = user.Create(&bucketType)
	if err != nil {
		log.Println(err)
	}
	return
}

//Create ...
func (user Users) Create(bucketType *Users) (err error) {

	if err = config.Get().BoltHold.Bolt().Update(func(tx *bolt.Tx) error {

		if bucketType.Createdate.IsZero() {
			bucketType.Createdate = time.Now()
			bucketType.Updatedate = bucketType.Createdate
		}

		if bucketType.ID == 0 {
			bucket := tx.Bucket([]byte(user.bucketName()))
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
func (user Users) List() (resultsALL []string) {
	var results []Users

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
func (user Users) GetFieldValue(Field string, Value interface{}) (results []Users, err error) {

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
