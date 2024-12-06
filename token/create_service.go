package token

// Create crea un nuevo token y lo almacena en la db
func Create(userID string, deps ...interface{}) (*Token, error) {
	token, err := insert(userID, deps...)
	if err != nil {
		return nil, err
	}

	cacheAdd(token)

	return token, nil
}
