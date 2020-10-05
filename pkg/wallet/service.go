package wallet

import (
	"github.com/google/uuid"
	"github.com/ayub94/wallet/pkg/types"
	"errors"
) 

type Service struct {
	nextAccountID int64
	accounts []*types.Account
	payments []*types.Payment
}

var (
	ErrAccountNotFound = errors.New("Account not found")
	ErrRegisteredPhone = errors.New("Phone already registered")
	ErrMustBePossitive = errors.New("Amount must be greater than zero")
	ErrNotEnoughpBalance = errors.New("not enoughp balance")
	ErrPaymentNotFound = errors.New("payment not found")
)

func (s *Service) RegisterAccount(phone types.Phone) (*types.Account, error) {
	for _, account := range s.accounts {
		if account.Phone == phone{
			return nil, ErrRegisteredPhone
		}	
	}
	s.nextAccountID++
	account :=  &types.Account{
		ID:   s.nextAccountID,
		Phone:   phone,
		Balance: 0,
	}
	s.accounts =append(s.accounts,account)
	return account, nil
}

func  (s *Service) Deposit(accountID int64, amount types.Money) error {
	if amount <= 0 {
		return ErrMustBePossitive
	}
	var account *types.Account
	
	for _, acnt := range s.accounts {
		if acnt.ID == accountID {
			account = acnt
			break
		}  
	}	
	if account == nil{
		return ErrAccountNotFound
	}
	
	account.Balance += amount	
	return nil
}

func (s *Service)Pay(acntID int64, amount types.Money, category types.PaymentCategory)(*types.Payment, error) {
	if amount <= 0 {
		return nil, ErrMustBePossitive
	}
	var account *types.Account
	
	for _, acnt := range s.accounts {
		if acnt.ID  == acntID {
			account = acnt
			break
		}  
	}
	if account == nil {
		return nil, ErrAccountNotFound
	}	
	if account.Balance <= 0 {
		return nil, ErrNotEnoughpBalance
	}
	account.Balance-=amount
	paymentID := uuid.New().String() 
	payment := &types.Payment{
		ID: paymentID,
		AccountID: acntID,
		Amount: amount,
		Category: category,
		Status: types.PaymentStatusInProgress,
	}
	s.payments = append(s.payments, payment)
	return payment, nil
	
}

func (s *Service) FindAccountByID(accountID int64)(*types.Account, error) {

	var account *types.Account

	for _, acc:= range s.accounts {
		if acc.ID == accountID {
			account = acc
			return account, nil
		}
	}
	return nil, ErrAccountNotFound
}

func (s *Service)Reject(paymentID string)  error{

	var findPayment *types.Payment
	var findAccount *types.Account

	for _, pmnt := range s.payments {
		if pmnt.ID == paymentID {
			findPayment = pmnt
			break
		}
	}	
	if findPayment==nil{
		return ErrPaymentNotFound
	}	
    for _, acnt := range s.accounts {
	    if acnt.ID == findPayment.AccountID{
			findAccount = acnt
			break	
		}
	}
	if findAccount==nil{
		return ErrAccountNotFound
	}
	findPayment.Status = types.PaymentStatusFail
	findAccount.Balance += findPayment.Amount
	return nil
}

func (s *Service) FindPaymentByID(paymentID string)(*types.Payment, error) {
	for _, payment := range s.payments {
		if payment.ID == paymentID {
			return payment, nil
		}
	}
	return nil, ErrPaymentNotFound
}
// Repeat zuri baday
func (s *Service) Repeat(paymentID string)(*types.Payment, error)  {
	payment, err := s.FindPaymentByID(paymentID)
	if err != nil {
		return nil, err
	}

	payment1, err := s.Pay(payment.AccountID, payment.Amount, payment.Category)
	if err != nil {
		return nil, err
	}
	return payment1, nil
}
