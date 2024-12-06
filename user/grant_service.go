package user

// Grant Le habilita los permisos enviados por parametros
func Grant(userID string, permissions []string, deps ...interface{}) error {
	user, err := findByID(userID, deps...)
	if err != nil {
		return err
	}

	for _, value := range permissions {
		user.grant(value)
	}
	err = update(user, deps...)

	return err
}
