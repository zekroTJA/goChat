package main

import (
	"crypto/sha256"
	"encoding/base64"
	"encoding/json"
	"os"
	"strings"
)

type Acc struct {
	Username string `json:"username"`
	Passhash string `json:"passhash"`
}

type AccountManager struct {
	FileName string
	accounts map[string]*Acc
}

func NewAccountManager(filename string) (*AccountManager, error) {
	accMgr := AccountManager{accounts: make(map[string]*Acc),
		FileName: filename,
	}
	file, err := os.Open(filename)
	defer file.Close()
	if !os.IsNotExist(err) {
		tmpMap := make(map[string]*Acc)
		decoder := json.NewDecoder(file)
		err = decoder.Decode(&tmpMap)
		if err != nil {
			return nil, err
		}
		accMgr.accounts = tmpMap
	}
	return &accMgr, nil
}

func (mgr *AccountManager) Register(username, password string) {
	passhash := createHash(password)
	mgr.accounts[strings.ToLower(username)] = &Acc{
		Username: username,
		Passhash: passhash,
	}
}

func (mgr *AccountManager) Get(username string) *Acc {
	username = strings.ToLower(username)
	acc := mgr.accounts[username]
	return acc
}

func (mgr *AccountManager) Check(username, password string) (*Acc, bool) {
	acc := mgr.Get(username)
	if acc == nil {
		return nil, false
	}
	valid := acc.Passhash == createHash(password)
	return acc, valid
}

func (mgr *AccountManager) Save() error {
	file, err := os.Create(mgr.FileName)
	defer file.Close()
	if err != nil {
		return err
	}
	encoder := json.NewEncoder(file)
	encoder.SetIndent("", "  ")
	err = encoder.Encode(mgr.accounts)
	return err
}

func createHash(data string) string {
	hasher := sha256.New()
	hasher.Write([]byte(data))
	return base64.URLEncoding.EncodeToString(hasher.Sum(nil))
}
