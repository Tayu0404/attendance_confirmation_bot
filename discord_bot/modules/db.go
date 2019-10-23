package module

import (
	"fmt"
	"github.com/jmoiron/sqlx"
	_"github.com/go-sql-driver/mysql"
)

type UserData struct {
	User_Name  int        `db:"user_name"`
	DATE       string     `db:"date"`
	Reason     string     `db:"reason"`
}

type Users struct {
	ID         int        `db:"id"`
	User_Name  string        `db:"user_name"`
}

func SelectDB (db *sqlx.DB) ([]UserData) {
	u := []UserData{}
	err := db.Select(&u,`
		SELECT
			users.user_name,
			data.date,
			data.reason
		FROM
			data
		INNER JOIN users
		ON data.user_id = users.id
	`)
	if err != nil {
		fmt.Println("SelectDB:", err)
	}
	return u
}

func SelectUserDB (db *sqlx.DB, uid string) ([]UserData) {
	u := []UserData{}
	err := db.Select(&u, fmt.Sprint(`
		SELECT
			users.user_name,
			data.date,
			data.reason
		FROM
			data
		INNER JOIN users
		ON data.user_id = users.id
		WHERE users.user_name='`, uid, `'
	`))
	if err != nil {
		fmt.Println("SelectUserDB : ",err)
	}
	return u
}

func AddToDB(db *sqlx.DB, user string, date string, reason string) (err error){
	id := UserCheckDB(db, user)
	fmt.Println(id)
	
	_, err = db.Query(fmt.Sprint(`
	INSERT INTO data (user_id, date, reason)
	VALUES ('`, id, `', '`, date, `', '`, reason, `')`))

	if err != nil {
		fmt.Println("AddToDB:", err)
	}
	return
}

func UserCheckDB (db *sqlx.DB, user string) int {
	u := []Users{}
	
	err := db.Select(&u,`
		SELECT *
		FROM users
	`)
	if err != nil {
		fmt.Println("UserCheckDB:", err)
	}

	id := contains(db, u, user)
	return id
}

func contains (db *sqlx.DB, arr []Users, user string) int {
	//user check
	for _, v := range arr {
		fmt.Println(v, v.ID, v.User_Name)
		if v.User_Name == user {
			return v.ID
		}
	}

	//create user
	_, err := db.Query(fmt.Sprint(`
		INSERT INTO users (user_name)
		VALUES ('`, user, `')`))
	if err != nil {
		fmt.Println("contains:", err)
	}
	var id Users
	err = db.Select(&id, fmt.Sprint(`
		SELECT id
		FROM users
		WHERE user_name='`, user, `'`))
	return id.ID
}

func CheckDate (db *sqlx.DB, userID string, date string) bool {
	u := []string{}

	fd := fmt.Sprintf("%s-%s-%s", date[:4], date[4:6], date[6:8])
	fmt.Println(fd)

	err := db.Select(&u, fmt.Sprint(`
		SELECT
			data.date
		FROM
			data
		INNER JOIN users ON data.user_id = users.id
		WHERE users.user_name='`, userID, `'
		AND data.date='`, fd, `'`))
	if err != nil {
		fmt.Println("CheckDate:", err)
		return false
	}
	if len(u) == 0 {
		return false
	}
	return true
}
