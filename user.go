package main

import (
	"database/sql"
	"log"
)

type user struct {
	id           int
	respectCount int
	connection   *sql.DB
}

func (u *user) incrementRespect() {
	u.respectCount++
	_, err := u.connection.Exec("UPDATE users SET respect_count = ? WHERE id = ?", u.respectCount, u.id)
	if err != nil {
		log.Fatalln(err)
	}
}

func newUser(id, respectCount int, db *sql.DB) *user {
	return &user{
		id:           id,
		respectCount: respectCount,
		connection:   db,
	}
}

func addUser(db *sql.DB, user *user) {
	stmt, _ := db.Prepare("INSERT INTO users (id,respect_count) VALUES(?,?)")
	stmt.Exec(user.id, user.respectCount)
	defer stmt.Close()
}

func getUser(db *sql.DB, id int) *user {
	u := newUser(id, 0, db)
	err := db.QueryRow("SELECT respect_count FROM users WHERE id = ?", id).Scan(&u.respectCount)
	if err != nil {
		addUser(db, u)
	}

	return u
}
