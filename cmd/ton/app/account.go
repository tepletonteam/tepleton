package app

import (
	sdk "github.com/tepleton/tepleton-sdk/types"
	"github.com/tepleton/tepleton-sdk/x/auth"
)

// State to Unmarshal
type GenesisState struct {
	Accounts []GenesisAccount `json:"accounts"`
}

// GenesisAccount doesn't need pubkey or sequence
type GenesisAccount struct {
	Address sdk.Address `json:"address"`
	Coins   sdk.Coins   `json:"coins"`
}

func NewGenesisAccount(acc auth.BaseAccount) GenesisAccount {
	return GenesisAccount{
		Address: acc.Address,
		Coins:   acc.Coins,
	}
}

// convert GenesisAccount to GaiaAccount
func (ga *GenesisAccount) ToAccount() (acc auth.BaseAccount) {
	return auth.BaseAccount{
		Address: ga.Address,
		Coins:   ga.Coins.Sort(),
	}
}
