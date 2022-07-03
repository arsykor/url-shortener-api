package in_memory

import (
	"errors"
	"fmt"
	"sync"
	"url-shortener-api/internal/domain/entities"
	"url-shortener-api/internal/domain/services"
)

const (
	errURLExists    = "the link for URL %s already exists, try using ID = %s"
	errLinkIdExists = "id exists"
	errNoLinkInDB   = "there is no passed link in DB"
)

type linkStorageInMemory struct {
	mutex        sync.Mutex
	mapLinkToURL map[string]string
	mapURLToLink map[string]string
	arrLinks     []entities.LinkInMemory
}

func NewStorageInMemory() services.Storage {
	return &linkStorageInMemory{mapLinkToURL: make(map[string]string), mapURLToLink: make(map[string]string)}
}

func (s *linkStorageInMemory) Create(URL string, linkId string) (*entities.Link, error) {
	ch := make(chan error)

	go s.createLink(ch, URL, linkId)

	res := <-ch

	if res != nil {
		return nil, res
	}

	return &entities.Link{Id: linkId, URL: URL}, nil
}

func (s *linkStorageInMemory) createLink(ch chan error, URL string, linkId string) {
	s.mutex.Lock()
	defer s.mutex.Unlock()

	_, ok := s.mapLinkToURL[linkId]
	if ok {
		ch <- errors.New(errLinkIdExists)
		return
	}

	v, ok := s.mapURLToLink[URL]
	if ok {
		ch <- errors.New(fmt.Sprintf(errURLExists, URL, v))
		return
	}

	s.mapLinkToURL[linkId] = URL
	s.mapURLToLink[URL] = linkId

	s.arrLinks = append(s.arrLinks, entities.LinkInMemory{
		Id:  linkId,
		URL: URL,
	})
	for _, item := range s.arrLinks {
		fmt.Println(item)
	}

	close(ch)
}

func (s *linkStorageInMemory) GetOne(id string) (*entities.Link, error) {
	v, ok := s.mapLinkToURL[id]
	if !ok {
		return nil, errors.New(errNoLinkInDB)
	}

	return &entities.Link{Id: id, URL: v}, nil
}
