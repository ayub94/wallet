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
	account, err := s.RegisterAccount("+992555551204")
	if err != nil {
		return nil, fmt.Errorf("account alrady reagistered, error = %v", err)
	}
	// deposit balance
	err = s.Deposit(account.ID, 1000)
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
func TestService_Repeat_success(t *testing.T){

	s := newTestServiceUser()
	acc, err := s.RegisterAccount("+992934251221")
	if err != nil {
		t.Errorf("RegisterAccountUser: cannot register account, error = %v", err)
		return 
	}
	err = s.Deposit(acc.ID, 100)
	if err != nil {
		t.Errorf("can not deposit account, error = %v", err)
		return 
	}
	payment, err := s.Pay(acc.ID, 10, "ice-cream")
	if err != nil {
		t.Errorf("can not pay, error = %v",err)
		return 
	}
	payment1, err := s.FindPaymentByID(payment.ID)

	if err != nil{
		t.Errorf("method FindPaymentByID returned not nil error, payment => %v", payment)
	}
	payment1, err = s.Repeat(payment.ID)
	if err != nil {
		t.Errorf("can not repeat payment, error = %v",err)
		return
	}
	if payment.Amount != payment1.Amount || payment.Category != payment1.Category {
		t.Error("wrong result")
	}
}

func TestService_FindAccountByID_success_user(t *testing.T) {
	s := newTestServiceUser()
	
    s.RegisterAccount("+992934251220")
	account, err := s.FindAccountByID(1)

	if err != nil{
		t.Errorf("method FindPaymentByID returned not nil error, payment => %v", account)
		return
	}
}

func TestService_FindAccountByID_notFound_user(t *testing.T) {
	s := newTestServiceUser()
	
    s.RegisterAccount("+992934251220")
	account, err := s.FindAccountByID(2)

	if err == nil{
		t.Errorf("method FindPaymentByID returned nil error, payment => %v", account)
		return
	}
}


func TestService_FavoritePayment_success(t *testing.T){
	s := newTestServiceUser()
	
	account, err := s.RegisterAccount("+992934251220")

	if err != nil{
		t.Errorf("method RegisterAccount returned not nil error, account => %v", account)
	}

	err = s.Deposit(account.ID, 100_00)
	if err != nil{
		t.Errorf("method Deposit returned not nil error, error => %v", err)
	}


	payment, err := s.Pay(account.ID, 10_00,"Cafe")

	if err != nil{
		t.Errorf("method Pay returned not nil error, account => %v", account)
	}



	favorite, err := s.FavoritePayment(payment.ID, "My Favorite")

	if err != nil{
		t.Errorf("method FavoritePayment returned not nil error, favorite => %v", favorite)
	}

	paymentFavorite, err := s.PayFromFavorite(favorite.ID)
	if err != nil{
		t.Errorf("method PayFromFavorite returned not nil error, paymentFavorite => %v", paymentFavorite)
	}
}	

func TestService_Export_success_user(t *testing.T) {
	var svc Service

	svc.RegisterAccount("+992000000001")
	svc.RegisterAccount("+992000000002")
	svc.RegisterAccount("+992000000003")

	err := svc.ExportToFile("export.txt")
	if err != nil {
		t.Errorf("method ExportToFile returned not nil error, err => %v", err)
	}

}

func TestService_Import_success_user(t *testing.T) {
	var svc Service


	err := svc.ImportFromFile("export.txt")
	
	if err != nil {
		t.Errorf("method ExportToFile returned not nil error, err => %v", err)
	}
}	

  
  func TestService_Import_success(t *testing.T) {
	var svc Service
	err := svc.ImportFromFile("export.txt")
	if err != nil {
	  t.Errorf("method ExportToFile returned not nil error, err => %v", err)
	}
  
  }

  func TestService_ExportImport_success_user(t *testing.T) {
	var svc Service

	svc.RegisterAccount("+992000000001")
	svc.RegisterAccount("+992000000002")
	svc.RegisterAccount("+992000000003")
	svc.RegisterAccount("+992000000004")
	
	err := svc.Export(".")
	if err != nil {
		t.Errorf("method ExportToFile returned not nil error, err => %v", err)
	}

	err = svc.Import(".")
	
	if err != nil {
		t.Errorf("method ImportToFile returned not nil error, err => %v", err)
	}

}
  

  func BenchmarkSumPayments_user(b *testing.B) {
	var svc Service  
	want:= types.Money(0)
	for i:=0 ; i < b.N ; i++ {
		result := svc.SumPayments(2)
		if result != types.Money(want) {
			b.Fatalf("invalid result, got %v, want %v", result, want)

		}
	}
}

