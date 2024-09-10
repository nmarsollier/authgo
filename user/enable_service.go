package user

// Enable habilita un usuario
func Enable(userID string, ctx ...interface{}) error {
	usr, err := findByID(userID, ctx...)
	if err != nil {
		return err
	}

	usr.Enabled = true
	_, err = update(usr, ctx...)

	return err
}
