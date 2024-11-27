package token

// Find busca un token en la db
func Find(tokenID string, deps ...interface{}) (*Token, error) {
	token, err := findByID(tokenID, deps...)
	if err != nil {
		return nil, err
	}

	cacheAdd(token)

	return token, nil
}
