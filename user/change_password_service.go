package user

// ChangePassword cambiar la contraseña del usuario indicado
func ChangePassword(userID string, current string, newPassword string, deps ...interface{}) error {
	user, err := findByID(userID, deps...)
	if err != nil {
		return err
	}

	if err = user.validatePassword(current); err != nil {
		return err
	}

	if err = user.setPasswordText(newPassword); err != nil {
		return err
	}

	_, err = update(user, deps...)

	return err
}
