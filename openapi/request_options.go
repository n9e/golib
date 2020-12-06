package openapi

import (
	"context"
	"crypto/tls"
	"fmt"
	"net/http"

	"github.com/yubo/golib/tlsutil"
	"github.com/yubo/golib/urlutil"
	"github.com/yubo/golib/util"
)

type RequestOptions struct {
	http.Client
	Url                string // https://example.com/api/v{version}/{model}/{object}?type=vm
	Method             string
	User               *string
	Pwd                *string
	Bearer             *string
	ApiKey             *string
	InputFile          *string // Priority InputFile > InputContent > Input
	InputContent       []byte
	Input              interface{}
	OutputFile         *string // Priority OutputFile > Output
	Output             interface{}
	Mime               string
	Ctx                context.Context
	header             http.Header
	CertFile           string
	KeyFile            string
	CaFile             string
	InsecureSkipVerify bool
}

func (p RequestOptions) String() string {
	return util.Prettify(p)
}

func (p RequestOptions) Transport() (tr *http.Transport, err error) {
	tr = &http.Transport{
		DisableCompression: true,
		Proxy:              http.ProxyFromEnvironment,
	}
	if (p.CertFile != "" && p.KeyFile != "") || p.CaFile != "" {
		tlsConf, err := p.TLSClientConfig()
		if err != nil {
			return nil, fmt.Errorf("can't create TLS config: %s", err.Error())
		}
		tr.TLSClientConfig = tlsConf
	}
	return tr, nil
}

func (p RequestOptions) TLSClientConfig() (*tls.Config, error) {
	serverName, err := urlutil.ExtractHostname(p.Url)
	if err != nil {
		return nil, err
	}

	return tlsutil.ClientConfig(tlsutil.Options{
		CaCertFile:         p.CaFile,
		KeyFile:            p.KeyFile,
		CertFile:           p.CertFile,
		InsecureSkipVerify: p.InsecureSkipVerify,
		ServerName:         serverName,
	})

}

type RequestOption interface {
	apply(*RequestOptions)
}

type funcRequestOption struct {
	f func(*RequestOptions)
}

func (p *funcRequestOption) apply(opt *RequestOptions) {
	p.f(opt)
}

func newFuncRequestOption(f func(*RequestOptions)) *funcRequestOption {
	return &funcRequestOption{
		f: f,
	}
}

func WithUrl(url string) RequestOption {
	return newFuncRequestOption(func(o *RequestOptions) {
		o.Url = url
	})
}

func WithMethod(method string) RequestOption {
	return newFuncRequestOption(func(o *RequestOptions) {
		o.Method = method
	})
}

func WithUser(user, pwd string) RequestOption {
	return newFuncRequestOption(func(o *RequestOptions) {
		o.User = &user
		o.Pwd = &pwd
	})
}

func WithBearer(bearer string) RequestOption {
	return newFuncRequestOption(func(o *RequestOptions) {
		o.Bearer = &bearer
	})
}

func WithApiKey(apiKey string) RequestOption {
	return newFuncRequestOption(func(o *RequestOptions) {
		o.ApiKey = &apiKey
	})
}

func WithInputFile(filePath string) RequestOption {
	return newFuncRequestOption(func(o *RequestOptions) {
		o.InputFile = &filePath
	})
}

func WithInputContent(body []byte) RequestOption {
	return newFuncRequestOption(func(o *RequestOptions) {
		o.InputContent = body
	})
}

func WithInput(input interface{}) RequestOption {
	return newFuncRequestOption(func(o *RequestOptions) {
		o.Input = input
	})
}

func WithOutputFile(filePath string) RequestOption {
	return newFuncRequestOption(func(o *RequestOptions) {
		o.OutputFile = &filePath
	})
}

func WithOutput(output interface{}) RequestOption {
	return newFuncRequestOption(func(o *RequestOptions) {
		o.Output = output
	})
}

func WithMime(mime string) RequestOption {
	return newFuncRequestOption(func(o *RequestOptions) {
		o.Mime = mime
	})
}

func WithCtx(ctx context.Context) RequestOption {
	return newFuncRequestOption(func(o *RequestOptions) {
		o.Ctx = ctx
	})
}

func WithHeader(k, v string) RequestOption {
	return newFuncRequestOption(func(o *RequestOptions) {
		o.header.Set(k, v)
	})
}

func WithTLSConfig(certFile, keyFile, caFile string) RequestOption {
	return newFuncRequestOption(func(o *RequestOptions) {
		o.CertFile = certFile
		o.KeyFile = keyFile
		o.CaFile = caFile
	})
}

func InsecureSkipVerify() RequestOption {
	return newFuncRequestOption(func(o *RequestOptions) {
		o.InsecureSkipVerify = true
	})
}
