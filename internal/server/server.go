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

func NewServer(s *store.Store, cfg config.Config, logger *slog.Logger) *Server {
	return &Server{
		store:  s,
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

	s.logger.InfoContext(ctx, "server started", "addr", s.config.Addr)

	go func() {
		<-ctx.Done()
		_ = l.Close()
	}()

	for {
		conn, err := l.Accept()
		if err != nil {
			if errors.Is(err, net.ErrClosed) {
				s.logger.InfoContext(ctx, "waiting for active connections to finish...")
				s.wg.Wait()
				s.logger.InfoContext(ctx, "server stopped")
				return nil
			}
			s.logger.ErrorContext(ctx, "accepting connection", "error", err)
			continue
		}

		s.wg.Go(func() {
			s.handleConnection(ctx, conn)
		})
	}
}

func (s *Server) handleConnection(ctx context.Context, conn net.Conn) {
	defer conn.Close()
	r := resp.NewResp(conn)

	for {
		arr, err := r.Read()

		if err != nil {
			if errors.Is(err, io.EOF) {
				s.logger.DebugContext(ctx, "client disconnected")
				return
			}
			s.logger.ErrorContext(ctx, "reading request", "error", err)
			return
		}
		if arr.Type != resp.ArrayType {
			s.logger.WarnContext(ctx, "invalid request type", "expected", "array", "got", arr.Type)
			continue
		}
		if len(arr.Array) == 0 {
			continue
		}

		output := commands.Handle(ctx, strings.ToUpper(arr.Array[0].Value), resp.ToStringSlice(arr.Array[1:]), s.store)
		_, err = io.WriteString(conn, output)
		if err != nil {
			s.logger.ErrorContext(ctx, "writing response", "error", err)
			return
		}
	}
}
