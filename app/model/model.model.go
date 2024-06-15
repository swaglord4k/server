package model

import (
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"strconv"
	"time"

	"github.com/julienschmidt/httprouter"
)

type Model interface{}

func CreateFromBody[T Model](r *http.Request) (*T, error) {
	var container T

	body, err := io.ReadAll(r.Body)
	if err != nil {
		return &container, err
	}
	if len(body) == 0 {
		return nil, nil
	}
	err = json.Unmarshal(body, &container)
	if err != nil {
		return &container, err
	}
	return &container, nil
}

func GetFloatFromParams(p httprouter.Params, name string, formats ...string) (*float64, error) {
	if param := p.ByName(name); param != "" {
		formats = append(formats, "2006-01-02 15:04:05")
		if floatParam, err := strconv.ParseFloat(formats[0], 32); err != nil {
			return &floatParam, nil
		}
		return nil, fmt.Errorf("%s=%s is not a valid float", name, param)
	}
	return nil, fmt.Errorf("%s param not found", name)
}

func GetTimeFromParams(p httprouter.Params, name string) (*time.Time, error) {
	if param := p.ByName(name); param != "" {
		if floatParam, err := time.Parse("2006", param); err != nil {
			return &floatParam, nil
		}
		return nil, fmt.Errorf("%s=%s is not a valid float", name, param)
	}
	return nil, fmt.Errorf("%s param not found", name)
}

func GetIntFromParams(p httprouter.Params, name string) (*int, error) {
	if param := p.ByName(name); param != "" {
		if intParam, err := strconv.Atoi(param); err != nil {
			return &intParam, nil
		}
		return nil, fmt.Errorf("%s=%s is not a valid int", name, param)
	}
	return nil, fmt.Errorf("%s param not found", name)
}

func GetStringFromParams(p httprouter.Params, name string) (*string, error) {
	if param := p.ByName(name); param != "" {
		return &param, nil
	}
	return nil, fmt.Errorf("%s param not found", name)
}

func GetStringArrayFromParams(p httprouter.Params, name string) ([]string, error) {
	if param := p.ByName(name); param != "" {
		return []string{param}, nil
	}
	return nil, fmt.Errorf("%s param not found", name)
}
