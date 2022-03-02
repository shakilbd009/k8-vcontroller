package main

import (
	"fmt"
	"net/http"
	"os"
	"time"

	"github.com/spf13/pflag"
	"k8s.io/apiserver/pkg/server"
	"k8s.io/apiserver/pkg/server/options"
	"k8s.io/component-base/cli/globalflag"
)

type Options struct {
	SecureServingOptions options.SecureServingOptions
}

type Config struct {
	SecureServeingInfo *server.SecureServingInfo
}

const (
	valKon = "val-kontroller"
)

func (o *Options) AddFlagSet(fs *pflag.FlagSet) {
	o.SecureServingOptions.AddFlags(fs)
}

func (o *Options) Config() *Config {
	if err := o.SecureServingOptions.MaybeDefaultWithSelfSignedCerts("0.0.0.0", nil, nil); err != nil {
		panic(err)
	}
	c := Config{}
	o.SecureServingOptions.ApplyTo(&c.SecureServeingInfo)
	return &c
}

func main() {
	options := NewDefaultOptions()
	fs := pflag.NewFlagSet(valKon, pflag.ExitOnError)
	globalflag.AddGlobalFlags(fs, valKon)
	options.AddFlagSet(fs)
	if err := fs.Parse(os.Args[:]); err != nil {
		panic(err)
	}
	c := options.Config()
	stopCH := server.SetupSignalHandler()
	mux := http.NewServeMux()
	mux.Handle("/", http.HandlerFunc(ServerKlusterValidation))
	ch, err := c.SecureServeingInfo.Serve(mux, time.Second*30, stopCH)
	if err != nil {
		panic(err)
	} else {
		<-ch
	}

}

func ServerKlusterValidation(w http.ResponseWriter, r *http.Request) {
	fmt.Println("valconroller was called")
}

func NewDefaultOptions() *Options {
	o := &Options{
		SecureServingOptions: *options.NewSecureServingOptions(),
	}
	o.SecureServingOptions.BindPort = 8443
	o.SecureServingOptions.ServerCert.PairName = valKon
	return o
}
