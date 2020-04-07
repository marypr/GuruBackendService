package src

import (
	"time"
)

// BetType - transaction type
type BetType string

const (
	// TypeWin - win bet
	TypeWin BetType = "Win"
	// TypeBet - making bet
	TypeBet BetType = "Bet"
)

//User - user model
type User struct {
	ID      uint64  `json:"id"`
	Balance float64 `json:"balance"`
	Token   string  `json:"token"`
}

// Statistics - statistics model
type Statistics struct {
	DepositCount uint    `json:"depositCount"`
	DepositSum   float64 `json:"depositSum"`
	BetCount     uint    `json:"betCount"`
	BetSum       float64 `json:"betSum"`
	WinCount     uint    `json:"winCount"`
	WinSum       float64 `json:"winSum"`
}

// Bet - bet model
type Bet struct {
	ID     uint64
	Type   BetType
	UserID uint64
	Amount float64
	Token  string
}

// Deposit - deposit model
type Deposit struct {
	ID            int       `json:"id"`
	UserID        uint64    `json:"userId"`
	Amount        float64   `json:"amount"`
	Token         string    `json:"token"`
	BalanceBefore float64   `json:"balanceBefore"`
	BalanceAfter  float64   `json:"balanceAfter"`
	Time          time.Time `json:"time"`
}

// Transaction - transaction model
type Transaction struct {
	ID            uint64    `json:"id"`
	UserID        uint64    `json:"userId"`
	BetID         uint      `json:"betID"`
	Amount        float64   `json:"amount"`
	BalanceBefore float64   `json:"balanceBefore"`
	BalanceAfter  float64   `json:"balanceAfter"`
	Time          time.Time `json:"time"`
}

// GetUserData - data for GetUser
type GetUserData struct {
	ID    uint64 `json:"id"`
	Token string `json:"token"`
}

// ReturnedGetUserData- data for return GetUser data
type ReturnedGetUserData struct {
	ID uint64 `json:"id"`
	Statistics
}

// AddDepositData - received data for AddDeposit
type AddDepositData struct {
	ID     int     `json:"id"`
	UserID uint64  `json:"userId"`
	Amount float64 `json:"amount"`
	Token  string  `json:"token"`
}

// ReturnedAddDepositData - returned data for AddDeposit
type ReturnedAddDepositData struct {
	Error   error   `json:"error"`
	Balance float64 `json:"balance"`
}

// TransactionData - received data for AddDeposit
type TransactionData struct {
	ID     uint64  `json:"id"`
	UserID uint64  `json:"userId"`
	Type   BetType `json:"betType"`
	Amount float64 `json:"amount"`
	Token  string  `json:"token"`
}
