package user

// Disable deshabilita un usuario
func Disable(userID string, ctx ...interface{}) error {
	usr, err := findByID(userID, ctx...)
	if err != nil {
		return err
	}

	usr.Enabled = false

	_, err = update(usr, ctx...)

	return err
}
