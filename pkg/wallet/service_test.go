package wallet

import (
	"fmt"
	"testing"
   "github.com/google/uuid"
	"github.com/ayub94/wallet/pkg/types"
	"reflect"

)
func TestService_Reject_success(t *testing.T) {
	// создаём сервис
	s := newTestServiceUser()

	_, payments, err := s.addAccountUser(defaultTestAccountUser)
	if err != nil {
		t.Error(err)
		return 
	} 
	//попробуем отменит платеж
	payment := payments[0]
	err = s.Reject(payment.ID)
	if err != nil {
		t.Errorf("Reject(): cannot reject payment, error = %v", err)
		return
	}
}
type testServiceUser struct { 
	*Service
}
func newTestServiceUser() *testServiceUser {
	return &testServiceUser{Service: &Service{}}
}
func (s *testServiceUser) addAccountWithBalance(phone types.Phone, balance types.Money)(*types.Account, error) {
	// sign in thre a user
	account, err := s.RegisterAccount(phone)
	if err != nil {
		return nil, fmt.Errorf("can not register account, error = %v", err)
	}
	// deposit balance
	err = s.Deposit(account.ID, balance)
	if err != nil {
		return nil, fmt.Errorf("can not deposit account, error = %v", err)
	}
	return account, nil
}
func TestService_FindPaymentByID_success(t *testing.T) {
    // создаём сервис
	s := newTestServiceUser()

	_, payments, err := s.addAccountUser(defaultTestAccountUser)
	if err != nil {
		t.Error(err)
		return
	}
	// попробуем найти платёж
	payment := payments[0]
	got, err := s.FindPaymentByID(payment.ID)
	if err != nil {
		t.Errorf("FindPaymentByID():  error = %v", err)
		return
	}
	// сравниваем платежи
	if !reflect.DeepEqual(payment, got) {
		t.Errorf("FindPaymentByID(): wrong payment returned = %v", err)
		return
	}
}
func TestService_FindPaymentByID_fail(t *testing.T) {
    // создаём сервис
	s := newTestServiceUser()
	_, _, err := s.addAccountUser(defaultTestAccountUser)
	if err != nil {
		t.Error(err)
		return
	}
	// попробуем найти не существуюший платёж
	_, err = s.FindPaymentByID(uuid.New().String())
	if err == nil {
		t.Errorf("FindPaymentByID():  must return error, returned nil")
		return
	}
	// сравниваем платежи
	if err != ErrPaymentNotFound {
		t.Errorf("FindPaymentByID(): cannot find payment = %v", err)
		return
	}
}
type testAccount struct {
	phone       types.Phone
	balance     types.Money
	payments    []struct {
		amount        types.Money
		category      types.PaymentCategory     
	}
}
func (s *testServiceUser) addAccountUser(data testAccount) (*types.Account, []*types.Payment, error) {
	// регистрируем там пользователья
	account, err := s.RegisterAccount(data.phone)
	if err != nil {
		return nil, nil, fmt.Errorf("cannot register account, error = %v", err)
	}
	// попольным его счет
	err = s.Deposit(account.ID, data.balance)
	if err != nil {
		return nil, nil, fmt.Errorf("cannot deposit account, error = %v", err)
	}

	// выполняем платежи
	// можем создат слайс сразу нужной длины, поскольку знаем размер
	payments := make([]*types.Payment, len(data.payments))
	for i, payment := range data.payments {
		// тогда здесь работаем просто через индекс, а не через append
		payments[i], err = s.Pay(account.ID, payment.amount, payment.category)
		if err != nil {
			return nil, nil, fmt.Errorf("cannot make paymennt, error = %v", err)
		} 
	}
	return account, payments, nil
}

var defaultTestAccountUser = testAccount {
	phone:        "+992934251222",
	balance:       10_000_00,
	payments:      []struct {
		amount        types.Money
		category      types.PaymentCategory
}{
	{amount: 1_000_00, category: "auto"},
},	
}
//func TestService_Repeat_success(t *testing.T){

//	s := newTestServiceUser()
//	account, err := s.RegisterAccount(data.phone)
//	if err != nil {
//		t.Errorf("RegisterAccount: cannot register account, error = %v", err)
//		return 
//	}
//}


