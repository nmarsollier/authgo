package user

// Grant Le habilita los permisos enviados por parametros
func Grant(userID string, permissions []string, ctx ...interface{}) error {
	user, err := findByID(userID, ctx...)
	if err != nil {
		return err
	}

	for _, value := range permissions {
		user.grant(value)
	}
	_, err = update(user, ctx...)

	return err
}
