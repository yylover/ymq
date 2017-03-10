package ymqd

import (
	"fmt"
	"io"
	"log"
	"net/http"
	"os"
	"ymq/internal/httpapi"
	"ymq/internal/version"

	"github.com/julienschmidt/httprouter"
)

/**
 * /pub  发送一个message
 * /mpub 发送多个message
 * /config 配置 哪些支持热配置
 * /info 信息查看
 * /ping
 * /status
 *
 *
 * /topic/create
 * /topic/delete
 * /topic/empty
 * /channel/create
 * /channel/delete
 * /channel/empty
 * /channel/pause 暂停?
 *
 * /debug/后续支持
 */

/**
 * 1. response 中写了w 之后，只能是ERROR
 */

//HTTPServer http 服务器
type HTTPServer struct {
	ctx    *context
	router http.Handler
	logger Logger
}

//LoggerLocal 本地输出
type LoggerLocal struct {
}

//Output 打印消息
func (l *LoggerLocal) Output(maxdepth int, s string) error {
	fmt.Printf("%d %s\n", maxdepth, s)
	return nil
}

//NewHTTPServer 创建http服务器
func NewHTTPServer(ctx *context) (*HTTPServer, error) {
	router := httprouter.New()
	logger := &LoggerLocal{}

	server := &HTTPServer{
		ctx:    ctx,
		router: router,
		logger: logger,
	}

	router.HandleMethodNotAllowed = true
	router.NotFound = httpapi.LogNotFoundHandler(logger)
	router.MethodNotAllowed = httpapi.LogMethodNotAllowedHandler(logger)
	router.PanicHandler = httpapi.LogPanicHandler(logger)

	// router.HandlerFunc("GET", "/hello", helloServer)
	router.GET("/pubtest", server.pub)
	router.GET("/hello/:name", server.helloName)
	router.GET("/info", httpapi.Decorate(server.doInfo, httpapi.Log(logger), httpapi.V1))
	router.GET("/ping", httpapi.Decorate(server.ping, httpapi.Log(logger), httpapi.V1))
	router.GET("/status", httpapi.Decorate(server.status, httpapi.Log(logger), httpapi.V1))
	router.GET("/config", httpapi.Decorate(server.config, httpapi.Log(logger), httpapi.V1))
	router.GET("/pub", httpapi.Decorate(server.doPub, httpapi.Log(logger), httpapi.V1))

	router.GET("/topic/create", httpapi.Decorate(server.topicCreate, httpapi.Log(logger), httpapi.V1))
	router.GET("/topic/delte", httpapi.Decorate(server.topicDelete, httpapi.Log(logger), httpapi.V1))
	router.GET("/topic/empty", httpapi.Decorate(server.topicEmpty, httpapi.Log(logger), httpapi.V1))

	router.GET("/channel/create", httpapi.Decorate(server.chanCreate, httpapi.Log(logger), httpapi.V1))
	router.GET("/channel/delete", httpapi.Decorate(server.chanDelete, httpapi.Log(logger), httpapi.V1))
	router.GET("/channel/empty", httpapi.Decorate(server.chanEmpty, httpapi.Log(logger), httpapi.V1))

	router.HandlerFunc("GET", "/hello", helloServer)

	return server, nil
}

// Serve 初始化服务器
func (h *HTTPServer) Serve() error {
	// h.server.ListenAndServe()
	// router := httprouter.New()
	// http.HandleFunc("/hello", helloServer)
	log.Fatal(http.ListenAndServe(":12345", h.router))
	return nil
}

//demo1 : go 最原先的处理原型
func helloServer(w http.ResponseWriter, req *http.Request) {
	// io.WriteString(w, "Hello world")
	panic("message")
}

//demo2 : httprouter 对路由的封装，添加了httprouter.Params参数
func (h *HTTPServer) pub(w http.ResponseWriter, req *http.Request, ps httprouter.Params) {
	io.WriteString(w, "pub message")
}

func (h *HTTPServer) helloName(w http.ResponseWriter, req *http.Request, ps httprouter.Params) {
	fmt.Fprintf(w, "hello, %s!\n", ps.ByName("name"))
}

//Demo3 : 在httprouter 的基础上，添加了装饰器，自动打印Log，自动对返回的数据进行处理
func (h *HTTPServer) doInfo(w http.ResponseWriter, req *http.Request, ps httprouter.Params) (interface{}, error) {
	hostname, err := os.Hostname()
	if err != nil {
		return nil, httpapi.Err{Code: 500, Text: err.Error()}
	}

	return struct {
		Version  string `json:"version"`
		Hostname string `json:"hostname"`
	}{
		Version:  version.String("ymqd"),
		Hostname: hostname,
	}, nil
}

func (h *HTTPServer) doPub(w http.ResponseWriter, req *http.Request, ps httprouter.Params) (interface{}, error) {
	return nil, nil
}

func (h *HTTPServer) topicCreate(w http.ResponseWriter, req *http.Request, ps httprouter.Params) (interface{}, error) {

	return nil, nil
}

func (h *HTTPServer) topicDelete(w http.ResponseWriter, req *http.Request, ps httprouter.Params) (interface{}, error) {
	return nil, nil
}

func (h *HTTPServer) topicEmpty(w http.ResponseWriter, req *http.Request, ps httprouter.Params) (interface{}, error) {
	return nil, nil
}

func (h *HTTPServer) chanCreate(w http.ResponseWriter, req *http.Request, ps httprouter.Params) (interface{}, error) {
	return nil, nil
}
func (h *HTTPServer) chanDelete(w http.ResponseWriter, req *http.Request, ps httprouter.Params) (interface{}, error) {
	return nil, nil
}
func (h *HTTPServer) chanEmpty(w http.ResponseWriter, req *http.Request, ps httprouter.Params) (interface{}, error) {
	return nil, nil
}
func (h *HTTPServer) ping(w http.ResponseWriter, req *http.Request, ps httprouter.Params) (interface{}, error) {
	return nil, nil
}
func (h *HTTPServer) status(w http.ResponseWriter, req *http.Request, ps httprouter.Params) (interface{}, error) {
	return nil, nil
}
func (h *HTTPServer) config(w http.ResponseWriter, req *http.Request, ps httprouter.Params) (interface{}, error) {
	return nil, nil
}
