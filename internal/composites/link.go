package composites

import (
	in_memory "url-shortener-api/internal/adapters/db/in-memory"
	tStorage "url-shortener-api/internal/adapters/db/postgresql"
	tHttp "url-shortener-api/internal/controllers/http"
	"url-shortener-api/internal/controllers/http/interface"
	"url-shortener-api/internal/domain/services"
	"url-shortener-api/internal/domain/usecases/link"
	tClient "url-shortener-api/pkg/clients/postgresql"
)

type LinkComposite struct {
	Storage services.Storage
	Service link.Service
	UseCase tHttp.UseCase
	Handler _interface.Handler
}

func NewLinkCompositePostgres(postgresqlClient tClient.Client, address string) *LinkComposite {
	storage := tStorage.NewStoragePostgres(postgresqlClient)
	service := services.NewService(storage)
	useCase := link.NewLinkUseCase(service, address)
	handler := tHttp.NewHandler(useCase)

	return &LinkComposite{
		Storage: storage,
		Service: service,
		UseCase: useCase,
		Handler: handler,
	}
}

func NewLinkCompositeInMemory(address string) *LinkComposite {
	storage := in_memory.NewStorageInMemory()
	service := services.NewService(storage)
	useCase := link.NewLinkUseCase(service, address)
	handler := tHttp.NewHandler(useCase)

	return &LinkComposite{
		Storage: storage,
		Service: service,
		UseCase: useCase,
		Handler: handler,
	}
}
