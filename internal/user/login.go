package user

import "go.mongodb.org/mongo-driver/bson"

func LoginId(info User) (bool, error) {
	var result User
	result, err := FindOneUser(
		bson.M{"id": info.Id},
	)
	return result.Password == info.Password, err
}

func LoginEmail(info User) (bool, error) {
	// var result User
	_, err := FindOneUser(
		bson.M{"email": info.Email, "password": info.Password},
	)
	return err != nil, err
}

func LoginName(info User) (User, error) {
	var result User
	result, err := FindOneUser(
		bson.M{"name": info.Name, "password": info.Password},
	)
	return result, err
}
