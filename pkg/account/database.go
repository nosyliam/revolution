package account

import (
	"encoding/json"
	"github.com/pkg/errors"
	"os"
	"path/filepath"
)

type Account struct {
	Token   string   `json:"token"`
	Servers []Server `json:"servers"`

	cookie string
}

func (a *Account) Load() error {
	return nil
}

func (a *Account) Refresh(secret []byte) error {
	return nil
}

func (a *Account) GenerateJoinUrl() (string, error) {
	return "", nil
}

type Database struct {
	Encrypt  bool      `json:"encrypt"`
	Accounts []Account `json:"accounts"`

	path   string
	secret []byte
}

func (d *Database) AddAccount() {

}

func (d *Database) Save() error {
	if d.path == "" {
		cwd, err := os.Getwd()
		if err != nil {
			return errors.Wrap(err, "failed to get working directory")
		}
		d.path = filepath.Join(cwd, "accounts.json")
	}

	f, err := os.OpenFile(d.path, os.O_RDWR|os.O_CREATE|os.O_TRUNC, 0755)
	if err != nil {
		return errors.Wrap(err, "failed to open account database")
	}

	data, err := json.Marshal(d)
	if err != nil {
		return errors.Wrap(err, "failed to marshal account database")
	}

	_, err = f.Write(data)
	if err != nil {
		return errors.Wrap(err, "failed to write to account database")
	}

	_ = f.Close()
}

func (d *Database) Load(path string, secret []byte) error {

}
