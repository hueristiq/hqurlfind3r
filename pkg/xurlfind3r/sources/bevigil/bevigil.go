package bevigil

import (
	"encoding/json"
	"fmt"

	"github.com/hueristiq/xurlfind3r/pkg/xurlfind3r/httpclient"
	"github.com/hueristiq/xurlfind3r/pkg/xurlfind3r/sources"
	"github.com/valyala/fasthttp"
)

type getURLsResponse struct {
	Domain string   `json:"domain"`
	URLs   []string `json:"urls"`
}

type Source struct{}

func (source *Source) Run(config *sources.Configuration, domain string) <-chan sources.Result {
	results := make(chan sources.Result)

	go func() {
		defer close(results)

		var err error

		var key string

		key, err = sources.PickRandom(config.Keys.Bevigil)
		if key == "" || err != nil {
			result := sources.Result{
				Type:   sources.Error,
				Source: source.Name(),
				Error:  err,
			}

			results <- result

			return
		}

		getURLsReqHeaders := map[string]string{}

		if len(config.Keys.Bevigil) > 0 {
			getURLsReqHeaders["X-Access-Token"] = key
		}

		getURLsReqURL := fmt.Sprintf("https://osint.bevigil.com/api/%s/urls/", domain)

		var getURLsRes *fasthttp.Response

		getURLsRes, err = httpclient.Get(getURLsReqURL, "", getURLsReqHeaders)
		if err != nil {
			result := sources.Result{
				Type:   sources.Error,
				Source: source.Name(),
				Error:  err,
			}

			results <- result

			return
		}

		var getURLsResData getURLsResponse

		err = json.Unmarshal(getURLsRes.Body(), &getURLsResData)
		if err != nil {
			result := sources.Result{
				Type:   sources.Error,
				Source: source.Name(),
				Error:  err,
			}

			results <- result

			return
		}

		for _, URL := range getURLsResData.URLs {
			if !sources.IsInScope(URL, domain, config.IncludeSubdomains) {
				continue
			}

			result := sources.Result{
				Type:   sources.URL,
				Source: source.Name(),
				Value:  URL,
			}

			results <- result
		}
	}()

	return results
}

func (source *Source) Name() string {
	return "bevigil"
}
