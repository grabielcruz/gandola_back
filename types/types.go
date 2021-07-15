package types

type TransactionWithBalance struct {
	Id          int
	Type        string
	Currency    string
	Amount      float32
	Description string
	USDBalance  float32
	VESBalance  float32
	Actor       struct {
		Id   int
		Name string
	}
	Executed  string
	CreatedAt string
}

type PendingTransaction struct {
	Id          int
	Type        string
	Currency    string
	Amount      float32
	Description string
	Actor       struct {
		Id   int
		Name string
	}
	CreatedAt string
}

type Actor struct {
	Id          int
	Name        string
	Description string
	IsCompany   bool
	CreatedAt   string
}

type IdResponse struct {
	Id int
}

var MaxTransactionAmount = 1e14
var MaxBalanceAmount = 1e19
