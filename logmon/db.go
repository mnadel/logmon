package logmon

import (
	"bytes"
	"encoding/binary"
	"log"

	"github.com/boltdb/bolt"
)

type Database struct {
	boltdb *bolt.DB
	config *Configuration
}

func NewDatabase(config *Configuration) *Database {
	db, err := bolt.Open(config.Db, 0600, nil)
	if err != nil {
		log.Fatal("error opening db:", config.Db, err.Error())
	}

	err = db.Update(func(tx *bolt.Tx) error {
		_, err := tx.CreateBucketIfNotExists([]byte("hashes"))
		if err != nil {
			return err
		}

		_, err = tx.CreateBucketIfNotExists([]byte("offsets"))

		return err
	})

	if err != nil {
		log.Fatal("error initializing database:", err.Error())
	}

	return &Database{
		boltdb: db,
		config: config,
	}
}

func (db *Database) Close() {
	db.boltdb.Close()
}

func (db *Database) updateFile(path string, offset uint64, hash []byte) error {
	return db.boltdb.Update(func(tx *bolt.Tx) error {
		hashBucket := tx.Bucket([]byte("hashes"))
		hashBucket.Put([]byte(path), hash)

		offsetBucket := tx.Bucket([]byte("offsets"))

		buf := make([]byte, 10)
		binary.PutUvarint(buf, offset)

		offsetBucket.Put([]byte(path), buf)

		return nil
	})
}

func (db *Database) getHash(path string) ([]byte, error) {
	var hash []byte

	err := db.boltdb.View(func(tx *bolt.Tx) error {
		bucket := tx.Bucket([]byte("hashes"))

		hash = bucket.Get([]byte(path))

		return nil
	})

	return hash, err
}

func (db *Database) getOffset(path string) (uint64, error) {
	var offset uint64

	err := db.boltdb.View(func(tx *bolt.Tx) error {
		bucket := tx.Bucket([]byte("offsets"))

		byteSlice := bucket.Get([]byte(path))
		if byteSlice != nil {
			v, err := binary.ReadUvarint(bytes.NewReader(byteSlice))
			if err != nil {
				return err
			}

			offset = v
		}

		return nil
	})

	return offset, err
}
