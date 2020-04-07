package src

import (
	"log"
	"net/http"
	"sync"
	"time"

	"github.com/gin-gonic/gin"
)

const tokenTest = "testtask"

// UserService is a struct that is is used to set repository for usersRepo (or its mocks)
type UserService struct {
	UsersRepo UserRepository
}

//NewUserService is a func to get new UserService with user's defined repository
func NewUserService(repository UserRepository) *UserService {
	return &UserService{
		UsersRepo: repository,
	}
}

func (s *UserService) AddUser(c *gin.Context) {
	var user User
	if err := c.BindJSON(&user); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"status": "bad request"})
		return
	}

	if !validateToken(user.Token) {
		log.Println("Error while validating token")
		c.JSON(http.StatusBadRequest, gin.H{"error": "bad token"})
		return
	}

	//save new user to cache
	if err := transactionCache.AddUser(&user); err != nil {
		log.Println("Error while saving user to cache")
		c.JSON(http.StatusInternalServerError, gin.H{"error": "can't save user"})
		return
	}

	//save new user in db
	if _, err := s.UsersRepo.AddUser(user); err != nil {
		log.Println("Error in AddUserHandler while adding user in db: ")
		log.Println(err)
		c.JSON(http.StatusBadRequest, gin.H{"error": "bad request"})
		return
	}
	c.JSON(http.StatusOK, gin.H{"error": ""})
}

func (s *UserService) updateUsers() {
	for _, user := range transactionCache.GetModifiedUsers() {
		log.Println("Saving user to cache", user.ID)
		if _, err := s.UsersRepo.AddUser(*user); err != nil {
			log.Println("Error in AddUser while adding user in db: ")
			log.Println(err)
			return
		}
	}
	transactionCache.ZeroingModifiedUsers()
}

func (s *UserService) GetUser(c *gin.Context) {
	var userData GetUserData
	if err := c.BindJSON(&userData); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"status": "bad request"})
		return
	}

	if !validateToken(userData.Token) {
		log.Println("Error while validating token")
		c.JSON(http.StatusBadRequest, gin.H{"error": "bad token"})
		return
	}

	user, err := transactionCache.GetUser(userData.ID)
	if err != nil {
		log.Println("Error while getting user from cache", err)
		c.JSON(http.StatusBadRequest, gin.H{"error": "can't get user"})
		return
	}
	userStatistics := transactionCache.GetUserStatistics(userData.ID)

	c.JSON(http.StatusOK, ReturnedGetUserData{
		ID:         user.ID,
		Statistics: *userStatistics,
	})
	return
}

func (s *UserService) AddDeposit(c *gin.Context) {
	var mu sync.Mutex
	var depositData AddDepositData
	if err := c.BindJSON(&depositData); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"status": "bad request"})
		return
	}

	if !validateToken(depositData.Token) {
		log.Println("Error while validating token")
		c.JSON(http.StatusBadRequest, gin.H{"error": "bad token"})
		return
	}

	user, err := transactionCache.GetUser(depositData.UserID)
	if err != nil {
		log.Println("Error while getting user from cache", err)
		c.JSON(http.StatusBadRequest, gin.H{"error": "can't get user"})
		return
	}

	var deposit Deposit
	deposit.Time = time.Now()
	deposit.BalanceBefore = user.Balance
	deposit.BalanceAfter = user.Balance + depositData.Amount

	mu.Lock()
	defer mu.Unlock()
	user.Balance += depositData.Amount

	transactionCache.AddModifiedUser(user)
	if _, err := transactionCache.UpdateStatistic(user.ID, depositData.Amount, AddDeposit); err != nil {
		log.Println("Error while update user statistic after adding deposit ", err)
	}
	c.JSON(http.StatusOK, ReturnedAddDepositData{
		Error:   nil,
		Balance: user.Balance,
	})
	return
}

func (s *UserService) MakeTransaction(c *gin.Context) {
	var mu sync.Mutex
	var transactionData TransactionData
	if err := c.BindJSON(&transactionData); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"status": "bad request"})
		return
	}

	if !validateToken(transactionData.Token) {
		log.Println("Error while validating token")
		c.JSON(http.StatusBadRequest, gin.H{"error": "bad token"})
		return
	}

	user, err := transactionCache.GetUser(transactionData.UserID)
	if err != nil {
		log.Println("Error while getting user from cache", err)
		c.JSON(http.StatusBadRequest, gin.H{"error": "can't get user"})
		return
	}

	var statistic *Statistics
	var transaction *Transaction

	switch transactionData.Type {
	case TypeWin:
		mu.Lock()
		user.Balance += transactionData.Amount
		if statistic, err = transactionCache.UpdateStatistic(user.ID, transactionData.Amount, Win); err != nil {
			log.Println("Error while update user stats", err)
		}
		mu.Unlock()
		transaction = &Transaction{
			ID:            transactionData.ID,
			UserID:        transactionData.UserID,
			BetID:         statistic.BetCount,
			Amount:        transactionData.Amount,
			BalanceBefore: user.Balance,
			BalanceAfter:  statistic.WinSum,
			Time:          time.Now(),
		}
	case TypeBet:
		if !validateBalance(user.Balance, transactionData.Amount) {
			log.Println("Error while making transaction - balance is lower then transaction amount")
			c.JSON(http.StatusBadRequest, gin.H{"error": "so sorry, but u haven't enough money to make this transaction"})
			return
		}
		mu.Lock()
		user.Balance -= transactionData.Amount
		if statistic, err = transactionCache.UpdateStatistic(user.ID, transactionData.Amount, MakeBet); err != nil {
			log.Println("Error while update user stats", err)
		}
		mu.Unlock()
		transaction = &Transaction{
			ID:            transactionData.ID,
			UserID:        transactionData.UserID,
			BetID:         statistic.BetCount,
			Amount:        transactionData.Amount,
			BalanceBefore: user.Balance,
			BalanceAfter:  statistic.BetSum,
			Time:          time.Now(),
		}
	}

	if err = transactionCache.AddTransaction(user.ID, transaction); err != nil {
		log.Println("Error while making transaction - can't save transaction ", err)
		c.JSON(http.StatusInternalServerError, gin.H{"error": "something went wrong, try a little bit later"})
	}
	c.JSON(http.StatusOK, gin.H{"error": "", "balance": user.Balance})
}

func validateToken(token string) bool {
	return token != tokenTest
}

func validateBalance(balance, amount float64) bool {
	return balance > amount
}
