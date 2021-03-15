package loop

import (
	"context"
	"fmt"
	"time"

	"bitbucket.org/crosschx/loop-netmon/httpmon"

	ldk "github.com/open-olive/loop-development-kit/ldk/go"
)

type Status struct {
	URL     string
	Success bool
}

type Loop struct {
	ctx      context.Context
	cancel   context.CancelFunc
	logger   *ldk.Logger
	sidekick ldk.Sidekick
	checker  chan bool
	statuses []Status
}

const (
	loopName    = "loop-netmon"
	refreshRate = 15 * time.Second
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

	l.SendWhisper("Network Monitor Loop Started", "Copy any URL to start monitoring it")

	l.statuses = []Status{
		{URL: "http://localhost:8000", Success: true},
	}
	l.checker = httpmon.Schedule(l.CheckUp, refreshRate)
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
	for _, status := range l.statuses {
		call := httpmon.IsURLUp(status.URL)
		l.logger.Info("checked URL status", "call", call)
		if status.Success == call.Success {
			continue
		}
		status.Success = call.Success
		l.SendWhisperCall(call)
	}
}

func (l *Loop) SendWhisperCall(call *httpmon.Call) {
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
