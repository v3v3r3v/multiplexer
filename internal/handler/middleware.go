package handler

import (
	"log"
	"net/http"
)

/// Limits processing of concurrent requests
/// all requests over limit will wait
func LimitConcurrentRequests(handlerFunc http.HandlerFunc, limit int) http.HandlerFunc {
	semaphore := make(chan struct{}, limit)

	return func (writer http.ResponseWriter, request *http.Request) {
		semaphore <- struct{}{}
		log.Println("Start processing request")
		defer func() {
			<-semaphore
			log.Println("Done processing request")
		}()
		handlerFunc(writer, request)
	}
}

func HttpMethods(handlerFunc http.HandlerFunc, methods... string) http.HandlerFunc {
	m := make(map[string]struct{}, len(methods))

	for _, v := range methods {
		m[v] = struct{}{}
	}

	return func (writer http.ResponseWriter, request *http.Request) {
		if _, ok := m[request.Method]; ok {
			handlerFunc(writer, request)
			return
		}

		writer.WriteHeader(http.StatusMethodNotAllowed)
	}
}
