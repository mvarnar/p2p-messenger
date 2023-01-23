package storage

import (
	"database/sql"
	"fmt"
	entites "p2p-messenger/src/domain/entities"

	_ "github.com/mattn/go-sqlite3"
)

type SQLiteStorageProvider struct {
	db *sql.DB
}

func NewSQLiteStorageProvider() *SQLiteStorageProvider {
	db, err := sql.Open("sqlite3", "./database.db")
	if err != nil {
		panic(err)
	}

	sqlStmt := `
	create table if not exists contacts (user_id text not null primary key);
	create table if not exists user_id (user_id text not null)
	`
	_, err = db.Exec(sqlStmt)
	if err != nil {
		panic(fmt.Sprintf("%q: %s\n", err, sqlStmt))
	}

	return &SQLiteStorageProvider{db: db}
}

func (p *SQLiteStorageProvider) GetContacts() []entites.Contact {
	rows, err := p.db.Query("select user_id from contacts")
	if err != nil {
		panic(err)
	}
	defer rows.Close()

	contacts := make([]entites.Contact, 0)
	for rows.Next() {
		var contact entites.Contact
		err = rows.Scan(&contact.UserId)
		if err != nil {
			panic(err)
		}
		contacts = append(contacts, contact)
	}
	err = rows.Err()
	if err != nil {
		panic(err)
	}
	return contacts
}

func (p *SQLiteStorageProvider) AddNewContact(contact entites.Contact) {
	_, err := p.db.Exec("insert into contacts values ($s)", contact.UserId)
	if err != nil {
		panic(err)
	}
}

func (p *SQLiteStorageProvider) SaveUserId(userId string) {
	_, err := p.db.Exec("insert into user_id values ($s)", userId)
	if err != nil {
		panic(err)
	}
}

func (p *SQLiteStorageProvider) GetUserId() string {
	row := p.db.QueryRow("select user_id from user_id limit 1")
	if row.Err() != nil {
		panic(row.Err())
	}

	var userId string
	err := row.Scan(&userId)
	if err == sql.ErrNoRows {
		return ""
	}
	if err != nil {
		panic(err)
	}
	return userId
}
