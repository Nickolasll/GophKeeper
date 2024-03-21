// Package sessionrepository содержит имплементацию интерфейса SessionRepositoryInterface
package sessionrepository

import (
	"encoding/json"

	bolt "go.etcd.io/bbolt"

	"github.com/sirupsen/logrus"

	"github.com/Nickolasll/goph-keeper/internal/client/domain"
)

const bucketName = "ActiveSession"
const keyName = "Session"

// SessionRepository - Имплементация репозитория сессий
type SessionRepository struct {
	// DB - Инстанс базы данных bbolt
	DB *bolt.DB
	// Crypto - Инстанс сервиса шифрования
	Crypto domain.CryptoServiceInterface
	log    *logrus.Logger
}

// Save - Сохраняет новую сессию
func (r SessionRepository) Save(session domain.Session) error {
	buf, err := json.Marshal(session)
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

	b, err := tx.CreateBucketIfNotExists([]byte(bucketName))
	if err != nil {
		return err
	}

	err = b.Put([]byte(keyName), encrypted)
	if err != nil {
		return err
	}
	_, err = tx.CreateBucketIfNotExists([]byte(session.UserID.String()))
	if err != nil {
		return err
	}

	err = tx.Commit()
	if err != nil {
		return err
	}

	return err
}

// Get - Возвращает последнюю сессию, если она существует
func (r SessionRepository) Get() (*domain.Session, error) {
	var session domain.Session
	var raw []byte

	err := r.DB.View(func(tx *bolt.Tx) error {
		root := tx.Bucket([]byte(bucketName))
		if root == nil {
			return domain.ErrEntityNotFound
		}
		raw = root.Get([]byte(keyName))

		return nil
	})
	if err != nil {
		return &session, err
	}

	if raw == nil {
		return &session, domain.ErrEntityNotFound
	}

	decrypted, err := r.Crypto.Decrypt(raw)
	if err != nil {
		return &session, err
	}

	err = json.Unmarshal(decrypted, &session)
	if err != nil {
		return &session, err
	}

	return &session, nil
}

// Delete - Удаляет существуюущую сессию
func (r SessionRepository) Delete() error {
	return r.DB.Update(func(tx *bolt.Tx) error {
		root := tx.Bucket([]byte(bucketName))
		if root != nil {
			err := root.Delete([]byte(keyName))
			if err != nil {
				return err
			}
		}

		return nil
	})
}

// New - Возвращает инстанс репозитория SessionRepository
func New(
	db *bolt.DB,
	crypto domain.CryptoServiceInterface,
	log *logrus.Logger,
) *SessionRepository {
	return &SessionRepository{
		DB:     db,
		Crypto: crypto,
		log:    log,
	}
}
