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
	create table if not exists key (bytes blob not null);
	create table if not exists messages
	(id integer not null primary key, sender_user_id text, receiver_user_id text, sent_datetime datetime, text text,
		FOREIGN KEY(receiver_user_id) REFERENCES contacts(user_id));
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
	_, err := p.db.Exec("insert into contacts values (?)", contact.UserId)
	if err != nil {
		panic(err)
	}
}

func (p *SQLiteStorageProvider) SaveKeyBytes(keyBytes []byte) {
	_, err := p.db.Exec("delete from key; insert into key values (?)", keyBytes)
	if err != nil {
		panic(err)
	}
}

func (p *SQLiteStorageProvider) GetKeyBytes() []byte {
	row := p.db.QueryRow("select bytes from key limit 1")
	if row.Err() != nil {
		panic(row.Err())
	}

	var keyBytes []byte
	err := row.Scan(&keyBytes)
	if err == sql.ErrNoRows {
		return nil
	}
	if err != nil {
		panic(err)
	}
	return keyBytes
}

func (p *SQLiteStorageProvider) SaveMessage(message entites.Message) {
	_, err := p.db.Exec(`insert into messages (sender_user_id, receiver_user_id, sent_datetime, text)
	values (?, ?, ?, ?)`,
		message.SenderContact.UserId, message.ReceiverContact.UserId, message.SentDatetime, message.Text)
	if err != nil {
		panic(err)
	}
}

func (p *SQLiteStorageProvider) GetMessages(userId string) []entites.Message {
	rows, err := p.db.Query(`select sender_user_id, receiver_user_id, sent_datetime, text from messages
	where sender_user_id = ? OR receiver_user_id = ?`, userId, userId)
	if err != nil {
		panic(err)
	}
	defer rows.Close()

	messages := make([]entites.Message, 0)
	for rows.Next() {
		var message entites.Message
		err = rows.Scan(
			&message.SenderContact.UserId,
			&message.ReceiverContact.UserId,
			&message.SentDatetime,
			&message.Text)
		if err != nil {
			panic(err)
		}
		messages = append(messages, message)
	}
	err = rows.Err()
	if err != nil {
		panic(err)
	}
	return messages
}
