// Copyright 2024 Google LLC
//
// Licensed under the Apache License, Version 2.0 (the "License");
// you may not use this file except in compliance with the License.
// You may obtain a copy of the License at
//
//     https://www.apache.org/licenses/LICENSE-2.0
//
// Unless required by applicable law or agreed to in writing, software
// distributed under the License is distributed on an "AS IS" BASIS,
// WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
// See the License for the specific language governing permissions and
// limitations under the License.

package accounts

import (
	"errors"
	"net/http"

	appOAuth2 "github.com/google/cloud-android-orchestration/pkg/app/oauth2"
)

const BasicAMType AMType = "basic"

// Implements the AccountManager interface using HTTP Basic Authentication,
// where the username and password are provided in the HTTP request header.
type BasicAccountManager struct {}

func NewBasicAccountManager() *BasicAccountManager {
	return &BasicAccountManager{}
}

func (m *BasicAccountManager) UserFromRequest(r *http.Request) (User, error) {
	return userFromRequest(r)
}

func (m *BasicAccountManager) OnOAuth2Exchange(w http.ResponseWriter, r *http.Request, tk appOAuth2.IDTokenClaims) (User, error) {
	rUser, err := userFromRequest(r)
	if err != nil {
		return nil, err
	}
	user, ok := tk["user"]
	if !ok {
		return nil, errors.New("no user in id token")
	}
	tkUser, ok := user.(string)
	if !ok {
		return nil, errors.New("malformed user in id token")
	}
	if rUser.Username() != tkUser {
		return nil, errors.New("logged in user doesn't match oauth2 user")
	}
	return rUser, nil
}

type BasicUser struct {
	username string
}

func (u *BasicUser) Username() string {
	return u.username
}

func userFromRequest(r *http.Request) (*BasicUser, error) {
	// TODO: verify the password
	username, _, ok := r.BasicAuth()
	if !ok {
		return nil, errors.New("cannot get username from the http request")
	}
	return &BasicUser{username}, nil
}
