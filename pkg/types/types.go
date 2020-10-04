package types

// Money представляет ссобой 
type Money int64

// PaymentCategory in string
type PaymentCategory string

// PaymentStatus is string
type PaymentStatus string


// Предопределенные статусы платежей
const (
    PaymentStatusOk    PaymentStatus = "Ok"
    PaymentStatusFail    PaymentStatus = "FAIL"
    PaymentStatusInProgress   PaymentStatus = "INPROGRESS"
)

// Payment provide struct
type Payment struct{
    AccountID          int64
    ID                 string
    Amount             Money
    Category           PaymentCategory
    Status             PaymentStatus
}

// Phone is string phone
type Phone string

// Account struct have struct
type Account struct {
	 ID int64
	 Phone Phone
	 Balance Money
}
