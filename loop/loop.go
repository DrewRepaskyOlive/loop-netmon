package loop

import ldk "github.com/open-olive/loop-development-kit/ldk/go"

type Loop struct {
	logger *ldk.Logger
}

const (
	loopName = "loop-netmon"
)

func Serve() error {
	log := ldk.NewLogger(loopName)
	loop, err := NewLoop(log)
	if err != nil {
		return err
	}
	// blocking call
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
	l.logger.Info("LoopStart called: " + loopName)
	return nil
}

func (l *Loop) LoopStop() error {
	return nil
}
