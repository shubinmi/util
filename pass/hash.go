package pass

import "golang.org/x/crypto/bcrypt"

func Hash(pass string) (hash string, err error) {
	password := []byte(pass)
	hashedPassword, err := bcrypt.GenerateFromPassword(password, bcrypt.DefaultCost)
	if err != nil {
		return
	}
	hash = string(hashedPassword)
	return
}

func Verify(hash, pass string) error {
	return bcrypt.CompareHashAndPassword([]byte(hash), []byte(pass))
}
