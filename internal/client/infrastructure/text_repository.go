// Package infrastructure содержит имплементацию репозиториев и клиентов
package infrastructure //nolint: dupl

import (
	"encoding/json"

	bolt "go.etcd.io/bbolt"

	"github.com/Nickolasll/goph-keeper/internal/client/domain"
)

// TextRepository - Имплементация репозитория для произвольных текстовых данных
type TextRepository struct {
	// DB - Интерфейс базы данных bbolt
	DB *bolt.DB
	// Crypto - Инстанс сервиса шифрования
	Crypto domain.CryptoServiceInterface
}

// Create - Сохраняет новые текстовые данные
func (r TextRepository) Create(userID string, text domain.Text) error {
	buf, err := json.Marshal(text)
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

	bkt, err := root.CreateBucketIfNotExists([]byte("Text"))
	if err != nil {
		return err
	}

	err = bkt.Put([]byte(text.ID), encrypted)
	if err != nil {
		return err
	}

	err = tx.Commit()
	if err != nil {
		return err
	}

	return err
}

// Update - Сохраняет существующие текстовые данные
func (r TextRepository) Update(userID string, text domain.Text) error {
	buf, err := json.Marshal(text)
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

		bkt := root.Bucket([]byte("Text"))
		if bkt == nil {
			return domain.ErrEntityNotFound
		}
		err = bkt.Put([]byte(text.ID), encrypted)
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

// Get - Возвращает текстовые данные по идентификатору данных и пользователя, если они существуют
func (r TextRepository) Get(userID, textID string) (domain.Text, error) {
	var text domain.Text
	var raw []byte

	err := r.DB.View(func(tx *bolt.Tx) error {
		root := tx.Bucket([]byte(userID))
		if root == nil {
			return domain.ErrEntityNotFound
		}

		bkt := root.Bucket([]byte("Text"))
		if bkt == nil {
			return domain.ErrEntityNotFound
		}
		raw = bkt.Get([]byte(textID))

		return nil
	})

	if err != nil {
		return text, err
	}

	if raw == nil {
		return text, domain.ErrEntityNotFound
	}

	decrypted, err := r.Crypto.Decrypt(raw)
	if err != nil {
		return text, err
	}

	err = json.Unmarshal(decrypted, &text)
	if err != nil {
		return text, err
	}

	return text, nil
}

// GetAll - возвращает все текстовые данные для пользователя
func (r TextRepository) GetAll(userID string) ([]domain.Text, error) {
	result := []domain.Text{}

	err := r.DB.View(func(tx *bolt.Tx) error {
		root := tx.Bucket([]byte(userID))
		if root == nil {
			return nil
		}

		bkt := root.Bucket([]byte("Text"))
		if bkt == nil {
			return nil
		}

		err := bkt.ForEach(func(_, v []byte) error {
			var text domain.Text
			decrypted, err := r.Crypto.Decrypt(v)
			if err != nil {
				return err
			}
			err = json.Unmarshal(decrypted, &text)
			if err != nil {
				return err
			}
			result = append(result, text)

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
