package user

// Granted verifica si el usuario tiene el permiso
func Granted(userID string, permission string, ctx ...interface{}) bool {
	usr, err := findByID(userID, ctx...)
	if err != nil {
		return false
	}

	return usr.granted(permission)
}
