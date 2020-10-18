package wallet

import (
	"sync"
	"io"
	"strconv"
	"github.com/google/uuid"
	"github.com/ayub94/wallet/pkg/types"
	"errors"
	"log"
	"os"
	"fmt"
	"strings"
) 

type Service struct {
	nextAccountID int64
	accounts []*types.Account
	payments []*types.Payment
	favorites []*types.Favorite
}

var (
	ErrAccountNotFound = errors.New("Account not found")
	ErrRegisteredPhone = errors.New("Phone already registered")
	ErrMustBePossitive = errors.New("Amount must be greater than zero")
	ErrNotEnoughpBalance = errors.New("not enoughp balance")
	ErrPaymentNotFound = errors.New("payment not found")
	ErrFavoriteNotFound = errors.New("favorite not found")
	ErrFileNotFound = errors.New("file not found")
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
// Repeat, repeat payment meth
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

func (s *Service)FavoritePayment(paymentID string, name string)(*types.Favorite, error) {
	payment, err := s.FindPaymentByID(paymentID)
	if err != nil {
		return nil, err
	}

	favoriteID := uuid.New().String()
	favorite := &types.Favorite{
		ID:          favoriteID,
		AccountID:   payment.AccountID,
		Name:        name,
		Amount:      payment.Amount,
		Category:    payment.Category,
	}
	s.favorites = append(s.favorites, favorite)
	return favorite, nil
}
func (s *Service)PayFromFavorite(favoriteID string)(*types.Payment, error) {
	var favorite *types.Favorite
	for _, fav := range s.favorites {
		if fav.ID == favoriteID {
		        favorite = fav
	    }
	}
	if favorite == nil {
		return nil, ErrFavoriteNotFound
	}		
	payment, err := s.Pay(favorite.AccountID, favorite.Amount, favorite.Category)
	if err != nil {
		return nil, err
	}
	return payment, nil
}

func (s *Service)ExportToFile(path string) error {
	file, err := os.Create(path)    
	if err != nil {
		log.Print(err)
		return ErrFileNotFound
	}
	defer func(){
		err := file.Close()
		if err != nil {
		log.Print(err)
		} 
	}()
	var strg string
	for _, account := range s.accounts{
	   strg +=  fmt.Sprint(account.ID) + ";"+ fmt.Sprint(account.Phone) +";"+ fmt.Sprint(account.Balance) +"|"
	}   
		_, err = file.WriteString(strg)
		if err != nil {
			return err
		}
	    return nil
}
func (s *Service)ImportFromFile(path string) error {
file, err := os.Open(path)
 if err != nil {
	log.Print(err)
    return ErrFileNotFound
}
defer func(){
if cerr := file.Close(); cerr != nil {
    log.Print(cerr)
}
}()
content :=make([]byte, 0)
buf := make([]byte, 4)
for{
  read, err := file.Read(buf)
  if err == io.EOF {
    break
  }

  if err!=nil {
    log.Print(err)
    return ErrFileNotFound
  }
content = append(content, buf[:read]...)
}
data:=string(content)
 
accounts :=strings.Split(data, "|")
accounts = accounts[:len(accounts)-1]
for _, account := range accounts {
 value := strings.Split(account, ";")
 id,err := strconv.Atoi(value[0])
    if err!=nil {
    return err
    }
phone :=types.Phone(value[1])
balance, err := strconv.Atoi(value[2])
if err!=nil {
    return err
}
addAccount := &types.Account {
   ID: int64(id),
   Phone: phone,
   Balance: types.Money(balance),
}

s.accounts = append(s.accounts, addAccount)
log.Print(account)
}
return nil
}

func (s *Service)Export(dir string) error {
	if len(s.accounts) != 0{
	filedir1, err := os.Create(dir + "/accounts.dump") 
	if err != nil {
		log.Print(err)
		return ErrFileNotFound
	}
	defer func(){
		err := filedir1.Close()
		if err != nil {
		log.Print(err)
		} 
	}()
	var str1 string
	for _, account := range s.accounts{
	   str1 +=  fmt.Sprint(account.ID) + ";"+ fmt.Sprint(account.Phone) +";"+ fmt.Sprint(account.Balance) +"|"
	}   
		_, err = filedir1.WriteString(str1)
	
	}	
	if len(s.payments) != 0{
	filedir2, err := os.Create(dir + "/payments.dump") 
	if err != nil {
		log.Print(err)
		return ErrFileNotFound
	}
	defer func(){
		err := filedir2.Close()
		if err != nil {
		log.Print(err)
		} 
	}()
	var str2 string
	for _, payment := range s.payments{
		str2 +=  fmt.Sprint(payment.AccountID) + ";"+ fmt.Sprint(payment.ID) + ";"+ fmt.Sprint(payment.Amount) +";"+ fmt.Sprint(payment.Category)+";"+ fmt.Sprint(payment.Status) +"|"
		}   
			_, err = filedir2.WriteString(str2)
		
	}	
	if len(s.favorites) != 0{
	filedir3, err := os.Create(dir + "/favorits.dump") 
	if err != nil {
		log.Print(err)
		return ErrFileNotFound
	}
	defer func(){
		err := filedir3.Close()
		if err != nil {
		log.Print(err)
		} 
	}()
	var str3 string
	for _, favorite := range s.favorites{
	   str3 +=  fmt.Sprint(favorite.ID) + ";"+ fmt.Sprint(favorite.AccountID) +";"+ fmt.Sprint(favorite.Name)+";"+ fmt.Sprint(favorite.Amount)+";"+ fmt.Sprint(favorite.Category) +"|"
	}   
		_, err = filedir3.WriteString(str3)
			
	}
	return nil
}	

func (s *Service)Import(dir string) error{
	fileaccounts, err := os.Open(dir + "/accounts.dump")
	if err != nil {
		log.Print(err)
		//return ErrFileNotFound
		err = ErrFileNotFound
	}
	if err != ErrFileNotFound{

	defer func(){
		if cerr := fileaccounts.Close() ; cerr !=nil {
			log.Print(cerr)
		}
	}()
	actcontent := make([]byte,0)
	actbuf := make([]byte, 4)
	for {
		read, err := fileaccounts.Read(actbuf)
		if err == io.EOF{
			break
		}
		if err != nil {
			log.Print(err)
			return ErrFileNotFound
		}
		actcontent = append(actcontent, actbuf[:read]...)
	}
	actdata := string(actcontent)

	accounts := strings.Split(actdata, "|")
	accounts = accounts[:len(accounts)-1]

	for _, account := range accounts {
		value := strings.Split(account, ";")
		id, err := strconv.Atoi(value[0])
		if err != nil {
			return err
		}
		phone := types.Phone(value[1])
		balance, err := strconv.Atoi(value[2])
		if err != nil {
			return err
		}
		addAccount := &types.Account {
			ID: int64(id),
			Phone: phone,
			Balance: types.Money(balance),
	    }
		 
		s.accounts = append(s.accounts, addAccount)
		log.Print(account)
	}
   }
	//return nil
	
	filepayments, err := os.Open(dir + "/payments.dump")
	if err != nil {
		log.Print(err)
		//return ErrFileNotFound
		err = ErrFileNotFound
	}
	if err != ErrFileNotFound {

	defer func(){
		if cerr := filepayments.Close(); cerr !=nil {
			log.Print(cerr)
		}
	}()
	pmtcontent := make([]byte,0)
	pmtbuf := make([]byte, 4)
	for {
		read, err := filepayments.Read(pmtbuf)
		if err == io.EOF{
			break
		}
		if err != nil {
			log.Print(err)
			return ErrFileNotFound
		}
		pmtcontent = append(pmtcontent, pmtbuf[:read]...)
	}
	pmtdata := string(pmtcontent)

	payments := strings.Split(pmtdata, "|")
	payments = payments[:len(payments)-1]

	for _, payment := range payments {
		val := strings.Split(payment, ";")
		accountID, err := strconv.Atoi(val[0])
		if err != nil {
			return err
		}
		paymentID := string(val[1])
		paymentAmount, err := strconv.Atoi(val[2])
		if err != nil {
			return err
		}
		paymentCategory := types.PaymentCategory(val[3])
		paymentStatus := types.PaymentStatus(val[4])
		addPayment := &types.Payment {
			AccountID: int64(accountID),
			ID:  paymentID,
			Amount:  types.Money(paymentAmount),
			Category:  types.PaymentCategory(paymentCategory),
			Status:   types.PaymentStatus(paymentStatus),    
	    }
		 
		s.payments = append(s.payments, addPayment)
		log.Print(payment)
	}
	//return nil
}
	

	filefavorites, err := os.Open(dir + "/favorites.dump")
	if err != nil {
		log.Print(err)
		//return ErrFileNotFound
		err = ErrFileNotFound
	}

	if err != ErrFileNotFound{

	defer func(){
		if cerr := filefavorites.Close() ; cerr !=nil {
			log.Print(cerr)
		}
	}()
	fvtcontent := make([]byte,0)
	fvtbuf := make([]byte, 4)
	for {
		read, err := filefavorites.Read(fvtbuf)
		if err == io.EOF{
			break
		}
		if err != nil {
			log.Print(err)
			return ErrFileNotFound
		}
		fvtcontent = append(fvtcontent,fvtbuf[:read]...)
	}
	fvtdata := string(fvtcontent)

	favorites := strings.Split(fvtdata, "|")
	favorites = favorites[:len(favorites)-1]

	for _, favorite := range favorites {
		v := strings.Split(favorite, ";")
		favID := string(v[0])
		favactID, err := strconv.Atoi(v[1])
		if err != nil {
			return err
		}
		favName := string(v[2])
		favAmount, err := strconv.Atoi(v[3])
		if err != nil {
			return err
		}
		favCategory := types.PaymentCategory(v[4])

		addFavorite := &types.Favorite {
			ID:  favID,
			AccountID: int64(favactID),
			Name:   favName,
			Amount:  types.Money(favAmount),
			Category:  types.PaymentCategory(favCategory),    
	    }
		 
		s.favorites = append(s.favorites, addFavorite)
		log.Print(favorite)
	}
}		
	return nil
}
// ExportAccountHistory exports the history of payments of 
func (s *Service)ExportAccountHistory(accountID int64)([]types.Payment, error) {
	account, err := s.FindAccountByID(accountID)
	if err != nil {
		return nil, err
	}
	var payments []types.Payment

	for _, payment := range s.payments {
		if payment.AccountID == account.ID{
			data := types.Payment{
				ID:     payment.ID,
				AccountID:    payment.AccountID,
				Amount:   payment.Amount,
				Category:   payment.Category,
				Status:  payment.Status,
			}
			payments = append(payments, data)
		}
	}
	return payments, nil
}
// HistoryToFiles save datas recieved from above method
func (s *Service) HistoryToFiles(payments []types.Payment, dir string, records int) error {

	if len(payments) > 0 {
		if len(payments) <= records {
			file, _ := os.OpenFile(dir+"/payments.dump", os.O_CREATE|os.O_WRONLY|os.O_TRUNC, 0666)
			defer func(){
				if cerr := file.Close(); cerr != nil {
					log.Print(cerr)
				}
			}()
			var str string
			for _, val := range payments {
				str += fmt.Sprint(val.ID) + ";" + fmt.Sprint(val.AccountID) + ";" + fmt.Sprint(val.Amount) + ";" + fmt.Sprint(val.Category) + ";" + fmt.Sprint(val.Status) + "\n"
			}
			file.WriteString(str)
		} else {

			var str string
			k := 0
			j := 1
			var file *os.File
			for _, val := range payments {
				if k == 0 {
					file, _ = os.OpenFile(dir+"/payments"+fmt.Sprint(j)+".dump", os.O_CREATE|os.O_WRONLY|os.O_TRUNC, 0666)
				}
				k++
				str = fmt.Sprint(val.ID) + ";" + fmt.Sprint(val.AccountID) + ";" + fmt.Sprint(val.Amount) + ";" + fmt.Sprint(val.Category) + ";" + fmt.Sprint(val.Status) + "\n"
				_, err := file.WriteString(str)
				if err!=nil {
					log.Print(err)
				}
				if k == records {
					str = ""
					j++
					k = 0
					file.Close()
				}
			}
		}
	}

	return nil
}


//SumPayments сумирует платежи
func (s *Service)SumPayments(goroutines int) types.Money {
	//goroutines = 2
	wg := sync.WaitGroup{}
	wg.Add(goroutines) // сколько горутин ждём
	mu := sync.Mutex{} //мютекс сразу пишут над теми данными, доступ к которым нужно закрытъ
	//var pmt *types.Payment
	var sumPeyments types.Money

	go func(){
		defer wg.Done() // сообщает что завершено
		for _, peyment := range s.payments {
			pmt := peyment
			sumPeyments += pmt.Amount	
		}
		mu.Lock()
		defer mu.Unlock()
	}()
	go func(){
		defer wg.Done() // сообщает что завершено
		for _, peyment := range s.payments {
			pmt := peyment
			sumPeyments += pmt.Amount	
		}
		mu.Lock()
		defer mu.Unlock()
	}()
	wg.Wait()
	return sumPeyments

}
