package main

import (
	"os"
	"time"

	"github.com/boltdb/bolt"
)

type Store struct {
	BucketName string
	Bucket     *bolt.Bucket
	Path       string
	Perms      os.FileMode
	db         *bolt.DB
	Timeout    time.Duration
}

func NewStore(path string, perms os.FileMode, time time.Duration) (*Store, error) {

	self := &Store{
		BucketName: "BroTop",
		Path:       path,
		Perms:      perms,
		Timeout:    time,
	}

	err := self.Open()

	return self, err
}

func (self *Store) Open() error {
	db, err := bolt.Open(self.Path, self.Perms, &bolt.Options{
		Timeout: self.Timeout,
	})

	if err != nil {
		return err
	}

	self.db = db

	self.db.Update(func(tx *bolt.Tx) error {
		self.Bucket, err = tx.CreateBucketIfNotExists([]byte(self.BucketName))

		if err != nil {
			return err
		}

		return nil
	})

	return nil
}

func (self *Store) Close() {
	self.db.Close()
}

func (self *Store) Get(key string) ([]byte, error) {

	var value []byte

	err := self.db.View(func(tx *bolt.Tx) error {
		value = self.Bucket.Get([]byte(key))
		return nil
	})

	return value, err
}

func (self *Store) Set(key, value string) error {

	err := self.db.Update(func(tx *bolt.Tx) error {
		err := self.Bucket.Put([]byte(key), []byte(value))
		return err
	})

	return err
}

func (self *Store) Update(key, value string) {

}

func (self *Store) Delete(key string) {

}
