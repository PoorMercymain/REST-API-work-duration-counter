package router

import (
	"context"
	"encoding/json"
	"errors"
	"net/http"
	"strconv"

	"github.com/julienschmidt/httprouter"
)

const keyParams = iota

var (
	ErrParamKey  = errors.New("error no key")
	ErrParamType = errors.New("error type key")
)

type params struct {
	value httprouter.Params
}

func Params(r *http.Request) *params {
	return &params{value: r.Context().Value(keyParams).(httprouter.Params)}
}

func (p params) String(key string) (string, error) {
	value := p.value.ByName(key)
	if value == "" {
		return "", ErrParamKey
	}
	return value, nil
}

func (p params) Int(key string) (int, error) {
	value := p.value.ByName(key)
	if value == "" {
		return 0, ErrParamKey
	}
	i, err := strconv.Atoi(value)
	if err != nil {
		return 0, ErrParamType
	}
	return i, nil
}

func (p params) Uint32(key string) (uint32, error) {
	i, err := p.Int(key)
	if err != nil {
		return 0, err
	}
	return uint32(i), nil
}

func (p params) Uint16(key string) (uint16, error) {
	i, err := p.Int(key)
	if err != nil {
		return 0, err
	}
	return uint16(i), nil
}

func (p params) Bool(key string) (bool, error) {
	value := p.value.ByName(key)
	if value == "" {
		return false, ErrParamKey
	}

	res, err := strconv.ParseBool(value)
	if err != nil {
		return false, ErrParamType
	}
	return res, nil
}

func Reply(w http.ResponseWriter, data interface{}) error {
	w.Header().Set("Content-Type", "application/json; charset=utf-8")
	w.WriteHeader(http.StatusOK)
	return json.NewEncoder(w).Encode(data)
}

func WrapHandler(h http.HandlerFunc) httprouter.Handle {
	return func(w http.ResponseWriter, r *http.Request, ps httprouter.Params) {
		ctx := context.WithValue(r.Context(), keyParams, ps)
		r = r.WithContext(ctx)
		h.ServeHTTP(w, r)
	}
}
