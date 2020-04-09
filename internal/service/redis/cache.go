package redis

import (
	"github.com/victorkt/flaggio/internal/service"
)

func shouldCacheEvaluation(req *service.EvaluationRequest) bool {
	return !req.IsDebug()
}
