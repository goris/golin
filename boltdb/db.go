package boltdb

/*
/   This is just a small adapter to do some particular wrapping of the
/   boltdb/bolt library so it fits with our needs
/
*/

import (
	"fmt"
	"github.com/boltdb/bolt"
)

func OpenBoltDB(name string) (*bolt.DB, error) {
	name = fmt.Sprintf("%s.db", name)
	boltDB, err := bolt.Open(name, 0644, nil)
	return boltDB, err
}

func UpdateBucket(db *bolt.DB, bucket string, data map[string]interface{}) error {
	fmt.Println(bucket, data)
	err := db.Update(func(tx *bolt.Tx) error {
		bucket, err := tx.CreateBucketIfNotExists([]byte(bucket))
		if err != nil {
			return fmt.Errorf("create bucket: %s", err)
		}
		for k, v := range data {
			if str, ok := v.(string); ok {
				err = bucket.Put([]byte(k), []byte(str))
			}
		}
		return err
	})
	return err
}
