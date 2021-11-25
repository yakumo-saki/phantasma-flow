package server

import (
	"context"
	"errors"
	"net"

	"github.com/rs/zerolog/log"
	"github.com/yakumo-saki/phantasma-flow/procman"
	"github.com/yakumo-saki/phantasma-flow/util"
)

type Server struct {
	procman.ProcmanModuleStruct

	rootCtx    context.Context
	rootCancel context.CancelFunc
	listener   net.Listener

	connections int32 // use with atomic.Add atomic.Load
}

func (sv *Server) IsInitialized() bool {
	return sv.Initialized
}

func (sv *Server) Initialize() error {
	sv.Name = "Server"
	sv.rootCtx, sv.rootCancel = context.WithCancel(context.Background())
	return nil
}

func (sv *Server) GetName() string {
	return sv.Name
}

func (sv *Server) startListen() error {
	// Finally start listening
	// TODO change port by config
	psock, err := net.Listen("tcp", ":5000")
	if err != nil {
		log.Error().Err(err).Msg("Failed listen tcp.")
		return err
	}
	sv.listener = psock
	return nil
}

func (sv *Server) Start(inCh <-chan string, outCh chan<- string) error {
	log := util.GetLoggerWithSource(sv.GetName(), "main")

	sv.FromProcmanCh = inCh
	sv.ToProcmanCh = outCh

	log.Info().Msg("Starting socket server.")

	err := sv.startListen()
	if err != nil {
		return err
	}
	log.Debug().Msg("TCP Socket start")

	go sv.awaitListener(sv.rootCtx)
	log.Info().Msg("Socket server started.")

	for {
		select {
		case v := <-sv.FromProcmanCh:
			log.Debug().Msgf("Got request from procman %s", v)
		case <-sv.rootCtx.Done():
			goto shutdown
		}
	}

shutdown:
	sv.listener.Close()
	log.Debug().Msg("Main thread exited.")
	return nil
}

func (sv *Server) Shutdown() {
	log := util.GetLoggerWithSource(sv.GetName(), "shutdown")
	sv.rootCancel()
	log.Debug().Msg("Shutdown initiated")
}

// Socket handling thread
func (sv *Server) awaitListener(ctx context.Context) {
	log := util.GetLoggerWithSource(sv.GetName(), "awaitListener")
	log.Info().Msg("Start Listener")

	sv.ToProcmanCh <- procman.RES_STARTUP_DONE

	for {
		log.Debug().Msg("Wait for client")

		// Accept() block execution.
		// continue when new client accepted or listener is closed (=shutdown)
		conn, err := sv.listener.Accept()
		if err != nil {
			if errors.Is(err, net.ErrClosed) {
				log.Debug().Err(err).Msg("Stop accept because of shutdown")
				goto shutdown
			} else {
				// Only network error. dont shutdown server
				log.Error().Err(err).Msg("Accept failed. continue")
				continue
			}
		}

		log.Debug().Msg("Accepted new client")
		go sv.dispatch(ctx, conn)

	}

shutdown:

	sv.ToProcmanCh <- procman.RES_SHUTDOWN_DONE
	log.Info().Msg("Socket server stopped.")
}
