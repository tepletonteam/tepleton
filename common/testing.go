//functions used in testing throughout
package common

import (
	"github.com/tepleton/basecoin/types"
	. "github.com/tepleton/go-common"
	"github.com/tepleton/go-crypto"
)

// Creates a PrivAccount from secret.
// The amount is not set.
func PrivAccountFromSecret(secret string) types.PrivAccount {
	privKey := crypto.GenPrivKeyEd25519FromSecret([]byte(secret))
	privAccount := types.PrivAccount{
		PrivKey: privKey,
		Account: types.Account{
			PubKey:   privKey.PubKey(),
			Sequence: 0,
		},
	}
	return privAccount
}

// Make `num` random accounts
func RandAccounts(num int, minAmount int64, maxAmount int64) []types.PrivAccount {
	privAccs := make([]types.PrivAccount, num)
	for i := 0; i < num; i++ {

		balance := minAmount
		if maxAmount > minAmount {
			balance += RandInt64() % (maxAmount - minAmount)
		}

		privKey := crypto.GenPrivKeyEd25519()
		pubKey := privKey.PubKey()
		privAccs[i] = types.PrivAccount{
			PrivKey: privKey,
			Account: types.Account{
				PubKey:   pubKey,
				Sequence: 0,
				Balance:  types.Coins{types.Coin{"", balance}},
			},
		}
	}

	return privAccs
}

//make input term for the AppTx or SendTx Types
func MakeInput(pubKey crypto.PubKey, coins types.Coins, sequence int) types.TxInput {
	input := types.TxInput{
		Address:  pubKey.Address(),
		PubKey:   pubKey,
		Coins:    coins,
		Sequence: sequence,
	}
	if sequence > 1 {
		input.PubKey = nil
	}
	return input
}
