package server

import (
	"context"
	"time"
)

func (s *Server) CloseDbConn() error {
	ctx, cancel := context.WithTimeout(context.Background(), 3*time.Second)
	defer cancel()

	s.conn.Close()
	// Check if the connection is closed
	err := s.conn.Ping(ctx)

	return err
}
