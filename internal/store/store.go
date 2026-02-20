package store

import (
	"fmt"

	badger "github.com/dgraph-io/badger/v4"
)

type Store struct {
	db *badger.DB
}

func NewStore(path string) (*Store, error) {
	opts := badger.DefaultOptions(path)
	db, err := badger.Open(opts)

	if err != nil {
		return nil, fmt.Errorf("could not create db: %w", err)
	}

	return &Store{db: db}, nil
}

func (s *Store) Exists(key string) (bool, error) {
	var exists bool

	err := s.db.View(func(txn *badger.Txn) error {
		if record, err := txn.Get([]byte(key)); err != nil {
			return err
		} else if record != nil {
			exists = true
		}

		return nil
	})

	return exists, err
}

func (s *Store) Get(key string) (string, error) {
	var value string

	err := s.db.View(func(txn *badger.Txn) error {
		record, err := txn.Get([]byte(key))

		if err != nil {
			return err
		}

		valCopy, err := record.ValueCopy(nil)

		if err != nil {
			return err
		}

		value = string(valCopy)

		return nil
	})

	return value, err
}

func (s *Store) Set(key string, value string) error {
	return s.db.Update(func(txn *badger.Txn) error {
		return txn.Set([]byte(key), []byte(value))
	})
}

func (s *Store) Close() error {
	return s.db.Close()
}
