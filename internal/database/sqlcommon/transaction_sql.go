// Copyright © 2022 Kaleido, Inc.
//
// SPDX-License-Identifier: Apache-2.0
//
// Licensed under the Apache License, Version 2.0 (the "License");
// you may not use this file except in compliance with the License.
// You may obtain a copy of the License at
//
//     http://www.apache.org/licenses/LICENSE-2.0
//
// Unless required by applicable law or agreed to in writing, software
// distributed under the License is distributed on an "AS IS" BASIS,
// WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
// See the License for the specific language governing permissions and
// limitations under the License.

package sqlcommon

import (
	"context"
	"database/sql"

	sq "github.com/Masterminds/squirrel"
	"github.com/hyperledger/firefly-common/pkg/ffapi"
	"github.com/hyperledger/firefly-common/pkg/fftypes"
	"github.com/hyperledger/firefly-common/pkg/i18n"
	"github.com/hyperledger/firefly-common/pkg/log"
	"github.com/hyperledger/firefly/internal/coremsgs"
	"github.com/hyperledger/firefly/pkg/core"
	"github.com/hyperledger/firefly/pkg/database"
)

var (
	transactionColumns = []string{
		"id",
		"ttype",
		"namespace",
		"created",
		"idempotency_key",
		"blockchain_ids",
	}
	transactionFilterFieldMap = map[string]string{
		"type":           "ttype",
		"idempotencykey": "idempotency_key",
		"blockchainids":  "blockchain_ids",
	}
)

const transactionsTable = "transactions"

func (s *SQLCommon) InsertTransaction(ctx context.Context, transaction *core.Transaction) (err error) {
	ctx, tx, autoCommit, err := s.BeginOrUseTx(ctx)
	if err != nil {
		return err
	}
	defer s.RollbackTx(ctx, tx, autoCommit)

	transaction.Created = fftypes.Now()
	var seq int64
	if seq, err = s.InsertTxExt(ctx, transactionsTable, tx,
		sq.Insert(transactionsTable).
			Columns(transactionColumns...).
			Values(
				transaction.ID,
				string(transaction.Type),
				transaction.Namespace,
				transaction.Created,
				transaction.IdempotencyKey,
				transaction.BlockchainIDs,
			),
		func() {
			s.callbacks.UUIDCollectionNSEvent(database.CollectionTransactions, core.ChangeEventTypeCreated, transaction.Namespace, transaction.ID)
		},
		transaction.IdempotencyKey != "", // on conflict we want to check for idempotency key mismatch to return a useful error
	); err != nil || seq < 0 {
		// Check for a duplicate idempotency key
		if transaction.IdempotencyKey != "" {
			fb := database.TransactionQueryFactory.NewFilter(ctx)
			existing, _, _ := s.GetTransactions(ctx, transaction.Namespace, fb.Eq("idempotencykey", (string)(transaction.IdempotencyKey)))
			if len(existing) > 0 {
				return i18n.NewError(ctx, coremsgs.MsgIdempotencyKeyDuplicateTransaction, transaction.IdempotencyKey, existing[0].ID)
			} else if err == nil {
				// If we don't have an error, and we didn't find an existing idempotency key match, then we don't know the reason for the conflict
				return i18n.NewError(ctx, coremsgs.MsgNonIdempotencyKeyConflictTxInsert, transaction.ID, transaction.IdempotencyKey)
			}
		}
		return err
	}

	return s.CommitTx(ctx, tx, autoCommit)
}

func (s *SQLCommon) transactionResult(ctx context.Context, row *sql.Rows) (*core.Transaction, error) {
	var transaction core.Transaction
	err := row.Scan(
		&transaction.ID,
		&transaction.Type,
		&transaction.Namespace,
		&transaction.Created,
		&transaction.IdempotencyKey,
		&transaction.BlockchainIDs,
	)
	if err != nil {
		return nil, i18n.WrapError(ctx, err, coremsgs.MsgDBReadErr, transactionsTable)
	}
	return &transaction, nil
}

func (s *SQLCommon) GetTransactionByID(ctx context.Context, namespace string, id *fftypes.UUID) (message *core.Transaction, err error) {

	rows, _, err := s.Query(ctx, transactionsTable,
		sq.Select(transactionColumns...).
			From(transactionsTable).
			Where(sq.Eq{"id": id, "namespace": namespace}),
	)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	if !rows.Next() {
		log.L(ctx).Debugf("Transaction '%s' not found", id)
		return nil, nil
	}

	transaction, err := s.transactionResult(ctx, rows)
	if err != nil {
		return nil, err
	}

	return transaction, nil
}

func (s *SQLCommon) GetTransactions(ctx context.Context, namespace string, filter ffapi.Filter) (message []*core.Transaction, fr *ffapi.FilterResult, err error) {

	query, fop, fi, err := s.FilterSelect(ctx, "", sq.Select(transactionColumns...).From(transactionsTable), filter, transactionFilterFieldMap, []interface{}{"sequence"}, sq.Eq{"namespace": namespace})
	if err != nil {
		return nil, nil, err
	}

	rows, tx, err := s.Query(ctx, transactionsTable, query)
	if err != nil {
		return nil, nil, err
	}
	defer rows.Close()

	transactions := []*core.Transaction{}
	for rows.Next() {
		transaction, err := s.transactionResult(ctx, rows)
		if err != nil {
			return nil, nil, err
		}
		transactions = append(transactions, transaction)
	}

	return transactions, s.QueryRes(ctx, transactionsTable, tx, fop, fi), err

}

func (s *SQLCommon) UpdateTransaction(ctx context.Context, namespace string, id *fftypes.UUID, update ffapi.Update) (err error) {

	ctx, tx, autoCommit, err := s.BeginOrUseTx(ctx)
	if err != nil {
		return err
	}
	defer s.RollbackTx(ctx, tx, autoCommit)

	query, err := s.BuildUpdate(sq.Update(transactionsTable), update, transactionFilterFieldMap)
	if err != nil {
		return err
	}
	query = query.Where(sq.Eq{"id": id, "namespace": namespace})

	_, err = s.UpdateTx(ctx, transactionsTable, tx, query, nil /* no change evnents for filter based updates */)
	if err != nil {
		return err
	}

	return s.CommitTx(ctx, tx, autoCommit)
}
