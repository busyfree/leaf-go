package hooks

import (
	"strings"

	"github.com/busyfree/leaf-go/util/conf"
	"github.com/sirupsen/logrus"
)

type FilterHook struct{}

func (h *FilterHook) Levels() []logrus.Level {
	return logrus.AllLevels
}

func (h *FilterHook) Fire(e *logrus.Entry) error {
	filters := conf.GetStringSlice("LOG_FILTERS")
	ignores := conf.GetStringSlice("LOG_IGNORES")
	if len(e.Message) > 0 {
		if len(ignores) > 0 {
			for _, v := range ignores {
				if strings.Contains(e.Message, v) {
					e.Message = "[IGNORES]"
					return nil
				}
			}
		}
		if len(filters) > 0 {
			for _, v := range filters {
				e.Message = strings.ReplaceAll(e.Message, v, "[FILTER]")
			}
		}
	}
	for k, v := range e.Data {
		if s, ok := v.(string); ok {
			if s == "" {
				delete(e.Data, k)
				continue
			}
			for _, v := range filters {
				if strings.Contains(s, v) {
					e.Data[k] = strings.ReplaceAll(s, v, `[FILTER]`)
				}
			}
		}
	}
	return nil
}
