package log

import (
	"context"
	"github.com/golibs-starter/golib/log/field"
	"github.com/golibs-starter/golib/web/constant"
	webContext "github.com/golibs-starter/golib/web/context"
	"github.com/golibs-starter/golib/web/event"
)

func ContextExtractor(ctx context.Context) []field.Field {
	if requestAttributes := webContext.GetRequestAttributes(ctx); requestAttributes != nil {
		return []field.Field{
			field.Object(constant.ContextReqMeta, NewContextAttributesFromReqAttr(requestAttributes)),
		}
	}
	if eventAttributes := event.GetAttributes(ctx); eventAttributes != nil {
		return []field.Field{
			field.Object(constant.ContextReqMeta, NewContextAttributesFromEventAttrs(eventAttributes)),
		}
	}
	return nil
}
