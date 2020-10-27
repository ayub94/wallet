package types

import (
	//"github.com/ayub94/wallet/pkg/types"
)

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

// Favorite - Изобранное
type Favorite struct{
    ID                 string
    AccountID          int64
    Name               string
    Amount             Money
    Category           PaymentCategory
}


