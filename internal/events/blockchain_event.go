// Copyright © 2023 Kaleido, Inc.
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

package events

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"net/http"

	"github.com/hyperledger/firefly-common/pkg/fftypes"
	"github.com/hyperledger/firefly-common/pkg/log"
	"github.com/hyperledger/firefly/pkg/blockchain"
	"github.com/hyperledger/firefly/pkg/core"
)

func buildBlockchainEvent(ns string, subID *fftypes.UUID, event *blockchain.Event, tx *core.BlockchainTransactionRef) *core.BlockchainEvent {
	fmt.Printf("[CrossChain]: The Event is %s \n", event.Name)

	out, err := json.Marshal(event.Output)
	if err != nil {
		fmt.Println("[CrossChain]: Event -> Error: ", err)
	}
	fmt.Printf("[CrossChain]: Output: %s \n", string(out))

	ev := &core.BlockchainEvent{
		ID:         fftypes.NewUUID(),
		Namespace:  ns,
		Listener:   subID,
		Source:     event.Source,
		ProtocolID: event.ProtocolID,
		Name:       event.Name,
		Output:     event.Output,
		Info:       event.Info,
		Timestamp:  event.Timestamp,
	}

	if event.Name == "PrimaryTxStatus" {
		fmt.Printf("[CrossChain]: PrimaryTxStatus on: txId: %s, satus: %s \n", event.Output.GetString("txId"), event.Output.GetInteger("status").String())
		if event.Output.GetInteger("status").String() == "4" {
			fmt.Printf("[CrossChain]: PrimaryTxStatus -> changeStatusNetwork on: networkUrl: %s, txId: %s \n", event.Output.GetString("networkUrl"), event.Output.GetString("txId"))
			changeStatusNetwork(event.Output.GetString("networkUrl"), event.Output.GetString("txId"), 4)
			fmt.Printf("[CrossChain]: Done PrimaryTxStatus -> changeStatusNetwork on: networkUrl: %s, txId: %s \n", event.Output.GetString("networkUrl"), event.Output.GetString("txId"))

			fmt.Printf("[CrossChain]: PrimaryTxStatus -> processConfirmTx on: networkUrl: %s, txId: %s \n", event.Output.GetString("networkUrl"), event.Output.GetString("txId"))
			err := processConfirmTx(event.Output.GetString("txId"), event.Output.GetString("networkUrl"), true)
			if err != nil {
				fmt.Println("[CrossChain]: PrimaryTxStatus -> processConfirmTx -> Error: ", err)
			}
			fmt.Printf("[CrossChain]: Done PrimaryTxStatus -> processConfirmTx on: networkUrl: %s, txId: %s \n", event.Output.GetString("networkUrl"), event.Output.GetString("txId"))
		}
		fmt.Printf("[CrossChain]: Done PrimaryTxStatus on: %s \n", event.Output.GetString("status"))
	}
	if event.Name == "NetworkTxStatus" {
		fmt.Printf("[CrossChain]: NetworkTxStatus on: txId: %s, status: %s \n", event.Output.GetString("txId"), event.Output.GetInteger("status").String())
		switch status := event.Output.GetInteger("status").String(); status {
		case "2":
			fmt.Printf("[CrossChain]: NetworkTxStatus -> NetworkTransactionStatusType.NETWORK_TRANSACTION_STARTED -> changeStatusPrimary on: primaryNetworkUrl %s, txId: %s \n", event.Output.GetString("primaryNetworkUrl"), event.Output.GetString("txId"))
			changeStatusPrimary(event.Output.GetString("primaryNetworkUrl"), event.Output.GetString("txId"), 2)
			fmt.Printf("[CrossChain]: Done NetworkTxStatus -> changeStatusPrimary on: primaryNetworkUrl %s, txId: %s \n", event.Output.GetString("primaryNetworkUrl"), event.Output.GetString("txId"))
		case "3":
			fmt.Printf("[CrossChain]: NetworkTxStatus -> NetworkTransactionStatusType.NETWORK_TRANSACTION_PREPARED -> changeStatusPrimary on: primaryNetworkUrl %s, txId: %s \n", event.Output.GetString("primaryNetworkUrl"), event.Output.GetString("txId"))
			changeStatusPrimary(event.Output.GetString("primaryNetworkUrl"), event.Output.GetString("txId"), 3)
			fmt.Printf("[CrossChain]: Done NetworkTxStatus -> NetworkTransactionStatusType.NETWORK_TRANSACTION_PREPARED -> changeStatusPrimary on: primaryNetworkUrl %s, txId: %s \n", event.Output.GetString("primaryNetworkUrl"), event.Output.GetString("txId"))

			fmt.Printf("[CrossChain]: NetworkTxStatus -> NetworkTransactionStatusType.NETWORK_TRANSACTION_PREPARED -> processConfirmTx on: primaryNetworkUrl %s, txId: %s \n", event.Output.GetString("primaryNetworkUrl"), event.Output.GetString("txId"))
			err := processConfirmTx(event.Output.GetString("txId"), event.Output.GetString("primaryNetworkUrl"), false)
			if err != nil {
				fmt.Println("[CrossChain]: NetworkTxStatus -> NetworkTransactionStatusType.NETWORK_TRANSACTION_PREPARED -> processConfirmTx -> Error: ", err)
			}
			fmt.Printf("[CrossChain]: Done NetworkTxStatus -> NetworkTransactionStatusType.NETWORK_TRANSACTION_PREPARED -> processConfirmTx on: primaryNetworkUrl %s, txId: %s \n", event.Output.GetString("primaryNetworkUrl"), event.Output.GetString("txId"))
		case "5":
			fmt.Printf("[CrossChain]: NetworkTxStatus -> NetworkTransactionStatusType.NETWORK_TRANSACTION_COMMITTED -> changeStatusPrimary on: primaryNetworkUrl %s, txId: %s \n", event.Output.GetString("primaryNetworkUrl"), event.Output.GetString("txId"))
			changeStatusPrimary(event.Output.GetString("primaryNetworkUrl"), event.Output.GetString("txId"), 5)
			fmt.Printf("[CrossChain]: Done NetworkTxStatus -> NetworkTransactionStatusType.NETWORK_TRANSACTION_COMMITTED -> changeStatusPrimary on: primaryNetworkUrl %s, txId: %s \n", event.Output.GetString("primaryNetworkUrl"), event.Output.GetString("txId"))
		}
		fmt.Printf("[CrossChain]: Done NetworkTxStatus on: %s \n", event.Output.GetString("status"))
	}

	if event.Name == "PreparePrimaryTransaction" {
		fmt.Printf("[CrossChain]: PreparePrimaryTransaction: TxID: %s, PrimaryNetworkID: %s, NetworkID: %s, URL: %s, InvocationID: %s, Args: %s \n", event.Output.GetString("txId"), event.Output.GetString("primaryNetworkId"), event.Output.GetString("networkId"), event.Output.GetString("url"), event.Output.GetString("invocationId"), event.Output.GetString("args"))
		ppTx := core.PreparePrimaryTx{
			TxID:             event.Output.GetString("txId"),
			PrimaryNetworkID: event.Output.GetString("primaryNetworkId"),
			NetworkID:        event.Output.GetString("networkId"),
			URL:              event.Output.GetString("url"),
			InvocationID:     event.Output.GetString("invocationId"),
			Args:             event.Output.GetString("args"),
		}
		fmt.Printf("[CrossChain]: PreparePrimaryTransaction -> processPreparePrimaryTransaction: TxID: %s, PrimaryNetworkID: %s, NetworkID: %s, URL: %s, InvocationID: %s, Args: %s \n", event.Output.GetString("txId"), event.Output.GetString("primaryNetworkId"), event.Output.GetString("networkId"), event.Output.GetString("url"), event.Output.GetString("invocationId"), event.Output.GetString("args"))
		err := processPreparePrimaryTransaction(ppTx)
		if err != nil {
			fmt.Println("[CrossChain]: PreparePrimaryTransaction -> processPreparePrimaryTransaction -> Error: ", err)
		}
		fmt.Printf("[CrossChain]: Done PreparePrimaryTransaction -> processPreparePrimaryTransaction: TxID: %s, PrimaryNetworkID: %s, NetworkID: %s, URL: %s, InvocationID: %s, Args: %s \n", event.Output.GetString("txId"), event.Output.GetString("primaryNetworkId"), event.Output.GetString("networkId"), event.Output.GetString("url"), event.Output.GetString("invocationId"), event.Output.GetString("args"))
	}
	if event.Name == "ConfirmNetworkTransaction" {
		fmt.Printf("[CrossChain]: ConfirmNetworkTransaction: TxID: %s, boolean success: %t, data: %s \n", event.Output.GetString("txId"), event.Output.GetBool("success"), event.Output.GetString("data"))
		err := prosessConfirmNetworkTransaction(event.Output.GetString("txId"), event.Output.GetBool("success"), event.Output.GetString("data"), event.Output.GetString("primaryNetworkUrl"))
		if err != nil {
			fmt.Println("[CrossChain]: ConfirmNetworkTransaction -> prosessConfirmNetworkTransaction -> Error: ", err)
		}
		fmt.Printf("[CrossChain]: Done ConfirmNetworkTransaction: TxID: %s, boolean success: %t, data: %s \n", event.Output.GetString("txId"), event.Output.GetBool("success"), event.Output.GetString("data"))
	}
	if tx != nil {
		ev.TX = *tx
	}
	return ev
}

