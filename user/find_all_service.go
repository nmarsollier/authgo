package user

// FindAllUsers wrapper para obtener todos los usuarios
func FindAllUsers(ctx ...interface{}) ([]*UserResponse, error) {

	user, err := findAll(ctx...)

	if err != nil {
		return nil, err
	}

	result := []*UserResponse{}
	for i := 0; i < len(user); i = i + 1 {
		result = append(result, &UserResponse{
			Id:          user[i].ID.Hex(),
			Name:        user[i].Name,
			Permissions: user[i].Permissions,
			Login:       user[i].Login,
			Enabled:     user[i].Enabled,
		})
	}

	return result, nil
}
