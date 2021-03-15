package loop

import (
	"context"
	"fmt"
	"time"

	"bitbucket.org/crosschx/loop-netmon/httpmon"

	ldk "github.com/open-olive/loop-development-kit/ldk/go"
)

type Loop struct {
	ctx      context.Context
	cancel   context.CancelFunc
	logger   *ldk.Logger
	sidekick ldk.Sidekick
	checker  chan bool
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
	l.ctx, l.cancel = context.WithCancel(context.Background())
	l.sidekick = sidekick

	httpmon.Schedule(l.CheckUp, refreshRate)
	return nil
}

func (l *Loop) LoopStop() error {
	l.logger.Info("stopping " + loopName)
	l.cancel()
	if l.checker != nil {
		close(l.checker)
	}
	return nil
}

func (l *Loop) CheckUp() {
	call := httpmon.IsURLUp(url)
	l.logger.Info("checked URL status", "call", call)
	l.SendWhisper(call)
}

func (l *Loop) SendWhisper(call *httpmon.Call) {
	whisper := ldk.WhisperContentMarkdown{
		Label:    "Netmon",
		Markdown: fmt.Sprintf("%v", call),
	}
	go func() {
		err := l.sidekick.Whisper().Markdown(l.ctx, &whisper)
		if err != nil {
			l.logger.Error("failed to emit whisper", "error", err)
		}
	}()
}
