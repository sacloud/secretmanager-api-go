// Copyright 2025- The sacloud/secretmanager-api-go authors
//
// Licensed under the Apache License, Version 2.0 (the "License");
// you may not use this file except in compliance with the License.
// You may obtain a copy of the License at
//
//	http://www.apache.org/licenses/LICENSE-2.0
//
// Unless required by applicable law or agreed to in writing, software
// distributed under the License is distributed on an "AS IS" BASIS,
// WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
// See the License for the specific language governing permissions and
// limitations under the License.

package secretmanager_test

import (
	"context"
	"fmt"
	"os"
	"sort"
	"strconv"

	sm "github.com/sacloud/secretmanager-api-go"
	v1 "github.com/sacloud/secretmanager-api-go/apis/v1"
)

var requriedEnvs = []string{
	"SAKURACLOUD_KMS_KEY_ID",
	"SAKURACLOUD_ACCESS_TOKEN",
	"SAKURACLOUD_ACCESS_TOKEN_SECRET",
}

func checkEnvs() {
	for _, env := range requriedEnvs {
		if os.Getenv(env) == "" {
			panic(env + " is not set")
		}
	}
}

func ExampleVaultAPI() {
	checkEnvs()

	client, err := sm.NewClient()
	if err != nil {
		panic(err)
	}

	ctx := context.Background()
	keyId := os.Getenv("SAKURACLOUD_KMS_KEY_ID") // コンパネやkms-api-goなどで取得
	vaultOp := sm.NewVaultOp(client)

	resCreate, err := vaultOp.Create(ctx, v1.CreateVault{
		Name:        "vault from go",
		Description: v1.NewOptString("vault from go client"),
		KmsKeyID:    keyId,
		Tags:        []string{"App", "Vault"},
	})
	if err != nil {
		panic(err)
	}
	fmt.Println(resCreate.Name)

	resList, err := vaultOp.List(ctx)
	if err != nil {
		panic(err)
	}

	// 生成順に帰ってくるとは限らないのでテストのためにソート
	sort.Slice(resList, func(i, j int) bool { return resList[i].ID < resList[j].ID })
	for _, vault := range resList {
		if vault.ID == resCreate.ID {
			fmt.Println(vault.Description.Value)
		}
	}

	_, err = vaultOp.Update(ctx, resCreate.ID, v1.Vault{
		Name:        "vault from go 2",
		Description: v1.NewOptString("vault from go client 2"),
		KmsKeyID:    keyId,
		Tags:        []string{"Test"},
	})
	if err != nil {
		panic(err)
	}

	resRead, err := vaultOp.Read(ctx, resCreate.ID)
	if err != nil {
		panic(err)
	}

	fmt.Println(resRead.Name)

	err = vaultOp.Delete(ctx, resCreate.ID)
	if err != nil {
		panic(err)
	}

	fmt.Println(resRead.Description.Value)
	// Output:
	// vault from go
	// vault from go client
	// vault from go 2
	// vault from go client 2
}

func ExampleSecretAPI() {
	checkEnvs()

	client, err := sm.NewClient()
	if err != nil {
		panic(err)
	}

	ctx := context.Background()
	keyId := os.Getenv("SAKURACLOUD_KMS_KEY_ID") // コンパネやkms-api-goなどで取得
	vaultOp := sm.NewVaultOp(client)

	vault, err := vaultOp.Create(ctx, v1.CreateVault{
		Name:        "vault for secret test",
		Description: v1.NewOptString("vault for secret test"),
		KmsKeyID:    keyId,
		Tags:        []string{"Test"},
	})
	if err != nil {
		panic(err)
	}

	secOp := sm.NewSecretOp(client, vault.ID)

	for i := 0; i < 2; i++ {
		resCreate, err := secOp.Create(ctx, v1.CreateSecret{
			Name:  "Sec1",
			Value: "SecretValue" + strconv.Itoa(i),
		})
		if err != nil {
			panic(err)
		}
		fmt.Println("version: " + strconv.Itoa(resCreate.LatestVersion))
	}
	resCreate, err := secOp.Create(ctx, v1.CreateSecret{
		Name:  "Sec2",
		Value: "SV22",
	})
	if err != nil {
		panic(err)
	}
	fmt.Println("version: " + strconv.Itoa(resCreate.LatestVersion))

	resList, err := secOp.List(ctx)
	if err != nil {
		panic(err)
	}

	sort.Slice(resList, func(i, j int) bool { return resList[i].Name < resList[j].Name })
	for _, sec := range resList {
		fmt.Println("name: " + sec.Name + ", version: " + strconv.Itoa(sec.LatestVersion))
	}

	for i := 0; i < 2; i++ {
		resUn, err := secOp.Unveil(ctx, v1.Unveil{
			Name:    "Sec1",
			Version: v1.NewOptNilInt(i + 1),
		})
		if err != nil {
			panic(err)
		}
		fmt.Println("value: " + resUn.Value)
	}
	fmt.Println("test end")
	// Output:
	// version: 1
	// version: 2
	// version: 1
	// name: Sec1, version: 2
	// name: Sec2, version: 1
	// value: SecretValue0
	// value: SecretValue1
	// test end

	err = secOp.Delete(ctx, v1.DeleteSecret{Name: "Sec1"})
	if err != nil {
		panic(err)
	}
	err = secOp.Delete(ctx, v1.DeleteSecret{Name: "Sec2"})
	if err != nil {
		panic(err)
	}

	err = vaultOp.Delete(ctx, vault.ID)
	if err != nil {
		panic(err)
	}
}
