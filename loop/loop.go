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

var currentStatus = true

func (l *Loop) CheckUp() {
	call := httpmon.IsURLUp(url)
	l.logger.Info("checked URL status", "call", call)
	l.SendWhisperCall(call)
}

func (l *Loop) SendWhisperCall(call *httpmon.Call) {
	if currentStatus == call.Success {
		return
	}
	currentStatus = call.Success
	label := "Network Monitor: call failed"
	if call.Success {
		label = "Network Monitor: call succeeded"
	}
	markdown := fmt.Sprintf(`Status change for %s
%s
`, call.URL, call.Error)

	l.SendWhisper(label, markdown)
}

func (l *Loop) SendWhisper(label string, markdown string) {
	whisper := ldk.WhisperContentMarkdown{
		Label:    label,
		Markdown: markdown,
	}
	go func() {
		err := l.sidekick.Whisper().Markdown(l.ctx, &whisper)
		if err != nil {
			l.logger.Error("failed to emit whisper", "error", err)
		}
	}()
}
