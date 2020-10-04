package main

import (
	"fmt"
	"github.com/ayub94/wallet/pkg/wallet"
)

func main() {
	svc := &wallet.Service{}
	account, err := svc.RegisterAccount("+992 555551204")
	if err != nil {
		fmt.Println(err)
		return
	}
	fmt.Println(account)
	
	err = svc.Deposit(account.ID, 100)
	if err != nil {
		switch err{
		case wallet.ErrAccountNotFound:
			fmt.Println("not found")
		case wallet.ErrRegisteredPhone:
			fmt.Println("already registered")	
		case wallet.ErrMustBePossitive:
			fmt.Println("nust be possitive")	
		}
	}	
	fmt.Println(account.Balance) //100

	account, err = svc.FindAccountByID(12345)
	if err != nil {
		fmt.Println(err)
	}
	fmt.Println(account)	

}

