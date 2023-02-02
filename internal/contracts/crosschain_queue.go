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

package contracts

import (
	"bytes"
	"encoding/json"
	"log"
	"net/http"

	"github.com/hyperledger/firefly/pkg/core"
)

type CrossChainManager struct{}

type ConfirmPrimaryTx struct{}

func NewCrossChainManager() *CrossChainManager {
	ccm := new(CrossChainManager)
	return ccm
}

func (ccm *CrossChainManager) ProcessTx(ppTx core.PreparePrimaryTx) {
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
		log.Fatal(err)
	}
	url := ppTx.URL + "/api/v1/namespaces/default/apis/cross-network/invoke/doNetwork"

	_, err = http.Post(url, "application/json", bytes.NewBuffer(jsonData)) //nolint
	if err != nil {
		log.Fatal(err)
	}
}
