package domain

type CryptoServiceInterface interface {
	Encrypt(value []byte) ([]byte, error)
	Decrypt(value []byte) ([]byte, error)
}
