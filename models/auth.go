package models

import (
	"fmt"
	"time"
)

func UserCreate(email string) {

	userExists := UserExist(email) 
	if userExists == false {
		db := RetrieveDB()
		fmt.Println("User create for email", email)
		sqlStatement := `INSERT INTO users (email, username, created_at, is_first_time) VALUES ($1, $2, $3, $4)`
		fmt.Println("sql statement", sqlStatement)
		sqlPreparedStatement, err := db.Prepare(sqlStatement)
		fmt.Println("sql prepared steate", sqlPreparedStatement)
		if err != nil {
			fmt.Println("SQL Statement is not prepared.", err)
			// return false
		}

		createdAt := time.Now().UTC()

		insertRow, err := sqlPreparedStatement.Exec(email, email, createdAt, true)

		if err != nil {
			fmt.Println("SQL Insert doesn't insert", err)
			// return false
		}
		fmt.Println("insert row", insertRow)

	}	
}

func UserExist(email string) (bool) {
	db := RetrieveDB()

	checkUser := `SELECT EXISTS (SELECT 1 FROM users WHERE email=$1)`
	statement , err := db.Prepare(checkUser)
	
	if err != nil {
		fmt.Println("SQL Statement is not prepared.", err)
		return false
	}

	var userStatus bool
	err = statement.QueryRow(email).Scan(&userStatus)

	if err != nil {
		fmt.Println("SQL Check doesn't insert", err)
		return false
	}
	return userStatus
}

func UserInfo() {

}
