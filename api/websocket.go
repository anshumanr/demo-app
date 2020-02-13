package api

import (
	"context"
	"fmt"
	"net"
	"net/http"
	"reflect"
	"time"

	"nhooyr.io/websocket"
	"nhooyr.io/websocket/wsjson"
)

var evChanW chan string

const (
	contentTypeJSON = "application/json"
	contentType     = "Content-Type"
)

//StartWebsocketServer ...
func StartWebsocketServer() {
	evChanW = make(chan string, 10)

	l, err := net.Listen("tcp", "localhost:5656")
	if err != nil {
		fmt.Println("failed to listen: ", err)
	}
	//defer l.Close()

	s := &http.Server{
		Handler:      http.HandlerFunc(handleWebsocket),
		ReadTimeout:  time.Second * 15,
		WriteTimeout: time.Second * 15,
	}
	//defer s.Close()

	// This starts the echo server on the listener.
	go func() {
		err := s.Serve(l)
		if err != nil {
			fmt.Printf("failed to listen and serve: %v", err)
		}
	}()
}

func handleWebsocket(w http.ResponseWriter, r *http.Request) {
	c, err := websocket.Accept(w, r, &websocket.AcceptOptions{
		InsecureSkipVerify: true,
	})
	if err != nil {
		fmt.Println(err)
		return
	}
	defer c.Close(websocket.StatusInternalError, "the sky is falling")

	ctx := r.Context()

	go read(ctx, c)

	//var v interface{}
	for {
		ev := <-evChanW
		err = wsjson.Write(ctx, c, ev)
	}

	//c.Close(websocket.StatusNormalClosure, "")
}

func read(ctx context.Context, c *websocket.Conn) {
	var (
		v   interface{}
		url string
	)

	m := make(map[string]string)

	for {
		err := wsjson.Read(ctx, c, &v)
		if err != nil {
			fmt.Println(err)
			return
		}

		val := reflect.ValueOf(v)
		if val.Kind() == reflect.Map {
			for _, key := range val.MapKeys() {
				k1 := key.String()
				v1 := fmt.Sprintf("%s", val.MapIndex(key))

				switch k1 {

				case "numberToDial":
					m["to"] = v1
				case "sessionId":
					m["request_uuid"] = v1
				default:
					m[k1] = v1
				}
			}
		}

		switch m["msgtype"] {
		case "MAKECALL":
			url = "localhost:8888/v1.0/accounts/123/call/"
			headers := map[string]string{
				contentType: contentTypeJSON,
			}
			post(m, url, headers)
		case "HOLD":
			url = "localhost:8888/v1.0/accounts/123/call/hold/" + m["request_uuid"]
			post(nil, url, nil)
		case "RESUME":
			url = "localhost:8888/v1.0/accounts/123/call/unhold/" + m["request_uuid"]
			post(nil, url, nil)
		case "DISCONNECT":
			url = "localhost:8888/v1.0/accounts/123/call/" + m["request_uuid"]
			delete(url)
		}

		fmt.Println("received: ", v, "---", m)
	}
}
