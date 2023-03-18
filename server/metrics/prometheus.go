package metrics

import (
	"context"
	"net/http"
	"runtime"
	"time"

	"github.com/prometheus/client_golang/prometheus"
	"github.com/prometheus/client_golang/prometheus/promhttp"
	"github.com/yakumo-saki/phantasma-flow/procman"
	"github.com/yakumo-saki/phantasma-flow/util"
)

type PrometeusExporterModule struct {
	procman.ProcmanModuleStruct

	serv *http.Server
}

// returns this instance is initialized or not.
// When procman.Add, Procman calls Initialize() if not initialized.
func (m *PrometeusExporterModule) IsInitialized() bool {
	return m.Initialized
}

// initialize this instance.
// Between Initialize and Start, no shutdown is called when error occures.
// so, dont initialize something needs shutdown sequence.
func (m *PrometeusExporterModule) Initialize() error {
	m.Initialized = true
	m.RootCtx, m.RootCancel = context.WithCancel(context.Background())
	return nil
}

func (m *PrometeusExporterModule) GetName() string {
	return "PrometeusExporterModule"
}

// lets roll! Do not forget to save procmanCh from parameter.
func (m *PrometeusExporterModule) Start(inCh <-chan string, outCh chan<- string) error {
	m.FromProcmanCh = inCh
	m.ToProcmanCh = outCh
	log := util.GetLoggerWithSource(m.GetName(), "main")

	log.Info().Msgf("Starting %s.", m.GetName())

	go m.startServer()
	m.ToProcmanCh <- procman.RES_STARTUP_DONE

	// wait for other message from Procman
	for {
		select {
		case v := <-m.FromProcmanCh:
			log.Debug().Msgf("Got request %s", v)
		case <-m.RootCtx.Done():
			goto shutdown
		}
	}

shutdown:

	log.Info().Msgf("%s Stopped.", m.GetName())
	m.ToProcmanCh <- procman.RES_SHUTDOWN_DONE
	return nil
}

func (m *PrometeusExporterModule) startServer() {
	log := util.GetLoggerWithSource(m.GetName(), "startServer")
	m.addPrometheusMetrics()

	mux := http.NewServeMux()
	mux.Handle("/metrics", promhttp.Handler())
	server := http.Server{Addr: ":5001", Handler: mux}

	m.serv = &server

	log.Info().Msg("Prometeus exporter started.")
	err := server.ListenAndServe() // block until server.shutdown() called
	if err != http.ErrServerClosed {
		log.Err(err).Msg("Prometeus exporter startup failed.")
	}
}

func (m *PrometeusExporterModule) addPrometheusMetrics() {
	// Subsystem: "runtime",
	// Name:      "goroutines_count",
	// Help:      "Number of goroutines that currently exist.",

	// In prometheus keys "-" is prohibited. use "_" instead
	err := prometheus.Register(
		prometheus.NewGaugeFunc(
			prometheus.GaugeOpts{
				Subsystem: "phantasmaflow",
				Name:      "messageHub_queue_length",
				Help:      "Number of messages queued in messageHub.",
			},
			func() float64 {
				return float64(runtime.NumGoroutine())
			},
		),
	)
	if err != nil {
		panic("Failed to register prometheus metrics (messagequeue)")
	}
}

func (sv *PrometeusExporterModule) Shutdown() {
	// When shutdown initiated, procman calls this function.
	// All modules must send SHUTDOWN_DONE to procman before timeout.
	// Otherwise procman is not stop or force shutdown.

	log := util.GetLogger()
	log.Debug().Msg("Shutdown initiated")

	ctx, cancel := context.WithTimeout(context.Background(), 3*time.Second)
	defer cancel()

	sv.serv.Shutdown(ctx) // Its no problem calling shutdown. if server startup failed
	log.Debug().Msg("Metrics server shutdown complete.")
	<-ctx.Done()

	sv.RootCancel()
}
