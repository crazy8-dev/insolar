/*
 *    Copyright 2018 Insolar
 *
 *    Licensed under the Apache License, Version 2.0 (the "License");
 *    you may not use this file except in compliance with the License.
 *    You may obtain a copy of the License at
 *
 *        http://www.apache.org/licenses/LICENSE-2.0
 *
 *    Unless required by applicable law or agreed to in writing, software
 *    distributed under the License is distributed on an "AS IS" BASIS,
 *    WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
 *    See the License for the specific language governing permissions and
 *    limitations under the License.
 */

package storage

import (
	"bytes"
	"context"
	"errors"

	"github.com/dgraph-io/badger"

	"github.com/insolar/insolar/core"
)

// iterstate stores iterator state
type iterstate struct {
	start  []byte
	prefix []byte
}

// ReplicaIter provides partial iterator over BadgerDB key/value pairs
// required for replication to Heavy Material node in provided pulse.
//
// Required kv-pairs are all records into and after provided pulse
// and all indexes available in database.
//
// "Partial" means it fetches data in chunks of the specified size.
// After a chunk has been fetched, it saves the current position.
// This is not an "honest" alogrithm, because the last record size can exceed the limit.
//
// Better implementation is for the future work.
type ReplicaIter struct {
	ctx        context.Context
	db         *DB
	limitBytes int
	istates    []*iterstate
}

// NewReplicaIter creates ReplicaIter started from provided pulse number with per iteration fetch limit.
func NewReplicaIter(ctx context.Context, db *DB, pulsenum core.PulseNumber, limit int) *ReplicaIter {
	recordsPrefix := []byte{scopeIDRecord}
	recordsStart := bytes.Join([][]byte{recordsPrefix, pulsenum.Bytes()}, nil)
	indexesPrefix := []byte{scopeIDLifeline}
	return &ReplicaIter{
		ctx:        ctx,
		db:         db,
		limitBytes: limit,

		istates: []*iterstate{
			&iterstate{recordsStart, recordsPrefix},
			&iterstate{indexesPrefix, indexesPrefix},
		},
	}
}

// NextRecords fetches next part of key value pairs.
func (r *ReplicaIter) NextRecords() ([]core.KV, error) {
	if r.isDone() {
		return nil, ErrReplicatorDone
	}
	fc := &fetchchunk{
		db:    r.db.db,
		limit: r.limitBytes,
	}
	for _, is := range r.istates {
		if is.start == nil {
			continue
		}
		var fetcherr error
		is.start, fetcherr = fc.fetch(r.ctx, is.prefix, is.start)
		if fetcherr != nil {
			return nil, fetcherr
		}
	}
	return fc.records, nil
}

// ErrReplicatorDone is returned by an Replicator NextRecords method when the iteration is complete.
var ErrReplicatorDone = errors.New("no more items in iterator")

func (r *ReplicaIter) isDone() bool {
	for _, is := range r.istates {
		if is.start != nil {
			return false
		}
	}
	return true
}

type fetchchunk struct {
	db      *badger.DB
	records []core.KV
	size    int
	limit   int
}

func (fc *fetchchunk) fetch(ctx context.Context, prefix []byte, start []byte) ([]byte, error) {
	if fc.size > fc.limit {
		return start, nil
	}

	var nextstart []byte
	err := fc.db.View(func(txn *badger.Txn) error {
		it := txn.NewIterator(badger.DefaultIteratorOptions)
		defer it.Close()

		for it.Seek(start); it.ValidForPrefix(prefix); it.Next() {
			item := it.Item()
			if item == nil {
				break
			}
			key := item.KeyCopy(nil)
			if fc.size > fc.limit {
				nextstart = key
				// inslogger.FromContext(ctx).Warnf("size > r.limit: %v > %v (nextstart=%v)",
				// 	fc.size, fc.limit, hex.EncodeToString(key))
				return nil
			}

			value, err := it.Item().ValueCopy(nil)
			if err != nil {
				return err
			}
			fc.records = append(fc.records, core.KV{K: key, V: value})
			fc.size += len(key) + len(value)
		}
		nextstart = nil
		return nil
	})
	return nextstart, err
}
