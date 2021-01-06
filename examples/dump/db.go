package main

import (
	"encoding/json"
	"fmt"
	"math/big"

	"github.com/boltdb/bolt"
	"github.com/patractlabs/go-patract/contracts/erc20"
	"github.com/patractlabs/go-patract/utils"
	"github.com/pkg/errors"
)

type AccountAsset struct {
	Sequence uint64  `json:"sequence"`
	Address  string  `json:"address"`
	Amount   big.Int `json:"amount"`
}

type erc20DB struct {
	*bolt.DB

	BucketKey []byte
}

func NewErc20DB(path string) *erc20DB {
	db, err := bolt.Open(path, 0600, nil)
	if err != nil {
		panic(err)
	}

	return &erc20DB{
		DB:        db,
		BucketKey: []byte{},
	}
}

func (db *erc20DB) Close() {
	db.Close()
}

func (db *erc20DB) init(hash string) error {
	// Start a writable transaction.
	tx, err := db.Begin(true)
	if err != nil {
		return err
	}
	defer tx.Rollback()

	db.BucketKey = []byte(fmt.Sprintf("erc20-%s", hash))

	// Use the transaction...
	_, err = tx.CreateBucketIfNotExists(db.BucketKey)
	if err != nil {
		return err
	}

	// Commit the transaction and check for error.
	if err := tx.Commit(); err != nil {
		return err
	}

	return nil
}

func (db *erc20DB) updateUserAsset(b *bolt.Bucket, address string, amtDeta *big.Int) error {
	a := b.Get([]byte(address))
	if a == nil {
		if amtDeta.Cmp(big.NewInt(0)) < 0 {
			return errors.New("asset can not be negative when created")
		}

		seq, _ := b.NextSequence()

		na := AccountAsset{
			Address:  address,
			Amount:   *amtDeta,
			Sequence: seq,
		}

		buf, err := json.Marshal(na)
		if err != nil {
			return err
		}

		return b.Put([]byte(address), buf)
	}

	na := AccountAsset{}

	err := json.Unmarshal(a, &na)
	if err != nil {
		return errors.Wrap(err, "unmarshal")
	}

	na.Amount = *na.Amount.Add(&na.Amount, amtDeta)

	buf, err := json.Marshal(na)
	if err != nil {
		return err
	}

	return b.Put([]byte(address), buf)
}

func (db *erc20DB) OnEventTransfer(evt *erc20.EventTransfer) error {
	return db.Update(func(tx *bolt.Tx) error {
		b := tx.Bucket(db.BucketKey)

		if evt.From.IsSome() {
			from, err := utils.EncodeAccountIDToSS58(evt.From.Value)
			if err != nil {
				return errors.Wrapf(err, "encode err")
			}

			amt := evt.Value.Int
			if amt.Cmp(big.NewInt(0)) != 0 {
				return db.updateUserAsset(b, from, amt.Neg(amt))
			}
		}

		if evt.To.IsSome() {
			to, err := utils.EncodeAccountIDToSS58(evt.To.Value)
			if err != nil {
				return errors.Wrapf(err, "encode err")
			}

			amt := evt.Value.Int
			if amt.Cmp(big.NewInt(0)) != 0 {
				return db.updateUserAsset(b, to, amt)
			}
		}

		return nil
	})
}
