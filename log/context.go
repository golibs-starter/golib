package log

import (
	"context"
	"gitlab.com/golibs-starter/golib/log/field"
)

type ContextExtractor func(ctx context.Context) []field.Field

type ContextExtractors []ContextExtractor

func (c ContextExtractors) IsExtractable() bool {
	return len(c) > 0
}

func (c ContextExtractors) Extract(ctx context.Context) []field.Field {
	fields := make([]field.Field, 0)
	for _, extractor := range c {
		fields = append(fields, extractor(ctx)...)
	}
	return fields
}
