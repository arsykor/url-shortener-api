package http

import (
	"context"
	"encoding/json"
	"fmt"
	"github.com/gin-gonic/gin"
	"github.com/pkg/errors"
	"io/ioutil"
	"net/http"
	"net/url"
	"unicode/utf8"
	"url-shortener-api/internal/domain/entities"
)

const (
	errNoValue       = "there is no %s specified in the body"
	errTooManyValues = "there is supposed to be one %s in the body"
	errEmptyValue    = "the value of %s is empty"
	errInvalidValue  = "invalid %s"
	errValueTooLong  = "%s has too many characters, max is %d"
)

type UseCase interface {
	GetURLById(ctx context.Context, id string) (string, error)
	CreateLink(ctx context.Context, URL string) (string, error)
}

type linkHandler struct {
	linkUseCase UseCase
}

type Error struct {
	Code    int    `json:"code"`
	Message string `json:"message"`
}

func NewHandler(linkUseCase UseCase) *linkHandler {
	return &linkHandler{linkUseCase: linkUseCase}
}

func (h *linkHandler) Register(router *gin.Engine) {
	router.POST("/", h.CreateLink)
	router.GET("/:id", h.GetURLById)
}

func (h *linkHandler) CreateLink(ctx *gin.Context) {
	var Link []entities.LinkRequest

	jsonData, err := ioutil.ReadAll(ctx.Request.Body)
	if err != nil {
		returnError(ctx, err)
		return
	}

	err = json.Unmarshal(jsonData, &Link)
	if err != nil {
		returnError(ctx, err)
		return
	}

	request := entities.LinkRequest{}
	paramName := request.ParamName()

	err = paramValidation(&Link, paramName)
	if err != nil {
		returnError(ctx, err)
		return
	}

	_, err = url.ParseRequestURI(Link[0].URL)
	if err != nil {
		returnError(ctx, errors.New(fmt.Sprintf(errInvalidValue, paramName)))
		return
	}

	linkOut, err := h.linkUseCase.CreateLink(ctx, Link[0].URL)
	if err != nil {
		returnError(ctx, err)
		return
	}

	ctx.IndentedJSON(http.StatusOK, entities.LinkResponse{Link: linkOut})
}

func (h *linkHandler) GetURLById(ctx *gin.Context) {
	id := ctx.Param("id")

	link, err := h.linkUseCase.GetURLById(ctx, id)
	if err != nil {
		returnError(ctx, err)
		return
	}

	ctx.Redirect(http.StatusFound, link)
}

func paramValidation(Link *[]entities.LinkRequest, paramName string) error {
	const maxURLLength = 100

	switch l := len(*Link); {
	case l == 0:
		return errors.New(fmt.Sprintf(errNoValue, paramName))
	case l > 1:
		return errors.New(fmt.Sprintf(errTooManyValues, paramName))
	}

	URL := (*Link)[0].URL
	switch l := utf8.RuneCountInString(URL); {
	case l == 0:
		return errors.New(fmt.Sprintf(errEmptyValue, paramName))
	case l > maxURLLength:
		return errors.New(fmt.Sprintf(errValueTooLong, paramName, maxURLLength))
	}

	return nil
}

func returnError(context *gin.Context, err error) {
	context.JSON(http.StatusBadRequest, Error{
		Code:    http.StatusBadRequest,
		Message: err.Error(),
	})
}
