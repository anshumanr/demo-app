package api

import "github.com/labstack/echo/v4"

func (s *Server) eventHandler(c echo.Context) error {
	req := new(response)

	if err := c.Bind(req); err != nil {
		s.log.LogError("", "error: %v", err)
		return err
	}

	return nil
}
