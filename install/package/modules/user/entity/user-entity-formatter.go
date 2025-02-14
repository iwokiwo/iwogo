package entity

type UserFormatter struct {
	ID     int    `json:"id"`
	Name   string `json:"name"`
	Email  string `json:"email"`
	Active int    `json:"active"`
	Role   string `json:"role"`
	Token  string `json:"token"`
}

func FormatUser(user User, token string) UserFormatter {
	formatter := UserFormatter{
		ID:     user.ID,
		Name:   user.Name,
		Email:  user.Email,
		Active: user.Active,
		Role:   user.Role,
		Token:  token,
	}

	return formatter
}

func FormatProfile(user User) UserFormatter {
	formatter := UserFormatter{
		ID:     user.ID,
		Name:   user.Name,
		Email:  user.Email,
		Active: user.Active,
	}
	return formatter
}

func FormatUsers(users []User) []UserFormatter {
	userFormatter := []UserFormatter{}

	for _, user := range users {
		formatter := FormatProfile(user)
		userFormatter = append(userFormatter, formatter)
	}

	return userFormatter
}
