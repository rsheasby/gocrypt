package gocrypt

type PasswordHasher interface {
	// HashPassword returns a hash of the provided password for storage in a database.
	HashPassword(password string) (hash string, err error)
	// ValidatePassword takes a password and the stored hash, and returns whether the password is valid.
	ValidatePassword(password string, hash string) (isValid bool, err error)
}
