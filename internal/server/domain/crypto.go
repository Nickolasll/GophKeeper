package domain

// CryptoServiceInterface - Интерфейс сервиса шифрования
type CryptoServiceInterface interface {
	// Encrypt - Зашифровывает данные
	Encrypt(value []byte) ([]byte, error)
	// Decrypt - Расщифровывает данные
	Decrypt(value []byte) ([]byte, error)
}
