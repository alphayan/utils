package xcjHttpClient

import "github.com/ddliu/go-httpclient"

//type XcjHttpClient struct {
//	httpclient.HttpClient
//}
var (
	httpClient *httpclient.HttpClient
)

func newHttpClient() {
	httpClient = httpclient.Defaults(httpclient.Map{
		httpclient.OPT_USERAGENT: "my awsome httpclient",
		"Accept-Language":        "en-us",
	})
}
func GetHttpClient() *httpclient.HttpClient {
	if httpClient == nil {
		newHttpClient()
	}
	return httpClient
}
