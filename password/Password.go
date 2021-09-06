package password

import (
	"github.com/meowalien/go-util/math"
	"golang.org/x/crypto/bcrypt"
)


const RecommendSaltLength = 24

// HashPassword will hash the given password with salt and Pepper, then return the hashed password and salt.
func HashPassword(password string, pepper string,saltLength int) (string, string, error) {
	if password =="" {
		panic("The password should not be empty")
	}
	if pepper =="" {
		panic("The pepper should not be empty")
	}
	if saltLength <= 0 {
		panic("The salt length should not be less than or equal to 0")
	}
	randomSalt := math.RandomString(saltLength)
	p := append([]byte(password), pepper+randomSalt...)
	bytes, err := bcrypt.GenerateFromPassword(p, bcrypt.DefaultCost)
	return string(bytes), randomSalt, err
}

// CheckPasswordHash check if the given password, hashedPassword, and salt match.
func CheckPasswordHash(password, hash, pepper, salt string) bool {
	if password =="" {
		panic("The password should not be empty")
	}
	if pepper =="" {
		panic("The pepper should not be empty")
	}
	if salt =="" {
		panic("The salt should not be empty")
	}
	err := bcrypt.CompareHashAndPassword([]byte(hash), append([]byte(password), pepper+salt...))
	return err == nil
}
