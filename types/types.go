package types

type TransactionWithBalance struct {
	Id          int
	Type        string
	Amount      float32
	Description string
	Balance     float32
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
	Amount      float32
	Description string
	Actor       struct {
		Id   int
		Name string
	}
	CreatedAt   string
}

type PartialTransaction struct {
	Description string
}

type PartialPendingTransaction struct {
	Type        string
	Amount      float32
	Description string
	Actor       int
}

type Actor struct {
	Id          int
	Name        string
	Description string
	IsCompany   bool
	CreatedAt   string
}

type PartialActor struct {
	Name        string
	Description string
}

type IdResponse struct {
	Id int
}
