package buckets

import (
	"strings"

	"github.com/coreos/bbolt"

	"uniport/config"
)

//Init ...
func Init() {
	Empty("/all")
}

var allowedBuckets = map[string]bool{
	"Accounts":     true,
	"Recipients":   true,
	"Transactions": true,
	"Wallets":      true,
}

//BucketList ...
func BucketList(bucketName string) (bucketList []string) {
	if len(bucketName) > 0 {
		bucketName = bucketName[1:]
	}
	switch bucketName {
	default:
		bucketList = append(bucketList, "Please Specify Bucket --> Bucket "+bucketName+" Invalid!!")

	case "Accounts":
		bucketList = append(bucketList, strings.Join(new(Accounts).List(), "\n"))

	case "Recipients":
		bucketList = append(bucketList, strings.Join(new(Recipients).List(), "\n"))

	case "Transactions":
		bucketList = append(bucketList, strings.Join(new(Transactions).List(), "\n"))

	case "Wallets":
		bucketList = append(bucketList, strings.Join(new(Wallets).List(), "\n"))
	}
	return
}

//Empty ...
func Empty(bucketName string) (Message []string) {

	switch bucketName {
	default:
		bucketName = bucketName[1:]
		if allowedBuckets[bucketName] {
			Message = append(Message, empty(bucketName))
		} else {
			Message = append(Message, "Please Specify Bucket")
		}

	case "/all":
		for bucket := range allowedBuckets {
			bucket = strings.Title(strings.ToLower(bucket))
			Message = append(Message, empty(bucket))
		}
		//Setup Users
	}
	return Message
}

func empty(bucketName string) string {
	if allowedBuckets[bucketName] {
		if err := config.Get().BoltHold.Bolt().Update(func(tx *bolt.Tx) (err error) {
			tx.DeleteBucket([]byte(bucketName))
			_, err = tx.CreateBucket([]byte(bucketName))
			return
		}); err != nil {
			return bucketName + " Bucket --> " + err.Error()
		}
		return bucketName + " Bucket -->  Emptied "
	}
	return bucketName + " Bucket -->  Does not Exist"
}
