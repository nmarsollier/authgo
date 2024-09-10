package user

// ChangePassword cambiar la contraseña del usuario indicado
func ChangePassword(userID string, current string, newPassword string, ctx ...interface{}) error {
	user, err := findByID(userID, ctx...)
	if err != nil {
		return err
	}

	if err = user.validatePassword(current); err != nil {
		return err
	}

	if err = user.setPasswordText(newPassword); err != nil {
		return err
	}

	_, err = update(user, ctx...)

	return err
}
