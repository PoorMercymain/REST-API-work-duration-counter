package service

import "github.com/PoorMercymain/REST-API-work-duration-counter/internal/handler"

func RedirectErrors(err error) {
	handler.PrintErrors(err)
}