func (em *eventManager) getChainListenerByProtocolIDCached(ctx context.Context, protocolID string) (*core.ContractListener, error) {
	return em.getChainListenerCached(fmt.Sprintf("pid:%s", protocolID), func() (*core.ContractListener, error) {
		return em.database.GetContractListenerByBackendID(ctx, em.namespace.Name, protocolID)
	})
}

func (em *eventManager) getChainListenerCached(cacheKey string, getter func() (*core.ContractListener, error)) (*core.ContractListener, error) {

	if cachedValue := em.chainListenerCache.Get(cacheKey); cachedValue != nil {
		return cachedValue.(*core.ContractListener), nil
	}
	listener, err := getter()
	if listener == nil || err != nil {
		return nil, err
	}
	em.chainListenerCache.Set(cacheKey, listener)
	return listener, err
}

func (em *eventManager) getTopicForChainListener(listener *core.ContractListener) string {
	if listener == nil {
		return core.SystemBatchPinTopic
	}
	var topic string
	if listener != nil && listener.Topic != "" {
		topic = listener.Topic
	} else {
		topic = listener.ID.String()
	}
	return topic
}

func (em *eventManager) maybePersistBlockchainEvent(ctx context.Context, chainEvent *core.BlockchainEvent, listener *core.ContractListener) error {
	if existing, err := em.txHelper.InsertOrGetBlockchainEvent(ctx, chainEvent); err != nil {
		return err
	} else if existing != nil {
		log.L(ctx).Debugf("Ignoring duplicate blockchain event %s", chainEvent.ProtocolID)
		// Return the ID of the existing event
		chainEvent.ID = existing.ID
		return nil
	}
	topic := em.getTopicForChainListener(listener)
	ffEvent := core.NewEvent(core.EventTypeBlockchainEventReceived, chainEvent.Namespace, chainEvent.ID, chainEvent.TX.ID, topic)
	err := em.database.InsertEvent(ctx, ffEvent)
	return err
}

