package task

import (
	"github.com/crdx/mission/args"
	"github.com/crdx/mission/logger"
)

type Action func(args args.Args, logger *logger.Logger) error
