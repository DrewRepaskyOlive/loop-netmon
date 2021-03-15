package loop

import (
	"time"

	"bitbucket.org/crosschx/loop-netmon/httpmon"

	ldk "github.com/open-olive/loop-development-kit/ldk/go"
)

type Loop struct {
	logger  *ldk.Logger
	checker chan bool
}

const (
	loopName    = "loop-netmon"
	refreshRate = 15 * time.Second
	url         = "https://loop-library-apiqa.oliveai.com/api/health/status"
)

func Serve() error {
	log := ldk.NewLogger(loopName)
	loop, err := NewLoop(log)
	if err != nil {
		return err
	}
	ldk.ServeLoopPlugin(log, loop)
	return nil
}

func NewLoop(logger *ldk.Logger) (*Loop, error) {
	logger.Info("NewLoop called: " + loopName)
	return &Loop{
		logger: logger,
	}, nil
}

func (l *Loop) LoopStart(sidekick ldk.Sidekick) error {
	l.logger.Info("starting " + loopName)

	httpmon.Schedule(l.CheckUp, refreshRate)
	return nil
}

func (l *Loop) LoopStop() error {
	l.logger.Info("stopping " + loopName)
	return nil
}

func (l *Loop) CheckUp() {
	call := httpmon.IsURLUp(url)
	l.logger.Info("checked URL status", "call", call)
}
