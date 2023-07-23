package log

import (
	"context"
	"gitlab.com/golibs-starter/golib/log/field"
	"gitlab.com/golibs-starter/golib/web/constant"
	webContext "gitlab.com/golibs-starter/golib/web/context"
	"gitlab.com/golibs-starter/golib/web/event"
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
