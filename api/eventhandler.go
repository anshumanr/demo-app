package api

import (
	"fmt"

	"github.com/labstack/echo/v4"
)

func (s *Server) eventHandler(c echo.Context) error {
	req := cbResponse{}

	if err := c.Bind(&req); err != nil {
		s.log.LogError("", "error: %v", err)
		return err
	}

	buf := fmt.Sprintf("%+v", req)

	go func() {
		fmt.Println(buf)
		select {
		case evChanW <- buf:
		default:
			fmt.Println("no event sent")
		}
	}()

	return c.String(200, "")
}
