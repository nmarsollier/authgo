package user

// FindAllUsers wrapper para obtener todos los usuarios
func FindAllUsers(deps ...interface{}) ([]*UserData, error) {

	user, err := findAll(deps...)

	if err != nil {
		return nil, err
	}

	result := []*UserData{}
	for i := 0; i < len(user); i = i + 1 {
		result = append(result, &UserData{
			Id:          user[i].ID,
			Name:        user[i].Name,
			Permissions: user[i].Permissions,
			Login:       user[i].Login,
			Enabled:     user[i].Enabled,
		})
	}

	return result, nil
}
