
type UserProfile struct{
	ID string `json:"id" bson:"_id"`
	Username string `json:"username" bson:"username"`
	hobbies []string `json:"password" bson:"password"`
}


func(u *UserProfile)createUser(username string, password string)error{
	bytes ,err := bcrypt.GenerateFromPassword([]byte(password),14)
	if err != nil {
		user.Password = string(bytes)
	}

}

func(u *UserProfile)login( username string, password string)error{
	
}
