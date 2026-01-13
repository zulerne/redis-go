package server

import (
	"context"
	"errors"
	"io"
	"log/slog"
	"net"
	"strings"
	"sync"

	"github.com/zulerne/redis-go/internal/commands"
	"github.com/zulerne/redis-go/internal/config"
	"github.com/zulerne/redis-go/internal/store"
	"github.com/zulerne/redis-go/pkg/resp"
)

type Server struct {
	store  *store.Store
	config config.Config
	logger *slog.Logger
	wg     sync.WaitGroup
}

func NewServer(store *store.Store, cfg config.Config, logger *slog.Logger) *Server {
	return &Server{
		store:  store,
		config: cfg,
		logger: logger,
	}
}

func (s *Server) Listen(ctx context.Context) error {
	l, err := net.Listen("tcp", s.config.Addr)
	if err != nil {
		return err
	}
	defer l.Close()

	s.logger.Info("server started", "addr", s.config.Addr)

	go func() {
		<-ctx.Done()
		l.Close()
	}()

	for {
		conn, err := l.Accept()
		if err != nil {
			if errors.Is(err, net.ErrClosed) {
				s.logger.Info("waiting for active connections to finish...")
				s.wg.Wait()
				s.logger.Info("server stopped")
				return nil
			}
			s.logger.Error("accepting connection", "error", err)
			continue
		}

		s.wg.Add(1)
		go func() {
			defer s.wg.Done()
			s.handleConnection(conn)
		}()
	}
}

func (s *Server) handleConnection(conn net.Conn) {
	defer conn.Close()
	r := resp.NewResp(conn)

	for {
		arr, err := r.Read()

		if err != nil {
			if errors.Is(err, io.EOF) {
				s.logger.Debug("client disconnected")
				return
			}
			s.logger.Error("reading request", "error", err)
			return
		}
		if arr.Type != resp.ArrayType {
			s.logger.Warn("invalid request type", "expected", "array", "got", arr.Type)
			continue
		}
		if len(arr.Array) == 0 {
			continue
		}

		output := commands.Handle(strings.ToUpper(arr.Array[0].Value), resp.ToStringSlice(arr.Array[1:]), s.store)
		_, err = io.WriteString(conn, output)
		if err != nil {
			s.logger.Error("writing response", "error", err)
			return
		}
	}
}
