package cli

import (
	domain "github.com/kenriortega/ngonx/internal/proxy/domain"
	handlers "github.com/kenriortega/ngonx/internal/proxy/handlers"
	services "github.com/kenriortega/ngonx/internal/proxy/services"
	"github.com/kenriortega/ngonx/pkg/badgerdb"
	"github.com/kenriortega/ngonx/pkg/genkey"
	"github.com/kenriortega/ngonx/pkg/httpsrv"
	"github.com/kenriortega/ngonx/pkg/logger"
	"github.com/kenriortega/ngonx/pkg/metric"
	"github.com/spf13/cobra"
)

var proxyCmd = &cobra.Command{
	Use:   "proxy",
	Short: "Run ngonx as a reverse proxy",
	Run: func(cmd *cobra.Command, args []string) {
		port, err := cmd.Flags().GetInt(flagPort)
		if err != nil {
			logger.LogError(err.Error())
		}
		generateApiKey, err := cmd.Flags().GetBool(flagGenApiKey)
		if err != nil {
			logger.LogError(err.Error())
		}
		prevKey, err := cmd.Flags().GetString(flagPrevKey)
		if err != nil {
			logger.LogError(err.Error())
		}
		// Exporter Metrics
		go metric.ExposeMetricServer(configFromYaml.ProxyGateway.PortExporterProxy)
		// proxy logic
		engine := configFromYaml.ProxyCache.Engine
		securityType := configFromYaml.ProxySecurity.Type
		key := configFromYaml.ProxyCache.Key + "_" + securityType

		var proxyRepository domain.ProxyRepository
		clientBadger := badgerdb.GetBadgerDB(false)
		proxyRepository = domain.NewProxyRepository(clientBadger)
		h := handlers.ProxyHandler{
			Service: services.NewProxyService(proxyRepository),
		}

		if generateApiKey {
			word := genkey.StringWithCharset()
			apiKey := genkey.ApiKeyGenerator(word)
			_, err := h.Service.SaveSecretKEY(engine, key, apiKey)
			if err != nil {
				logger.LogError("genkey: Failed " + err.Error())
			}
			logger.LogInfo("genkey: Susscefull")
		}
		if prevKey != "" {
			_, err := h.Service.SaveSecretKEY(engine, key, prevKey)
			if err != nil {
				logger.LogError("prevKey: Failed " + err.Error())
			}
			logger.LogInfo("prevKey: Susscefull")
		}

		for _, endpoints := range configFromYaml.ProxyGateway.EnpointsProxy {

			h.ProxyGateway(endpoints, engine, key, securityType)
		}

		if configFromYaml.ProxySSL.Enable {
			portSSL := configFromYaml.ProxyGateway.Port + configFromYaml.ProxySSL.SSLPort
			server := httpsrv.NewServerSSL(
				configFromYaml.ProxyGateway.Host,
				portSSL,
				nil,
			)
			server.StartSSL(
				configFromYaml.ProxySSL.CrtFile,
				configFromYaml.ProxySSL.KeyFile,
			)
		} else {
			port = configFromYaml.ProxyGateway.Port + port
			server := httpsrv.NewServer(
				configFromYaml.ProxyGateway.Host,
				port,
				nil,
			)
			server.Start()
		}
	},
}

func init() {
	proxyCmd.Flags().Int(flagPort, 5000, "Port to serve to run proxy")
	proxyCmd.Flags().Bool(flagGenApiKey, false, "Action for generate hash for protected routes")
	proxyCmd.Flags().String(flagPrevKey, "", "Action for save a previous hash for protected routes to validate JWT")
	rootCmd.AddCommand(proxyCmd)

}
