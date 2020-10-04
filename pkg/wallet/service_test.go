package wallet

import (
	"reflect"
	"testing"

	"github.com/ayub94/wallet/pkg/types"
)

func TestService_FindAccountByID(t *testing.T) {
	type args struct {
		accountID int64
	}
	tests := []struct {
		name    string
		s       *Service
		args    args
		want    *types.Account
		wantErr bool
	}{
		// TODO: Add test cases.
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := tt.s.FindAccountByID(tt.args.accountID)
			if (err != nil) != tt.wantErr {
				t.Errorf("Service.FindAccountByID() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if !reflect.DeepEqual(got, tt.want) {
				t.Errorf("Service.FindAccountByID() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestService_Reject(t *testing.T) {
	type args struct {
		paymentID string
	}
	tests := []struct {
		name    string
		s       *Service
		args    args
		wantErr bool
	}{
		// TODO: Add test cases.
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if err := tt.s.Reject(tt.args.paymentID); (err != nil) != tt.wantErr {
				t.Errorf("Service.Reject() error = %v, wantErr %v", err, tt.wantErr)
			}
		})
	}
}

func TestService_Reject_success(t *testing.T){
	// создаём сервис
    s := &Service{}
	
	//регистрируем там пользователья
	phone := types.Phone("+992555551204")
	account, err := s.RegisterAccount(phone)
	if err != nil {
		t.Errorf("Reject(): cannot register account, error = %v",err)
		return
	}
	
	//паполняем ее счет
	err = s.Deposit(account.ID, 1_000_000)
	if err != nil {
		t.Errorf("Reject(): cannot deposit account, error = %v",err)
		return
	}

	//осушествляем платкж на его счет
	payment, err := s.Pay(account.ID, 100_000, "auto")
	if err != nil {
		t.Errorf("Reject(): cannot create payment, error = %v",err)
		return
	}

	//попробуем отменит платеж
	err = s.Reject(payment.ID)
	if err != nil {
		t.Errorf("Reject(): error = %v",err)
		return
	}
}
