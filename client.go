// Copyright 2025- The sacloud/secretmanager-api-go authors
//
// Licensed under the Apache License, Version 2.0 (the "License");
// you may not use this file except in compliance with the License.
// You may obtain a copy of the License at
//
//      http://www.apache.org/licenses/LICENSE-2.0
//
// Unless required by applicable law or agreed to in writing, software
// distributed under the License is distributed on an "AS IS" BASIS,
// WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
// See the License for the specific language governing permissions and
// limitations under the License.

package secretmanager

import (
	"context"
	"fmt"
	"runtime"

	client "github.com/sacloud/api-client-go"
	v1 "github.com/sacloud/secretmanager-api-go/apis/v1"
)

// DefaultAPIRootURL デフォルトのAPIルートURL
const DefaultAPIRootURL = "https://secure.sakura.ad.jp/cloud/zone/tk1a/api/cloud/1.1"

// const DefaultAPIRootURL = "https://secure.sakura.ad.jp/cloud-test/zone/is1y/api/cloud/1.1"

// UserAgent APIリクエスト時のユーザーエージェント
var UserAgent = fmt.Sprintf(
	"secretmanager-api-go/%s (%s/%s; +https://github.com/sacloud/secretmanager-api-go) %s",
	Version,
	runtime.GOOS,
	runtime.GOARCH,
	client.DefaultUserAgent,
)

// SecuritySourceはOpenAPI定義で使用されている認証のための仕組み。api-client-goが処理するので、ogen用はダミーで誤魔化す
type DummySecuritySource struct {
	Username string
	Password string
}

func (ss DummySecuritySource) BasicAuth(ctx context.Context, operationName v1.OperationName) (v1.BasicAuth, error) {
	return v1.BasicAuth{Username: ss.Username, Password: ss.Password, Roles: nil}, nil
}

func NewClient() (*v1.Client, error) {
	return NewClientWithApiUrl(DefaultAPIRootURL)
}

func NewClientWithApiUrl(apiUrl string) (*v1.Client, error) {
	c, err := client.NewClient(apiUrl, client.WithUserAgent(UserAgent))
	if err != nil {
		return nil, err
	}

	v1Client, err := v1.NewClient(c.ServerURL(), DummySecuritySource{Username: "", Password: ""}, v1.WithClient(c.NewHttpRequestDoer()))
	if err != nil {
		return nil, fmt.Errorf("failed to create secretmanager client: %w", err)
	}

	return v1Client, nil
}
