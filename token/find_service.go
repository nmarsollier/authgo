package token

// Find busca un token en la db
func Find(tokenID string, ctx ...interface{}) (*Token, error) {
	token, err := findByID(tokenID, ctx...)
	if err != nil {
		return nil, err
	}

	cacheAdd(token)

	return token, nil
}
