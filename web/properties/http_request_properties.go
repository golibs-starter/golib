package properties

import (
	"fmt"
	"github.com/golibs-starter/golib/config"
	"regexp"
)

func NewHttpRequestLogProperties(loader config.Loader) (*HttpRequestLogProperties, error) {
	props := HttpRequestLogProperties{}
	err := loader.Bind(&props)
	return &props, err
}

type HttpRequestLogProperties struct {
	Disabled               bool
	PredefinedDisabledUrls []*UrlMatching `default:"[{\"UrlPattern\":\"^/actuator/.*\"}]"`
	DisabledUrls           []*UrlMatching
	allDisabledUrls        []*UrlMatching
}

func (h *HttpRequestLogProperties) AllDisabledUrls() []*UrlMatching {
	return h.allDisabledUrls
}

func (h *HttpRequestLogProperties) Prefix() string {
	return "app.httpRequest.logging"
}

func (h *HttpRequestLogProperties) PostBinding() error {
	if len(h.PredefinedDisabledUrls) > 0 {
		for _, disabledUrl := range h.PredefinedDisabledUrls {
			regex, err := regexp.Compile(disabledUrl.UrlPattern)
			if err != nil {
				return fmt.Errorf("url pattern [%s] is not valid in regex format, error [%v]",
					disabledUrl.UrlPattern, err)
			}
			disabledUrl.urlRegexp = regex
		}
	}
	if len(h.DisabledUrls) > 0 {
		for _, disabledUrl := range h.DisabledUrls {
			regex, err := regexp.Compile(disabledUrl.UrlPattern)
			if err != nil {
				return fmt.Errorf("url pattern [%s] is not valid in regex format, error [%v]",
					disabledUrl.UrlPattern, err)
			}
			disabledUrl.urlRegexp = regex
		}
	}
	h.allDisabledUrls = append(h.PredefinedDisabledUrls, h.DisabledUrls...)
	return nil
}

type UrlMatching struct {
	Method     string
	UrlPattern string
	urlRegexp  *regexp.Regexp
}

func (u UrlMatching) UrlRegexp() *regexp.Regexp {
	return u.urlRegexp
}