func (em *eventManager) emitBlockchainEventMetric(event *blockchain.Event) {
	if em.metrics.IsMetricsEnabled() && event.Location != "" && event.Signature != "" {
		em.metrics.BlockchainEvent(event.Location, event.Signature)
	}
}

func (em *eventManager) BlockchainEvent(event *blockchain.EventWithSubscription) error {
	return em.retry.Do(em.ctx, "persist blockchain event", func(attempt int) (bool, error) {
		err := em.database.RunAsGroup(em.ctx, func(ctx context.Context) error {
			listener, err := em.getChainListenerByProtocolIDCached(ctx, event.Subscription)
			if err != nil {
				return err
			}
			if listener == nil {
				log.L(ctx).Warnf("Event received from unknown subscription %s", event.Subscription)
				return nil // no retry
			}
			if listener.Namespace != em.namespace.Name {
				log.L(em.ctx).Debugf("Ignoring blockchain event from different namespace '%s'", listener.Namespace)
				return nil
			}
			listener.Namespace = em.namespace.Name

			chainEvent := buildBlockchainEvent(listener.Namespace, listener.ID, &event.Event, &core.BlockchainTransactionRef{
				BlockchainID: event.BlockchainTXID,
			})

			if err := em.maybePersistBlockchainEvent(ctx, chainEvent, listener); err != nil {
				return err
			}
			em.emitBlockchainEventMetric(&event.Event)
			return nil
		})
		return err != nil, err
	})
}

