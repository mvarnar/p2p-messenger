package storage

import (
	entites "p2p-messenger/src/domain/entities"
	_ "github.com/mattn/go-sqlite3"
	"database/sql"
	"fmt"
)

type SQLiteStorageProvider struct {
	db *sql.DB
}

func NewSQLiteStorageProvider() SQLiteStorageProvider {
	db, err := sql.Open("sqlite3", "./database.db")
	if err != nil {
		panic(err)
	}

	sqlStmt := `
	create table if not exists contacts (user_id text not null primary key);
	`
	_, err = db.Exec(sqlStmt)
	if err != nil {
		panic(fmt.Sprintf("%q: %s\n", err, sqlStmt))
	}

	return SQLiteStorageProvider{db: db}
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
