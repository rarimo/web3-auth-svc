package cli

import (
	"context"
	"os"
	"os/signal"
	"sync"
	"syscall"

	"github.com/alecthomas/kingpin"
	"github.com/rarimo/web3-auth-svc/internal/config"
	"github.com/rarimo/web3-auth-svc/internal/service"
	"gitlab.com/distributed_lab/kit/kv"
	"gitlab.com/distributed_lab/logan/v3"
)

func Run(args []string) bool {
	log := logan.New()

	defer func() {
		if rvr := recover(); rvr != nil {
			log.WithRecover(rvr).Error("app panicked")
		}
	}()

	cfg := config.New(kv.MustFromEnv())
	log = cfg.Log()

	app := kingpin.New("web3-auth-svc", "")

	runCmd := app.Command("run", "run command")
	serviceCmd := runCmd.Command("service", "run service") // you can insert custom help

	// custom commands go here...

	cmd, err := app.Parse(args[1:])
	if err != nil {
		log.WithError(err).Error("failed to parse arguments")
		return false
	}

	ctx, stop := signal.NotifyContext(context.Background(), os.Interrupt, syscall.SIGTERM)
	defer stop()

	var wg sync.WaitGroup
	run := func(f func(context.Context, config.Config)) {
		wg.Add(1)
		go func() {
			f(ctx, cfg)
			wg.Done()
		}()
	}

	switch cmd {
	case serviceCmd.FullCommand():
		run(service.Run)
	// handle any custom commands here in the same way
	default:
		log.Errorf("unknown command %s", cmd)
		return false
	}

	gracefulStop := make(chan os.Signal, 1)
	signal.Notify(gracefulStop, syscall.SIGTERM, syscall.SIGINT)

	wgch := make(chan struct{})
	go func() {
		wg.Wait()
		close(wgch)
	}()

	select {
	case <-ctx.Done():
		cfg.Log().WithError(ctx.Err()).Info("Interrupt signal received")
		stop()
		<-wgch
	case <-wgch:
		cfg.Log().Warn("all services stopped")
	}

	return true
}
