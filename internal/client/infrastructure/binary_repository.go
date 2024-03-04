// Package infrastructure содержит имплементацию репозиториев и клиентов
package infrastructure //nolint: dupl

import (
	"encoding/json"

	bolt "go.etcd.io/bbolt"

	"github.com/Nickolasll/goph-keeper/internal/client/domain"
)

// BinaryRepository - Имплементация репозитория для произвольных бинарных данных
type BinaryRepository struct {
	// DB - Интерфейс базы данных bbolt
	DB *bolt.DB
	// Crypto - Инстанс сервиса шифрования
	Crypto domain.CryptoServiceInterface
}

// Create - Сохраняет новые бинарные данные
func (r BinaryRepository) Create(userID string, bin domain.Binary) error {
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

	root := tx.Bucket([]byte(userID))
	if root == nil {
		return domain.ErrEntityNotFound
	}

	bkt, err := root.CreateBucketIfNotExists([]byte("Binary"))
	if err != nil {
		return err
	}

	err = bkt.Put([]byte(bin.ID), encrypted)
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
func (r BinaryRepository) Update(userID string, bin domain.Binary) error {
	buf, err := json.Marshal(bin)
	if err != nil {
		return err
	}

	encrypted, err := r.Crypto.Encrypt(buf)
	if err != nil {
		return err
	}

	err = r.DB.Update(func(tx *bolt.Tx) error {
		root := tx.Bucket([]byte(userID))
		if root == nil {
			return domain.ErrEntityNotFound
		}

		bkt := root.Bucket([]byte("Binary"))
		if bkt == nil {
			return domain.ErrEntityNotFound
		}
		err = bkt.Put([]byte(bin.ID), encrypted)
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
func (r BinaryRepository) Get(userID, binID string) (domain.Binary, error) {
	var bin domain.Binary
	var raw []byte

	err := r.DB.View(func(tx *bolt.Tx) error {
		root := tx.Bucket([]byte(userID))
		if root == nil {
			return domain.ErrEntityNotFound
		}

		bkt := root.Bucket([]byte("Binary"))
		if bkt == nil {
			return domain.ErrEntityNotFound
		}
		raw = bkt.Get([]byte(binID))

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
func (r BinaryRepository) GetAll(userID string) ([]domain.Binary, error) {
	result := []domain.Binary{}

	err := r.DB.View(func(tx *bolt.Tx) error {
		root := tx.Bucket([]byte(userID))
		if root == nil {
			return nil
		}

		bkt := root.Bucket([]byte("Binary"))
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
