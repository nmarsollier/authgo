package user

// Enable habilita un usuario
func Enable(userID string, deps ...interface{}) error {
	usr, err := findByID(userID, deps...)
	if err != nil {
		return err
	}

	usr.Enabled = true
	_, err = update(usr, deps...)

	return err
}
