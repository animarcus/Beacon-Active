package Model

import (
	"encoding/binary"
	"fmt"
)

type UserId uint32
type UserName string

type InvalidUserNameError struct{
	UserName UserName
}

func (err *InvalidUserNameError) Error() string {
	return fmt.Sprintf("Error: username %s already taken", err.UserName)
}

type User struct {
	Id UserId `json:"user_id"`
	UserName UserName `json:"username"`
}

func (user *User) AddUser() (err error) {
	//users := map[UserId]*User{}
	//users, err = GetAllUsers()
	//if err!= nil {
	//	err = fmt.Errorf("Error adding user with username %s", user.UserName)
	//	return
	//}
	//for _, u := range users {
	//	if user.UserName == u.UserName {
	//		err = &InvalidUserNameError{UserName: user.UserName}
	//		return
	//	}
	//}

	statement := "insert into beaconactivedb.public.users (username) values ($1) returning id"
	stmt, err := Db.Prepare(statement)
	if err != nil {
		return
	}
	defer stmt.Close()
	err = stmt.QueryRow(user.UserName).Scan(&user.Id)
	return
}

func GetAllUsers() (users map[UserId]*User, err error) {
	users = map[UserId]*User{}
	rows, err := Db.Query("select id, username from beaconactivedb.public.users limit $1", 10)
	if err != nil {
		return
	}
	for rows.Next() {
		user := User{}
		err = rows.Scan(&user.Id, &user.UserName)
		if err != nil {
			return
		}
		users[user.Id] = &user
	}
	err = rows.Close()
	return
}

func GetUser(userId UserId) (user *User, err error) {
	users := map[UserId]*User{}
	users, err = GetAllUsers()
	if err != nil {
		return
	}
	user = users[userId]
	return
}

func (user *User) Bytes() []byte {
	var buf []byte
	binary.PutVarint(buf, int64(user.Id))
	return append(buf, []byte(user.UserName)...)
}