// Package fashionmgr configures micro for a fashionmgr environment
package fashionmgr

import (
	"os"
	"strings"
	"sync"

	"github.com/urfave/cli/v2"

	"github.com/micro/micro/v3/profile"
	"github.com/micro/micro/v3/service/metrics"
	"github.com/micro/micro/v3/service/registry"
	microRuntime "github.com/micro/micro/v3/service/runtime"
	"github.com/micro/micro/v3/service/runtime/local"
	"github.com/micro/micro/v3/service/store"
	"github.com/micro/micro/v3/service/store/memory"

	// plugins
	consul "github.com/micro/micro/plugin/consul/v3"
	prometheus "github.com/micro/micro/plugin/prometheus/v3"
)

func init() {
	profile.Register("fashionmgr", fashionmgr)
}

var (
	// hack because setup func is called multiple times due to before func weirdness
	fmOnce sync.Once
)

var fashionmgr = &profile.Profile{
	Name: "fashionmgr",
	Setup: func(ctx *cli.Context) error {
		var retError error
		fmOnce.Do(func() {
			var addresses []string
			addr := os.Getenv("CONSUL_ADDRESSES")
			if addr == "" {
				addresses = []string{"localhost:8500"}
			} else {
				addresses = strings.Split(addr, ",")
			}

			profile.SetupRegistry(consul.NewRegistry(registry.Addrs(addresses...)))

			// Set up a default metrics reporter (being careful not to clash with any that have already been set):
			if !metrics.IsSet() {
				prometheusReporter, err := prometheus.New()
				if err != nil {
					retError = err
					return
				}
				metrics.SetDefaultMetricsReporter(prometheusReporter)
				microRuntime.DefaultRuntime = local.NewRuntime()
			}

			store.DefaultStore = memory.NewStore()
		})
		return retError
	},
}
