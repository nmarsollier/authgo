package user

// FindAllUsers wrapper para obtener todos los usuarios
func FindAllUsers(ctx ...interface{}) ([]*User, error) {
	return findAll(ctx...)
}
