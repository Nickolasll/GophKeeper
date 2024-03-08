// Package binaryrepository содержит имплементацию интерфейса BinaryRepositoryInterface
package binaryrepository

import (
	"encoding/json"
	"errors"

	bolt "go.etcd.io/bbolt"

	"github.com/google/uuid"
	"github.com/sirupsen/logrus"

	"github.com/Nickolasll/goph-keeper/internal/client/domain"
)

const bucketName = "Binary"

// BinaryRepository - Имплементация репозитория для произвольных бинарных данных
type BinaryRepository struct {
	// DB - Интерфейс базы данных bbolt
	DB *bolt.DB
	// Crypto - Инстанс сервиса шифрования
	Crypto domain.CryptoServiceInterface
	log    *logrus.Logger
}

// Create - Сохраняет новые бинарные данные
func (r BinaryRepository) Create(userID uuid.UUID, bin domain.Binary) error {
	buf, err := json.Marshal(bin)
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

	err = bkt.Put([]byte(bin.ID.String()), encrypted)
	if err != nil {
		return err
	}

	err = tx.Commit()
	if err != nil {
		return err
	}

	return err
}

// Update - Сохраняет существующие бинарные данные
func (r BinaryRepository) Update(userID uuid.UUID, bin domain.Binary) error {
	buf, err := json.Marshal(bin)
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
		err = bkt.Put([]byte(bin.ID.String()), encrypted)
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

// Get - Возвращает бинарные данные по идентификатору данных и пользователя, если они существуют
func (r BinaryRepository) Get(userID, binID uuid.UUID) (domain.Binary, error) {
	var bin domain.Binary
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
		raw = bkt.Get([]byte(binID.String()))

		return nil
	})

	if err != nil {
		return bin, err
	}

	if raw == nil {
		return bin, domain.ErrEntityNotFound
	}

	decrypted, err := r.Crypto.Decrypt(raw)
	if err != nil {
		return bin, err
	}

	err = json.Unmarshal(decrypted, &bin)
	if err != nil {
		return bin, err
	}

	return bin, nil
}

// GetAll - возвращает все бинарные данные для пользователя
func (r BinaryRepository) GetAll(userID uuid.UUID) ([]domain.Binary, error) {
	result := []domain.Binary{}

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
			var bin domain.Binary
			decrypted, err := r.Crypto.Decrypt(v)
			if err != nil {
				return err
			}
			err = json.Unmarshal(decrypted, &bin)
			if err != nil {
				return err
			}
			result = append(result, bin)

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

// ReplaceAll - Заменяет все локальные бинарные данные пользователя на новые
func (r BinaryRepository) ReplaceAll(userID uuid.UUID, bins []domain.Binary) error {
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

	err = root.DeleteBucket([]byte(bucketName))
	if err != nil && !errors.Is(err, bolt.ErrBucketNotFound) {
		return err
	}

	bkt, err := root.CreateBucketIfNotExists([]byte(bucketName))
	if err != nil {
		return err
	}

	for _, v := range bins {
		buf, err := json.Marshal(v)
		if err != nil {
			return err
		}

		encrypted, err := r.Crypto.Encrypt(buf)
		if err != nil {
			return err
		}

		err = bkt.Put([]byte(v.ID.String()), encrypted)
		if err != nil {
			return err
		}
	}

	err = tx.Commit()
	if err != nil {
		return err
	}

	return err
}

// New - Возвращает инстанс репозитория BinaryRepository
func New(
	db *bolt.DB,
	crypto domain.CryptoServiceInterface,
	log *logrus.Logger,
) *BinaryRepository {
	return &BinaryRepository{
		DB:     db,
		Crypto: crypto,
		log:    log,
	}
}
