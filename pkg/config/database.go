package config

import (
	"fmt"
	"github.com/pkg/errors"
)

type Account struct {
	Name           string `yaml:"name" key:"true"`
	ServerID       string `yaml:"serverID,omitempty"`
	LinkCode       string `yaml:"linkCode,omitempty"`
	Invalid        bool   `yaml:"invalid,omitempty"`
	WindowConfigID string `yaml:"windowConfigID,omitempty"`
}

func (a *Account) Load() error {
	// Load account information & private servers from token
	return nil
}

func (a *Account) Refresh() error {
	return nil
}

func (a *Account) GenerateJoinUrl(ignoreLink bool) (string, error) {
	if a.Invalid {
		return "", errors.New(fmt.Sprintf("The session for account %s has expired", a.Name))
	}

	//

	return "", nil
}

func NewAccount(db *AccountDatabase, code string) (*Account, error) {
	return nil, nil
}

type AccountDatabase struct {
	Accounts *List[Account] `yaml:"accounts"`
	Servers  *List[Server]  `yaml:"servers"`
}

func (d *AccountDatabase) Add(code string) {
	/*
		account, err := NewAccount(d, code)
		if err != nil {
			// TODO: send error to UI
		}
		d.accountMap[code] = account*/
}

func (d *AccountDatabase) Delete(name string) {
	/*
		if _, ok := d.accountMap[name]; !ok {
			return
		}
		delete(d.accountMap, name)
		for i, account := range d.Accounts {
			if account.Name == name {
				d.Accounts = append(d.Accounts[:i], d.Accounts[i+1:]...)
			}
		}*/
}

func (d *AccountDatabase) Get(name string) *Account {
	/*
		if account, ok := d.accountMap[name]; !ok {
			return nil
		} else {
			return account
		}*/
	return nil
}

func NewAccountDatabase() (*AccountDatabase, error) {
	/*
		db := &AccountDatabase{configFile: configFile{path: "accounts.json", format: JSON}}
		if err := db.load(); err != nil {
			return nil, errors.Wrap(err, "Failed to load account database")
		}
		db.accountMap = make(map[string]*Account)
		for _, account := range db.Accounts {
			db.accountMap[account.Name] = account
			account.db = db
		}

	*/
	return nil, nil
}
