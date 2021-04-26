package responsetypes

type WalletType struct {
	PrivateKey string `json:"PrivateKey"`
	PublicKey  string `json:"PublicKey"`
	Address    string `json:"Address"`
}

type TransactionParams struct {
	SenderPrivateKey string `json:"SenderPrivateKey"`
	RecipientAddress string `json:"RecipientAddress"`
	SenderAddress    string `json:"SenderAddress"`
	Value            string `json:"Value"`
	SenderPublicKey  string `json:"SenderPublicKey"`
}

type TransactionType struct {
	RecipientAddress string `json:"RecipientAddress"`
	SenderAddress    string `json:"SenderAddress"`
	Value            uint64 `json:"Value"`
	SenderPublicKey  string `json:"SenderPublicKey"`
	Signature        string `json:"Signature"`
}

type GetRequestWalletAmountType struct {
	Address string `json:"Address"`
}

type GetResponseWalletAmountType struct {
	Amount uint64 `json:"Amount"`
}
