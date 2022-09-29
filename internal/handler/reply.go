package handler

import (
	"net/http"

	"github.com/PoorMercymain/REST-API-work-duration-counter/pkg/router"
)

func reply(w http.ResponseWriter, data interface{}) error {
	return router.Reply(w, data)
}
