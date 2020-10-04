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

func TestService_FindPaymentByID(t *testing.T) {
	type args struct {
		paymentID string
	}
	tests := []struct {
		name    string
		s       *Service
		args    args
		want    *types.Payment
		wantErr bool
	}{
		// TODO: Add test cases.
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := tt.s.FindPaymentByID(tt.args.paymentID)
			if (err != nil) != tt.wantErr {
				t.Errorf("Service.FindPaymentByID() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if !reflect.DeepEqual(got, tt.want) {
				t.Errorf("Service.FindPaymentByID() = %v, want %v", got, tt.want)
			}
		})
	}
}
