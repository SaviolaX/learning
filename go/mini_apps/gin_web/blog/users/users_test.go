package users

import (
	"testing"

	"github.com/quii/go-graceful-shutdown/assert"
)

func TestCreate(t *testing.T) {
	testDB := UsersInMemoryDB{}

	newUser := User{
		Id:       1,
		Username: "testUser",
		Password: "testPassword",
	}

	testDB.Create(newUser)

	wantUsersLen := 1

	assert.Equal(t, len(testDB.Users), wantUsersLen)
}

func TestUpdate(t *testing.T) {
	t.Run("user updated", func(t *testing.T) {
		testDB := UsersInMemoryDB{}
		testDB.Users = append(testDB.Users, User{
			Id:       1,
			Username: "testUser",
			Password: "testPassword",
		})

		updatedUser := User{
			Id:       1,
			Username: "updatedUser",
			Password: "updatedPassword",
		}

		correctID := 1

		testDB.Update(correctID, updatedUser)

		wantUserUsername := "updatedUser"
		wantUserPassword := "updatedPassword"

		assert.Equal(t, testDB.Users[0].Username, wantUserUsername)
		assert.Equal(t, testDB.Users[0].Password, wantUserPassword)
	})
	t.Run("user not found", func(t *testing.T) {
		testDB := UsersInMemoryDB{}

		updatedUser := User{}

		incorrectID := 0000

		want := "User not found"

		err := testDB.Update(incorrectID, updatedUser)

		assert.Equal(t, err.Error(), want)

	})
}

func TestGetById(t *testing.T) {
	t.Run("got user by id", func(t *testing.T) {
		testDB := UsersInMemoryDB{}
		testDB.Users = append(testDB.Users, User{
			Id:       1,
			Username: "testUser",
			Password: "testPassword",
		})

		correctID := 1

		user, _ := testDB.GetById(correctID)

		want := User{
			Id:       1,
			Username: "testUser",
			Password: "testPassword",
		}

		assert.Equal(t, user, want)
	})
	t.Run("user not found", func(t *testing.T) {
		testDB := UsersInMemoryDB{}

		incorrectID := 0000

		want := "User not found"

		_, err := testDB.GetById(incorrectID)

		assert.Equal(t, err.Error(), want)

	})
}
