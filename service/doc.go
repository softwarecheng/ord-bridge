package service

import (
	"strings"

	"github.com/softwarecheng/ord-bridge/docs"
)

//	@contact.name	API Support
//	@contact.url	https://ordx.space
//	@contact.email	support@tinyverse.space

// @securityDefinitions.apikey	Bearer
// @in							header
// @name						Authorization
func (s *ServiceManager) InitApiDoc(swaggerHost, schemes, basePath string) {
	docs.SwaggerInfo.Title = "ordx api"
	docs.SwaggerInfo.Version = "v0.1.0"
	schemeList := strings.Split(schemes, ",")
	for _, scheme := range schemeList {
		if scheme == "http" {
			docs.SwaggerInfo.Schemes = append(docs.SwaggerInfo.Schemes, "http")
		} else if scheme == "https" {
			docs.SwaggerInfo.Schemes = append(docs.SwaggerInfo.Schemes, "https")
		}
	}
	if len(docs.SwaggerInfo.Schemes) == 0 {
		docs.SwaggerInfo.Schemes = []string{"http"}
	}

	docs.SwaggerInfo.Description = "ordx api docs for develper"
	docs.SwaggerInfo.Host = swaggerHost
	docs.SwaggerInfo.BasePath = basePath
}
