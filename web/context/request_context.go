package context

import (
	"context"
	"gitlab.id.vin/vincart/golib/web/constants"
)

func GetRequestAttributes(ctx context.Context) *RequestAttributes {
	reqAttrCtxValue := ctx.Value(constants.ContextReqAttribute)
	if reqAttrCtxValue == nil {
		return nil
	}
	requestAttributes, ok := reqAttrCtxValue.(*RequestAttributes)
	if !ok {
		return nil
	}
	return requestAttributes
}
