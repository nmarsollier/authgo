package user

// Granted verifica si el usuario tiene el permiso
func Granted(userID string, permission string, deps ...interface{}) bool {
	usr, err := findByID(userID, deps...)
	if err != nil {
		return false
	}

	return usr.granted(permission)
}
