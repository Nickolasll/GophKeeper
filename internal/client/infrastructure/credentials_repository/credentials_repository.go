// Package credentialsrepository содержит имплементацию интерфейса CredentialsRepositoryInterface
package credentialsrepository

import (
	"encoding/json"

	bolt "go.etcd.io/bbolt"

	"github.com/google/uuid"
	"github.com/sirupsen/logrus"

	"github.com/Nickolasll/goph-keeper/internal/client/domain"
)

const bucketName = "Credentials"

// CredentialsRepository - Имплементация репозитория для логинов и паролей
type CredentialsRepository struct {
	// DB - Интерфейс базы данных bbolt
	DB *bolt.DB
	// Crypto - Инстанс сервиса шифрования
	Crypto domain.CryptoServiceInterface
	log    *logrus.Logger
}

// Create - Сохраняет новую пару логина и пароля
func (r CredentialsRepository) Create(userID uuid.UUID, cred domain.Credentials) error {
	buf, err := json.Marshal(cred)
	if err != nil {
		return err
	}

	encrypted, err := r.Crypto.Encrypt(buf)
	if err != nil {
		return err
	}

	tx, err := r.DB.Begin(true)
	if err != nil {
		return err
	}
	defer func() {
		err = tx.Rollback()
	}()

	root := tx.Bucket([]byte(userID.String()))
	if root == nil {
		return domain.ErrEntityNotFound
	}

	bkt, err := root.CreateBucketIfNotExists([]byte(bucketName))
	if err != nil {
		return err
	}

	err = bkt.Put([]byte(cred.ID.String()), encrypted)
	if err != nil {
		return err
	}

	err = tx.Commit()
	if err != nil {
		return err
	}

	return err
}

// Update - Сохраняет существующую пару логина и пароля
func (r CredentialsRepository) Update(userID uuid.UUID, cred domain.Credentials) error {
	buf, err := json.Marshal(cred)
	if err != nil {
		return err
	}

	encrypted, err := r.Crypto.Encrypt(buf)
	if err != nil {
		return err
	}

	err = r.DB.Update(func(tx *bolt.Tx) error {
		root := tx.Bucket([]byte(userID.String()))
		if root == nil {
			return domain.ErrEntityNotFound
		}

		bkt := root.Bucket([]byte(bucketName))
		if bkt == nil {
			return domain.ErrEntityNotFound
		}
		err = bkt.Put([]byte(cred.ID.String()), encrypted)
		if err != nil {
			return err
		}

		return nil
	})
	if err != nil {
		return err
	}

	return nil
}

// Get - Возвращает логин и пароль по идентификатору данных и пользователя, если они существуют
func (r CredentialsRepository) Get(userID, credID uuid.UUID) (domain.Credentials, error) {
	var cred domain.Credentials
	var raw []byte

	err := r.DB.View(func(tx *bolt.Tx) error {
		root := tx.Bucket([]byte(userID.String()))
		if root == nil {
			return domain.ErrEntityNotFound
		}

		bkt := root.Bucket([]byte(bucketName))
		if bkt == nil {
			return domain.ErrEntityNotFound
		}
		raw = bkt.Get([]byte(credID.String()))

		return nil
	})

	if err != nil {
		return cred, err
	}

	if raw == nil {
		return cred, domain.ErrEntityNotFound
	}

	decrypted, err := r.Crypto.Decrypt(raw)
	if err != nil {
		return cred, err
	}

	err = json.Unmarshal(decrypted, &cred)
	if err != nil {
		return cred, err
	}

	return cred, nil
}

// GetAll - возвращает все логины и пароли для пользователя
func (r CredentialsRepository) GetAll(userID uuid.UUID) ([]domain.Credentials, error) {
	result := []domain.Credentials{}

	err := r.DB.View(func(tx *bolt.Tx) error {
		root := tx.Bucket([]byte(userID.String()))
		if root == nil {
			return nil
		}

		bkt := root.Bucket([]byte(bucketName))
		if bkt == nil {
			return nil
		}

		err := bkt.ForEach(func(_, v []byte) error {
			var cred domain.Credentials
			decrypted, err := r.Crypto.Decrypt(v)
			if err != nil {
				return err
			}
			err = json.Unmarshal(decrypted, &cred)
			if err != nil {
				return err
			}
			result = append(result, cred)

			return nil
		})
		if err != nil {
			return err
		}

		return nil
	})

	if err != nil {
		return result, err
	}

	return result, nil
}

// New - Возвращает инстанс репозитория CredentialsRepository
func New(
	db *bolt.DB,
	crypto domain.CryptoServiceInterface,
	log *logrus.Logger,
) *CredentialsRepository {
	return &CredentialsRepository{
		DB:     db,
		Crypto: crypto,
		log:    log,
	}
}
