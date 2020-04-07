package src

import (
	"errors"
	"sync"
)

type cache struct {
	mx            sync.RWMutex
	modifiedUsers map[uint64]*User
	users         map[uint64]*User
	statistics    map[uint64]*Statistics
	transactions  map[uint64]*Transaction
}

// transactionCache - cache
var transactionCache cache

const (
	// AddDeposit - deposit changes
	AddDeposit = "addDeposit"
	// MakeBet - bet made
	MakeBet = "makeBet"
	// Win - win
	Win = "userWin"
)

func init() {
	transactionCache = cache{
		users:         make(map[uint64]*User),
		statistics:    make(map[uint64]*Statistics),
		transactions:  make(map[uint64]*Transaction),
		modifiedUsers: make(map[uint64]*User),
	}
}

func (c *cache) AddUser(user *User) error {
	c.mx.Lock()
	defer c.mx.Unlock()
	if _, ok := c.users[user.ID]; ok {
		return errors.New("user already exist")
	}

	c.users[user.ID] = user
	return nil
}

func (c *cache) GetUser(id uint64) (*User, error) {
	c.mx.RLock()
	defer c.mx.RUnlock()
	if user, ok := c.users[id]; ok {
		return user, nil
	}
	return nil, errors.New("user doesn't exist")
}

func (c *cache) GetModifiedUsers() map[uint64]*User {
	c.mx.RLock()
	defer c.mx.RUnlock()
	return c.modifiedUsers
}

func (c *cache) AddModifiedUser(user *User) {
	c.mx.Lock()
	c.modifiedUsers[user.ID] = user
	c.mx.Unlock()
}

func (c *cache) AddTransaction(userID uint64, transaction *Transaction) error {
	c.mx.Lock()
	c.transactions[userID] = transaction
	c.mx.Unlock()
	return nil
}

func (c *cache) GetUserStatistics(id uint64) *Statistics {
	c.mx.Lock()
	if _, ok := c.statistics[id]; !ok {
		c.statistics[id] = &Statistics{}
	}
	c.mx.Unlock()
	return c.statistics[id]
}

func (c *cache) ZeroingModifiedUsers() {
	c.mx.Lock()
	for k := range c.modifiedUsers {
		delete(c.modifiedUsers, k)
	}
	c.mx.Unlock()
}

func (c *cache) UpdateStatistic(userID uint64, amount float64, changesType string) (*Statistics, error) {
	statistics := c.GetUserStatistics(userID)
	c.mx.Lock()
	defer c.mx.Unlock()

	switch changesType {
	case AddDeposit:
		statistics.DepositSum += amount
		statistics.DepositCount++
	case MakeBet:
		statistics.BetSum += amount
		statistics.BetCount++
	case Win:
		statistics.WinSum += amount
		statistics.WinCount++
	default:
		return &Statistics{}, errors.New("invalid operation type")
	}
	return statistics, nil
}