func changeStatusPrimary(primaryURL string, txID string, _status int) {
	fmt.Printf("[CrossChain]: changeStatusPrimary on: primaryURL: %s, txId: %s, status: %d, \n", primaryURL, txID, _status)
	values := map[string]map[string]interface{}{
		"input": {
			"txId":    txID,
			"_status": _status,
		},
	}
	jsonData, _ := json.Marshal(values)
	url := fmt.Sprintf("%s/api/v1/namespaces/default/apis/cross-chain/invoke/changeStatus", primaryURL)
	_, err := http.Post(url, "application/json", bytes.NewBuffer(jsonData)) //nolint
	if err != nil {
		fmt.Println("[CrossChain]: changeStatusPrimary -> Error: ", err)
	}
	fmt.Printf("[CrossChain]: Done changeStatusPrimary on: primaryURL: %s, txId: %s, status: %d, \n", primaryURL, txID, _status)
}

func changeStatusNetwork(networkURL string, txID string, _status int) {
	fmt.Printf("[CrossChain]: changeStatusNetwork on: networkURL: %s, txId: %s, status: %d, \n", networkURL, txID, _status)
	values := map[string]map[string]interface{}{
		"input": {
			"txId":    txID,
			"_status": _status,
		},
	}
	jsonData, _ := json.Marshal(values)
	url := fmt.Sprintf("%s/api/v1/namespaces/default/apis/cross-network/invoke/changeStatus", networkURL)
	_, err := http.Post(url, "application/json", bytes.NewBuffer(jsonData)) //nolint
	if err != nil {
		fmt.Println("[CrossChain]: changeStatusNetwork -> Error: ", err)
	}
	fmt.Printf("[CrossChain]: Done changeStatusNetwork on: networkURL: %s, txId: %s, status: %d, \n", networkURL, txID, _status)
}

func processPreparePrimaryTransaction(ppTx core.PreparePrimaryTx) error {
	fmt.Printf("[CrossChain]: processPreparePrimaryTransaction \n")

	values := map[string]map[string]interface{}{
		"input": {
			"txId":             ppTx.TxID,
			"primaryNetworkId": ppTx.PrimaryNetworkID,
			"networkId":        ppTx.NetworkID,
			"invocationId":     ppTx.InvocationID,
			"args":             ppTx.Args,
		},
	}
	jsonData, err := json.Marshal(values)
	if err != nil {
		fmt.Println("[CrossChain]: processPreparePrimaryTransaction -> Error: ", err)
	}
	url := fmt.Sprintf("%s/api/v1/namespaces/default/apis/cross-network/invoke/doNetwork", ppTx.URL)

	_, err = http.Post(url, "application/json", bytes.NewBuffer(jsonData)) //nolint
	if err != nil {
		fmt.Println("[CrossChain]: processPreparePrimaryTransaction -> Error: ", err)
	}
	fmt.Printf("[CrossChain]: Done processPreparePrimaryTransaction \n")
	return err
}

func processConfirmTx(txID string, typeURL string, isNetwork bool) error {
	fmt.Printf("[CrossChain]: processConfirmTx \n")
	values := map[string]map[string]interface{}{
		"input": {
			"txId": txID,
		},
	}
	jsonData, _ := json.Marshal(values)

	url := fmt.Sprintf("%s/api/v1/namespaces/default/apis/%s/invoke/%s", typeURL, parseURL(isNetwork, true), parseURL(isNetwork, false))
	_, err := http.Post(url, "application/json", bytes.NewBuffer(jsonData)) //nolint
	if err != nil {
		fmt.Println("[CrossChain]: processConfirmTx -> Error: ", err)
	}
	fmt.Printf("[CrossChain]: Done processConfirmTx \n")
	return err
}

func parseURL(isNetwork bool, first bool) string {
	switch {
	case isNetwork && first:
		return "cross-network"
	case isNetwork && !first:
		return "confirmDoNetwork"
	case !isNetwork && first:
		return "cross-chain"
	default:
		return "confirmDoCross"
	}
}

func prosessConfirmNetworkTransaction(txID string, success bool, data string, primaryNetworkUrl string) error {
	fmt.Printf("[CrossChain]: prosessConfirmNetworkTransaction \n")
	values := map[string]map[string]interface{}{
		"input": {
			"txId": txID,
			"data": nil,
		},
	}
	if success {
		values["input"]["data"] = data
	}
	jsonData, _ := json.Marshal(values)

	url := fmt.Sprintf("%s/api/v1/namespaces/default/apis/cross-chain/invoke/finishPrimaryTransaction", primaryNetworkUrl)
	_, err := http.Post(url, "application/json", bytes.NewBuffer(jsonData)) //nolint
	if err != nil {
		fmt.Println("[CrossChain]: prosessConfirmNetworkTransaction -> Error: ", err)
	}
	fmt.Printf("[CrossChain]: prosessConfirmNetworkTransaction \n")
	return err
}
