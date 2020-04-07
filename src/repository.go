package src

import "database/sql"

//UserRepository is repo for User
type UserRepository interface {
	AddUser(u User) (lastID int, err error)
	GetUser(id int) (u User, err error)
}

type postgresUsersRepository struct {
	DB *sql.DB
}

//NewPostgresUsersRepo is a function to get New postgresUsersRepository which uses given connection
func NewPostgresUsersRepo(db *sql.DB) UserRepository {
	return &postgresUsersRepository{db}
}

//InsertUser is a function that inserts a user entity into a database
func (p *postgresUsersRepository) AddUser(u User) (lastID int, err error) {
	err = p.DB.QueryRow("INSERT INTO guru.users (id,balance,token) values($1,$2,$3) RETURNING id",
		u.ID, u.Balance, u.Token).Scan(&lastID)
	return
}

//GetUser is a function that get user data
func (p *postgresUsersRepository) GetUser(id int) (u User, err error) {
	rows := p.DB.QueryRow(`SELECT id,balance FROM guru.users where idl=$1`, id)
	err = rows.Scan(&u.ID, &u.Balance)
	if err != nil {
		return
	}
	return
}
