package httpapi

import (
	"encoding/json"
	"fmt"
	"net/http"
	"time"
	"ymq/internal/app"

	"github.com/julienschmidt/httprouter"
)

//APIHandler 接口处理函数原型,参数是httprouter 的参数,返回值包括Err 类型的error 和interface{}类型，装饰器自动根据返回值
type APIHandler func(http.ResponseWriter, *http.Request, httprouter.Params) (interface{}, error)

//Decorator 装饰器
type Decorator func(APIHandler) APIHandler

//Err 返回值错误
type Err struct {
	Code int
	Text string
}

func (e Err) Error() string {
	return e.Text
}

//Decorate 装饰函数
func Decorate(f APIHandler, ds ...Decorator) httprouter.Handle {
	decorated := f
	for _, decorate := range ds {
		decorated = decorate(decorated)
	}

	return func(w http.ResponseWriter, req *http.Request, ps httprouter.Params) {
		decorated(w, req, ps)
	}
}

//LogNotFoundHandler 没有找到时处理函数
func LogNotFoundHandler(l app.Logger) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, req *http.Request) {
		// fmt.Fprint(w, "methond not found")
		Decorate(func(w http.ResponseWriter, req *http.Request, ps httprouter.Params) (interface{}, error) {
			return nil, Err{404, "Not Found"}
		}, Log(l), V1)(w, req, nil)
	})
}

// LogMethodNotAllowedHandler 方法不允许
func LogMethodNotAllowedHandler(l app.Logger) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, req *http.Request) {
		// fmt.Fprint(w, "methond not allowed")
		Decorate(func(w http.ResponseWriter, req *http.Request, ps httprouter.Params) (interface{}, error) {
			return nil, Err{404, "Not Allowed"}
		}, Log(l), V1)(w, req, nil)
	})
}

//LogPanicHandler panic handler
func LogPanicHandler(l app.Logger) func(http.ResponseWriter, *http.Request, interface{}) {
	return func(w http.ResponseWriter, req *http.Request, ps interface{}) {
		l.Output(2, fmt.Sprintf("ERROR: panic in HTTP handler - %s", ps))
		Decorate(func(w http.ResponseWriter, req *http.Request, ps httprouter.Params) (interface{}, error) {
			return nil, Err{500, "INTERNAL_ERROR"}
		}, Log(l), V1)(w, req, nil)
	}
}

//Log 装饰器函数
func Log(l app.Logger) Decorator {
	return func(f APIHandler) APIHandler {
		return func(w http.ResponseWriter, req *http.Request, ps httprouter.Params) (interface{}, error) {
			start := time.Now()

			response, err := f(w, req, ps)
			elapsed := time.Since(start)
			status := 200
			if e, ok := err.(Err); ok {
				status = e.Code
			}

			l.Output(2, fmt.Sprintf("%d %s %s (%s) %s", status, req.Method, req.URL.RequestURI(), req.RemoteAddr, elapsed))
			return response, err
		}
	}
}

//V1 第一版本的处理器
func V1(f APIHandler) APIHandler {
	return func(w http.ResponseWriter, req *http.Request, ps httprouter.Params) (interface{}, error) {
		data, err := f(w, req, ps)

		if err != nil {
			responseV1(w, err.(Err).Code, err)
			return nil, nil
		}

		responseV1(w, 200, data)
		return nil, nil
	}
}

func responseV1(w http.ResponseWriter, code int, data interface{}) {
	var isJSON bool
	var err error
	var response []byte

	if code == 200 {
		switch data.(type) {
		case string:
			response = []byte(data.(string))
		case []byte:
			response = data.([]byte)
		case nil:
			response = []byte{}
		default:
			isJSON = true
			response, err = json.Marshal(data)
			if err != nil {
				code = 500
			}
		}
	}

	if err != nil {
		isJSON = false
		response = []byte(fmt.Sprintf(`{"message":"%s"}`, err.Error()))
	}

	if isJSON {
		w.Header().Set("Content-Type", "application/json; charset=utf-8")
	}
	w.Header().Set("X-YYMQ-Content-Type", "ymq; version=1.0")
	w.WriteHeader(code)
	w.Write(response)
}
