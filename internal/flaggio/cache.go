package flaggio

import (
	"fmt"
	"strings"
)

const (
	namespace         = "flaggio"
	flagNamespace     = "flag"
	segmentNamespace  = "segment"
	evaluateNamespace = "eval"
)

func cacheKey(model string, parts ...string) string {
	var sb strings.Builder
	sb.WriteString(fmt.Sprintf("%s:%s", namespace, model))
	for _, part := range parts {
		sb.WriteString(fmt.Sprintf(":%s", part))
	}
	return sb.String()
}

func FlagCacheKey(parts ...string) string {
	return cacheKey(flagNamespace, parts...)
}

func SegmentCacheKey(parts ...string) string {
	return cacheKey(segmentNamespace, parts...)
}

func EvalCacheKey(parts ...string) string {
	return cacheKey(evaluateNamespace, parts...)
}
