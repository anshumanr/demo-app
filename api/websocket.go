package api

import (
	"context"
	"fmt"
	"io"
	"net/http"
	"time"

	"github.com/labstack/echo/v4"
	"nhooyr.io/websocket"
)

func (s *Server) handleWebsocket(c echo.Context) error {
	r := c.Request()
	w := c.Response().Writer

	return echoServer(w, r)
}

func echoServer(w http.ResponseWriter, r *http.Request) error {
	c, err := websocket.Accept(w, r, &websocket.AcceptOptions{
		InsecureSkipVerify: true,
	})
	if err != nil {
		fmt.Println("Error: ", err)
		return err
	}
	defer c.Close(websocket.StatusInternalError, "the sky is falling")

	err = readws(r.Context(), c)
	if websocket.CloseStatus(err) == websocket.StatusNormalClosure {
		return nil
	}
	if err != nil {
		return fmt.Errorf("failed to echo with %v: %w", r.RemoteAddr, err)
	}

	return err
}

func readws(ctx context.Context, c *websocket.Conn) error {
	ctx, cancel := context.WithTimeout(ctx, time.Second*10)
	defer cancel()

	typ, r, err := c.Reader(ctx)
	if err != nil {
		return err
	}

	w, err := c.Writer(ctx, typ)
	if err != nil {
		return err
	}

	_, err = io.Copy(w, r)
	if err != nil {
		return fmt.Errorf("failed to io.Copy: %w", err)
	}

	err = w.Close()
	return err
}
