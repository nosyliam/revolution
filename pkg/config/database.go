package config

import (
	"fmt"
	"github.com/pkg/errors"
)

type Account struct {
	Name           string  `json:"name"`
	ServerID       *string `json:"serverID,omitempty"`
	LinkCode       *string `json:"linkCode"`
	Invalid        bool    `json:"invalid"`
	WindowConfigID *string `json:"windowConfigID,omitempty"`

	db *AccountDatabase
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
	configFile
	Accounts []*Account `json:"accounts"`
	Servers  []*Server  `json:"servers"`

	accountMap map[string]*Account
}

func (d *AccountDatabase) Add(code string) {
	account, err := NewAccount(d, code)
	if err != nil {
		// TODO: send error to UI
	}
	d.Accounts = append(d.Accounts, account)
	d.accountMap[code] = account
}

func (d *AccountDatabase) Delete(name string) {
	if _, ok := d.accountMap[name]; !ok {
		return
	}
	delete(d.accountMap, name)
	for i, account := range d.Accounts {
		if account.Name == name {
			d.Accounts = append(d.Accounts[:i], d.Accounts[i+1:]...)
		}
	}
}

func (d *AccountDatabase) Get(name string) *Account {
	if account, ok := d.accountMap[name]; !ok {
		return nil
	} else {
		return account
	}
}

func (d *AccountDatabase) RegisterListeners() {

}

func (d *AccountDatabase) Default() {}

func NewAccountDatabase() (*AccountDatabase, error) {
	db := &AccountDatabase{configFile: configFile{path: "accounts.json", format: JSON}}
	if err := db.load(); err != nil {
		return nil, errors.Wrap(err, "Failed to load account database")
	}
	db.accountMap = make(map[string]*Account)
	for _, account := range db.Accounts {
		db.accountMap[account.Name] = account
		account.db = db
	}
	return db, nil
}
