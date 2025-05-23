package users

import "errors"

type User struct {
	Id       int
	Username string
	Password string
}

type UsersInMemoryDB struct {
	Users []User
}

func (db *UsersInMemoryDB) Create(user User) {
	db.Users = append(db.Users, user)
}

func (db *UsersInMemoryDB) Update(id int, updUser User) error {
	for i, user := range db.Users {
		if user.Id == id {
			db.Users[i] = updUser
			return nil
		}
	}
	return errors.New("User not found")
}

func (db *UsersInMemoryDB) GetById(id int) (User, error) {
	for _, user := range db.Users {
		if user.Id == id {
			return user, nil
		}
	}
	return User{}, errors.New("User not found")
}
