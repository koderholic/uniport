package api

import (
	"compress/gzip"
	"context"
	"io"
	"os"
	"strings"

	"log"
	"net/http"

	"github.com/julienschmidt/httprouter"
	"github.com/justinas/alice"
	"github.com/rs/cors"

	"uniport/buckets"
	"uniport/config"
	"uniport/utils"
)

//Message ...
type Message struct {
	Code           int
	Message, Error string
	Body           interface{}
}

//Router ...
type Router struct { // Router struct would carry the httprouter instance,
	*httprouter.Router //so its methods could be verwritten and replaced with methds with wraphandler
}

//Get ...
func (router *Router) Get(path string, handler http.Handler) {
	router.GET(path, wrapHandler(handler)) // Get is an endpoint to only accept requests of method GET
}

//Post is an endpoint to only accept requests of method POST
func (router *Router) Post(path string, handler http.Handler) {
	router.POST(path, wrapHandler(handler))
}

//Put is an endpoint to only accept requests of method PUT
func (router *Router) Put(path string, handler http.Handler) {
	router.PUT(path, wrapHandler(handler))
}

//Delete is an endpoint to only accept requests of method DELETE
func (router *Router) Delete(path string, handler http.Handler) {
	router.DELETE(path, wrapHandler(handler))
}

//NewRouter is a wrapper that makes the httprouter struct a child of the router struct
func NewRouter() *Router {
	return &Router{httprouter.New()}
}

type gzipResponseWriter struct {
	io.Writer
	http.ResponseWriter
}

func (w gzipResponseWriter) Write(b []byte) (int, error) {
	return w.Writer.Write(b)
}

func gzipWrite(dataBytes []byte, httpRes http.ResponseWriter) {
	//httpRes.Header().Set("Transfer-Encoding", "gzip")
	httpRes.Header().Set("Content-Encoding", "gzip")
	gzipHandler := gzip.NewWriter(httpRes)
	defer gzipHandler.Close()
	httpResGzip := gzipResponseWriter{Writer: gzipHandler, ResponseWriter: httpRes}
	httpResGzip.Write(dataBytes)
}

func wrapHandler(httpHandler http.Handler) httprouter.Handle {
	return func(httpRes http.ResponseWriter, httpReq *http.Request, httpParams httprouter.Params) {
		ctx := context.WithValue(httpReq.Context(), "params", httpParams)
		httpReq = httpReq.WithContext(ctx)

		if !strings.Contains(httpReq.Header.Get("Accept-Encoding"), "gzip") {
			httpHandler.ServeHTTP(httpRes, httpReq)
			return
		}

		httpRes.Header().Set("Content-Encoding", "gzip")
		gzipHandler := gzip.NewWriter(httpRes)
		defer gzipHandler.Close()
		httpResGzip := gzipResponseWriter{Writer: gzipHandler, ResponseWriter: httpRes}
		httpHandler.ServeHTTP(httpResGzip, httpReq)
	}
}

func fileServe(httpRes http.ResponseWriter, httpReq *http.Request) {
	urlPath := strings.Split(httpReq.URL.String()[1:], "?")[0]
	urlPath = strings.Replace(urlPath, "//", "/", -1)

	switch config.Get().OS {
	case "ios", "android":
		if _, err := config.Asset(urlPath); err != nil {
			urlPath = "index.html"
		}
	default:
		if _, err := os.Stat(urlPath); os.IsNotExist(err) {
			urlPath = "index.html"
		}
	}

	httpRes.Header().Set("Cache-Control", "max-age=0, must-revalidate")
	httpRes.Header().Set("Pragma", "no-cache")
	httpRes.Header().Set("Expires", "0")

	if dataBytes, err := config.Asset(urlPath); err == nil {
		httpRes.Header().Add("Content-Type", config.ContentType(urlPath))
		if !strings.Contains(httpReq.Header.Get("Accept-Encoding"), "gzip") {
			httpRes.Write(dataBytes)
			return
		}
		gzipWrite(dataBytes, httpRes)
	}
}

//StartRouter ...
func StartRouter() {

	middlewares := alice.New()
	router := NewRouter()

	router.GET("/bucket/list/*bucket", func(httpRes http.ResponseWriter, httpReq *http.Request, httpParams httprouter.Params) {
		if claims := utils.VerifyJWT(httpRes, httpReq); claims == nil {
			bucketList := buckets.BucketList(httpParams.ByName("bucket"))
			httpRes.Write([]byte(strings.Join(bucketList, "\n")))
		}
	})

	router.POST("/bucket/empty/*bucket", func(httpRes http.ResponseWriter, httpReq *http.Request, httpParams httprouter.Params) {
		if claims := utils.VerifyJWT(httpRes, httpReq); claims == nil {
			errMessages := buckets.Empty(httpParams.ByName("bucket"))
			httpRes.Write([]byte(strings.Join(errMessages, "\n")))
		}
	})

	handler := cors.New(cors.Options{
		AllowedOrigins:   []string{"*"},
		AllowedMethods:   []string{"GET", "POST"},
		AllowCredentials: true,
		AllowedHeaders:   []string{"Accept", "Content-Type", "X-Auth-Token", "*"},
		Debug:            false,
	}).Handler(router)

	apiHandler(middlewares, router)
	pageHandler(middlewares, router)

	// wrapHandlerBlacklists := http.HandlerFunc(
	// 	func(httpRes http.ResponseWriter, httpReq *http.Request) {
	// 		handler.ServeHTTP(httpRes, httpReq)
	// 	})

	sMessage := "serving @ " + config.Get().Address
	println(sMessage)
	log.Println(sMessage)
	log.Fatal(http.ListenAndServe(config.Get().Address, handler))
}
