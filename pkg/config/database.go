package config

type Account struct {
	Token      string   `json:"token"`
	Servers    []Server `json:"servers"`
	PresetName string   `json:"presetName"`

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
	configFile
	Encrypt  bool      `json:"encrypt"`
	Accounts []Account `json:"accounts"`
}

func (d *Database) Decrypt() error {
	return nil
}

func (d *Database) AddAccount() {

}

func NewAccountDatabase() *Database {
	return &Database{configFile: configFile{path: "accounts.json", format: JSON}}
}
