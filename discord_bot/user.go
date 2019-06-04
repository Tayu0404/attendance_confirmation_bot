package modeles

import (
	"fmt"
	"log"

	"github.com/jmoiron/sqlx"
)

type User struct {
	ID		int		`db:"uuid"`
	Date	string	`db:"date"`
	Reason	string	`db:reason"`
}

type Userlist []User

func main() {
	var userlist Userlist
	
	db, err := sqlx.Open("mysql", "root@/attendance_rec_db")
	if err != nil {
		log.Fatal(err)
	}

	rows, err := db.Queryx("SELECT * FROM user")
	if err != nil {
		log.Fatal(err)
	}

	var user User
	for rows.Next() {
		err := rows.StructScan(&user)
		if err != nil {
			userlist = append(userlist, user)
		}
	}

	fmt.Println(userlist)
}

