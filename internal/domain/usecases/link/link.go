package link

import (
	"crypto/rand"
	"errors"
	"fmt"
	"math/big"
	"url-shortener-api/internal/domain/entities"
)

const (
	chars  = "_abcdefghijklmnopqrstuvwxyzABCDEFGHIJKLMNOPQRSTUVWXYZ0123456789"
	number = 10
)

const (
	errLinkIdExists = "id exists"
	errRandomizer   = "randomizer method generated the same value %d times in a row. Try increasing of chars constant"
)

type Service interface {
	GetById(id string) (*entities.Link, error)
	Create(URL string, linkId string) (*entities.Link, error)
}

type linkUseCase struct {
	service Service
	address string
}

func NewLinkUseCase(service Service, address string) *linkUseCase {
	return &linkUseCase{service: service, address: address}
}

func (l *linkUseCase) GetURLById(id string) (string, error) {
	link, err := l.service.GetById(id)
	if err != nil {
		return "", err
	}

	return link.URL, nil
}

func (l *linkUseCase) CreateLink(URL string) (string, error) {
	const maxAttempts = 3

	for i := 0; i < maxAttempts; i++ {
		linkId, err := randomizer(number, chars)
		if err != nil {
			return "", err
		}

		link, err := l.service.Create(URL, linkId)
		if err != nil {
			if err.Error() == errLinkIdExists {
				continue
			}
			return "", err
		}

		return fmt.Sprintf("http://" + l.address + "/" + link.Id), nil
	}

	return "", errors.New(fmt.Sprintf(errRandomizer, maxAttempts))
}

func randomizer(n int, chars string) (string, error) {
	ret := make([]byte, n)
	length := big.NewInt(int64(len(chars)))

	for i := 0; i < n; i++ {
		num, err := rand.Int(rand.Reader, length)
		if err != nil {
			return "", err
		}
		ret[i] = chars[num.Int64()]
	}

	return string(ret), nil
}
