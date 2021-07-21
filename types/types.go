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
	Id         int
	Type       string
	Name       string
	NationalId string
	Address    string
	Notes      string
	CreatedAt  string
}

type Note struct {
	Id          int
	Description string
	Urgency     string
	Attended    bool
	CreatedAt   string
	AttendedAt  string
}

type Bill struct {
	Id        int
	Url       string
	Date      string
	Company   string
	Charged   bool
	CreatedAt string
}

type Trip struct {
	Id     int
	Date   string
	Origin struct {
		Id         int
		Name       string
		NationalId string
		Address    string
	}
	Destination struct {
		Id         int
		Name       string
		NationalId string
		Address    string
	}
	Cargo     string
	Driver    string
	Truck     string
	Bill      string
	Voucher   string
	Completed bool
	Notes     string
}

type IdResponse struct {
	Id int
}

var MaxTransactionAmount = 1e14
var MaxBalanceAmount = 1e19