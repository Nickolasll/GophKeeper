// Package jwkrepository содержит имплементацию интерфейса JWKRepositoryInterface
package jwkrepository

import (
	"encoding/json"

	bolt "go.etcd.io/bbolt"

	"github.com/lestrrat-go/jwx/v2/jwk"
	"github.com/sirupsen/logrus"

	"github.com/Nickolasll/goph-keeper/internal/client/domain"
)

// JWKRepository - Имплементация репозитория публичного ключа
type JWKRepository struct {
	// DB - Инстанс базы данных bbolt
	DB *bolt.DB
	// Crypto - Инстанс сервиса шифрования
	Crypto domain.CryptoServiceInterface
	log    *logrus.Logger
}

// Save - Сохраняет публичный ключ
func (r JWKRepository) Save(key jwk.Key) error {
	buf, err := json.Marshal(key)
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

	b, err := tx.CreateBucketIfNotExists([]byte("JWK"))
	if err != nil {
		return err
	}

	err = b.Put([]byte("Keys"), encrypted)
	if err != nil {
		return err
	}

	err = tx.Commit()
	if err != nil {
		return err
	}

	return err
}

// Get - Возвращает публичный ключ, если он существует
func (r JWKRepository) Get() (jwk.Key, error) {
	var key jwk.Key
	var raw []byte

	err := r.DB.View(func(tx *bolt.Tx) error {
		root := tx.Bucket([]byte("JWK"))
		if root == nil {
			return domain.ErrEntityNotFound
		}
		raw = root.Get([]byte("Keys"))

		return nil
	})
	if err != nil {
		return key, err
	}

	if len(raw) == 0 {
		return key, domain.ErrEntityNotFound
	}

	decrypted, err := r.Crypto.Decrypt(raw)
	if err != nil {
		return key, err
	}

	key, err = jwk.ParseKey(decrypted)
	if err != nil {
		return key, err
	}

	return key, nil
}

// Delete - Удаляет существующий ключ
func (r JWKRepository) Delete() error {
	return r.DB.Update(func(tx *bolt.Tx) error {
		root := tx.Bucket([]byte("JWK"))
		if root != nil {
			err := root.Delete([]byte("Keys"))
			if err != nil {
				return err
			}
		}

		return nil
	})
}

// New - Возвращает инстанс репозитория JWKRepository
func New(
	db *bolt.DB,
	crypto domain.CryptoServiceInterface,
	log *logrus.Logger,
) *JWKRepository {
	return &JWKRepository{
		DB:     db,
		Crypto: crypto,
		log:    log,
	}
}
