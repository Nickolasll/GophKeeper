// Package bankcardrepository содержит имплементацию интерфейса BankCardRepositoryInterface
package bankcardrepository

import (
	"encoding/json"
	"errors"

	bolt "go.etcd.io/bbolt"

	"github.com/google/uuid"
	"github.com/sirupsen/logrus"

	"github.com/Nickolasll/goph-keeper/internal/client/domain"
)

const bucketName = "BankCard"

// BankCardRepository - Имплементация репозитория для банковских карт
type BankCardRepository struct {
	// DB - Интерфейс базы данных bbolt
	DB *bolt.DB
	// Crypto - Инстанс сервиса шифрования
	Crypto domain.CryptoServiceInterface
	log    *logrus.Logger
}

// Create - Сохраняет новую банковскую карту
func (r BankCardRepository) Create(userID uuid.UUID, card *domain.BankCard) error {
	buf, err := json.Marshal(card)
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

	err = bkt.Put([]byte(card.ID.String()), encrypted)
	if err != nil {
		return err
	}

	err = tx.Commit()
	if err != nil {
		return err
	}

	return err
}

// Update - Сохраняет существующую банковскую карту
func (r BankCardRepository) Update(userID uuid.UUID, card *domain.BankCard) error {
	buf, err := json.Marshal(card)
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
		err = bkt.Put([]byte(card.ID.String()), encrypted)
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

// Get - Возвращает банковскую карту по идентификатору данных и пользователя, если она существуют
func (r BankCardRepository) Get(userID, cardID uuid.UUID) (domain.BankCard, error) {
	var card domain.BankCard
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
		raw = bkt.Get([]byte(cardID.String()))

		return nil
	})

	if err != nil {
		return card, err
	}

	if raw == nil {
		return card, domain.ErrEntityNotFound
	}

	decrypted, err := r.Crypto.Decrypt(raw)
	if err != nil {
		return card, err
	}

	err = json.Unmarshal(decrypted, &card)
	if err != nil {
		return card, err
	}

	return card, nil
}

// GetAll - возвращает все банковские карты для пользователя
func (r BankCardRepository) GetAll(userID uuid.UUID) ([]domain.BankCard, error) {
	result := []domain.BankCard{}

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
			var card domain.BankCard
			decrypted, err := r.Crypto.Decrypt(v)
			if err != nil {
				return err
			}
			err = json.Unmarshal(decrypted, &card)
			if err != nil {
				return err
			}
			result = append(result, card)

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

// ReplaceAll - Заменяет все локальные банковские карты пользователя на новые
func (r BankCardRepository) ReplaceAll(userID uuid.UUID, cards []domain.BankCard) error {
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

	for _, v := range cards {
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

// New - Возвращает инстанс репозитория CredentialsRepository
func New(
	db *bolt.DB,
	crypto domain.CryptoServiceInterface,
	log *logrus.Logger,
) *BankCardRepository {
	return &BankCardRepository{
		DB:     db,
		Crypto: crypto,
		log:    log,
	}
}
