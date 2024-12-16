package app

import (
	"context"
	"errors"
	"net/http"
	"os"
	"os/signal"
	"strconv"
	"sync"
	"time"

	"github.com/go-chi/chi/v5"
	"github.com/jackc/pgx/v5/pgxpool"
	"go.uber.org/zap"
	"go.uber.org/zap/zapcore"

	"gihub.com/bongerka/sberPDRIS/internal/config"
	"gihub.com/bongerka/sberPDRIS/internal/utils"
)

type App struct {
	cfg *config.Config

	conn *pgxpool.Pool

	wg        *sync.WaitGroup
	ctx       context.Context
	cancelCtx context.CancelFunc

	pdrisProvider *serviceProvider

	server *http.Server
}

func NewApp() (*App, error) {
	ctx, cancel := context.WithCancel(context.Background())

	a := &App{
		wg:        &sync.WaitGroup{},
		ctx:       ctx,
		cancelCtx: cancel,
	}

	err := a.initDeps(ctx)
	if err != nil {
		return nil, err
	}

	return a, nil
}

func (a *App) Run() {
	defer a.conn.Close()
	defer a.cancelCtx()

	ctx, cancel := signal.NotifyContext(a.ctx, os.Interrupt)
	defer cancel()

	go func() {
		if err := a.server.ListenAndServe(); err != nil {
			if errors.Is(err, http.ErrServerClosed) {
				zap.L().Info("successfully shutdown")
				return
			}
			zap.L().Fatal("shutdown with error", zap.Error(err))
		}
	}()

	zap.L().Sugar().Infof("server successfully started on %s", a.cfg.Server.HTTPAddr)

	select {
	case <-ctx.Done():
		if err := a.server.Shutdown(ctx); err != nil {
			zap.L().Error("error when shutdown", zap.Error(ctx.Err()))
		}
	}
}

func (a *App) initDeps(ctx context.Context) error {
	inits := []func(context.Context) error{
		a.initConfig,
		a.initLogger,
		a.initDB,
		a.initService,
		a.initHTTPServer,
	}

	for _, f := range inits {
		err := f(ctx)
		if err != nil {
			return err
		}
	}

	return nil
}

func (a *App) initConfig(_ context.Context) error {
	cfg, err := config.NewConfig[config.Config]()
	if err != nil {
		return err
	}

	a.cfg = cfg

	return nil
}

func (a *App) initLogger(_ context.Context) error {
	if a.cfg.Server.IsDev() {
		minLevel := zapcore.DebugLevel

		developmentCfg := zap.NewDevelopmentEncoderConfig()
		developmentCfg.EncodeLevel = zapcore.CapitalColorLevelEncoder
		developmentCfg.EncodeTime = zapcore.TimeEncoderOfLayout(time.StampMilli)
		developmentCfg.EncodeCaller = func(caller zapcore.EntryCaller, encoder zapcore.PrimitiveArrayEncoder) {
			rel, err := utils.GetRelPath(caller.String())
			if err != nil {
				encoder.AppendString(caller.String())
				return
			}

			encoder.AppendString(rel)
		}

		consoleEncoder := zapcore.NewConsoleEncoder(developmentCfg)
		core := zapcore.NewCore(consoleEncoder, zapcore.AddSync(os.Stdout), minLevel)

		lg := zap.New(core, zap.AddCaller(), zap.AddStacktrace(zapcore.ErrorLevel))
		zap.ReplaceGlobals(lg)
	}

	return nil
}

func (a *App) initDB(_ context.Context) error {
	conn, err := pgxpool.New(a.ctx, "user=postgres password=password host=pg_pdris4 port=5432 dbname=postgres")
	if err != nil {
		return err
	}

	a.conn = conn
	return nil
}

func (a *App) initService(_ context.Context) error {
	a.pdrisProvider = NewServiceProvider(a.conn)

	return nil
}

func (a *App) initHTTPServer(_ context.Context) error {
	mux := chi.NewMux()

	a.server = &http.Server{
		Addr:              a.cfg.Server.HTTPAddr,
		Handler:           mux,
		ReadTimeout:       10 * time.Second,
		WriteTimeout:      10 * time.Second,
		ReadHeaderTimeout: 5 * time.Second,
		IdleTimeout:       10 * time.Minute,
	}

	// unclear and unscalable code architecture but good enough for hw
	mux.Get("/update/{value}", func(writer http.ResponseWriter, request *http.Request) {
		strValue := chi.URLParam(request, "value")
		value, err := strconv.Atoi(strValue)
		if err != nil {
			writer.Write([]byte("wrong input"))
			writer.WriteHeader(http.StatusBadRequest)
			return
		}

		if err := a.pdrisProvider.Service().UpdateValue(a.ctx, value); err != nil {
			writer.WriteHeader(http.StatusInternalServerError)
		}

		writer.Write([]byte("successfully updated"))
		writer.WriteHeader(http.StatusOK)
	})

	mux.Get("/", func(writer http.ResponseWriter, request *http.Request) {
		value, err := a.pdrisProvider.Service().GetValue(a.ctx)
		if err != nil {
			writer.WriteHeader(http.StatusInternalServerError)
		}

		writer.Write([]byte(strconv.Itoa(value)))
		writer.WriteHeader(http.StatusOK)
	})

	return nil
}
