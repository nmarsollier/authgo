package user

// Disable deshabilita un usuario
func Disable(userID string, deps ...interface{}) error {
	usr, err := findByID(userID, deps...)
	if err != nil {
		return err
	}

	usr.Enabled = false

	err = update(usr, deps...)

	return err
}
