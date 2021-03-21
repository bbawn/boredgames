package services

import (
	"net/http"

	daoerr "github.com/bbawn/boredgames/internal/dao/errors"
	"github.com/bbawn/boredgames/internal/games/set"
)

func httpStatus(err error) int {
	switch err.(type) {
	case daoerr.AlreadyExistsError:
		return http.StatusConflict
	case daoerr.InternalError:
		return http.StatusInternalServerError
	case daoerr.NotFoundError:
		return http.StatusNotFound
	case set.InvalidArgError:
		return http.StatusBadRequest
	case set.InvalidStateError:
		return http.StatusConflict
	default:
		return http.StatusInternalServerError
	}
}
