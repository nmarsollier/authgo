package user

// Revoke Le revoca los permisos enviados por parametros
func Revoke(userID string, permissions []string, deps ...interface{}) error {
	user, err := findByID(userID, deps...)
	if err != nil {
		return err
	}

	for _, value := range permissions {
		user.revoke(value)
	}
	_, err = update(user, deps...)

	return err
}
