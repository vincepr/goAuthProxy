package main

import "fmt"

/*
* interface could make switching to config-file or database
* or supporting multiple dbs easy.
 */
 type Storage interface{
	AddAccount() error
	GetAccountByName(string) (*Account, error)
}

// quick account storage that stores accounts just in memory without any file-reads
type AccountStorage struct{
	Accounts map[string]Account
}

// maps our users that can log in and consume our content
type Account struct {
	Name string 		`json:"name"`		// unique username garantteed by map
	PasswordHash []byte `json:"-"`
	IsAdmin bool 		`json:"isAdmin"`
}

// Iplementation of our default Storage type, thats basically just a map living in memory atm.
func NewAccountStorage() (*AccountStorage){
	storage := make(map[string]Account)
	return &AccountStorage{
		Accounts: storage,
	}
}

func (st *AccountStorage) AddAccount(newAcc *Account) (error){
	// check if username already exists
	if _, ok := st.Accounts[newAcc.Name]; !ok{
		st.Accounts[newAcc.Name] = *newAcc
		return nil
	}
	return fmt.Errorf("User already exists")
}

func (st *AccountStorage) GetAccountByName(name string) (*Account, error){
	if val, ok := st.Accounts[name]; ok{
		return &val, nil
	}
	return nil, fmt.Errorf("User not found")
}
