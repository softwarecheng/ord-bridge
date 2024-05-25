package service

import (
	"fmt"
	"io"
	"os"
	"path/filepath"
	"strings"
	"sync"
	"time"

	"github.com/didip/tollbooth/v7/limiter"
	"github.com/gin-contrib/cors"
	"github.com/gin-contrib/logger"
	"github.com/gin-gonic/gin"
	rotatelogs "github.com/lestrrat-go/file-rotatelogs"
	"github.com/rs/zerolog"
	"github.com/softwarecheng/ord-bridge/indexer"
	"github.com/softwarecheng/ord-bridge/service/bitcoind"
	serverCommon "github.com/softwarecheng/ord-bridge/service/common"
	basic "github.com/softwarecheng/ord-bridge/service/index"
	"github.com/softwarecheng/ord-bridge/service/ord"
	swaggerFiles "github.com/swaggo/files"
	ginSwagger "github.com/swaggo/gin-swagger"
)

type RateLimit struct {
	limit    *limiter.Limiter
	reqCount int
}

type ServiceManager struct {
	basicSrvice  *basic.Service
	ordService   *ord.Service
	btcdService  *bitcoind.Service
	apiConf      *serverCommon.API
	initApiConf  bool
	apiConfMutex sync.Mutex
	apiLimitMap  sync.Map
}

func NewManager(indexer *indexer.Indexer) *ServiceManager {
	return &ServiceManager{
		basicSrvice: basic.NewService(indexer),
		ordService:  ord.NewService(),
		btcdService: bitcoind.NewService(),
	}
}

func (s *ServiceManager) Start(rpcUrl, swaggerHost, swaggerSchemes, rpcProxy, rpcLogFile string, apiConf any) error {
	gin.SetMode(gin.ReleaseMode)
	r := gin.Default()
	var writers []io.Writer
	if rpcLogFile != "" {
		exePath, _ := os.Executable()
		executableName := filepath.Base(exePath)
		if strings.Contains(executableName, "debug") {
			executableName = "debug"
		}
		executableName += ".rpc"
		fileHook, err := rotatelogs.New(
			rpcLogFile+"/"+executableName+".%Y%m%d.log",
			rotatelogs.WithLinkName(rpcLogFile+"/"+executableName+".log"),
			rotatelogs.WithMaxAge(24*time.Hour),
			rotatelogs.WithRotationTime(1*time.Hour),
		)
		if err != nil {
			return fmt.Errorf("failed to create RotateFile hook, error %s", err)
		}
		writers = append(writers, fileHook)
	}
	writers = append(writers, os.Stdout)
	gin.DefaultWriter = io.MultiWriter(writers...)
	r.Use(logger.SetLogger(
		logger.WithLogger(logger.Fn(func(c *gin.Context, l zerolog.Logger) zerolog.Logger {
			if c.Request.Header["Authorization"] == nil {
				return l
			}
			return l.With().
				Str("Authorization", c.Request.Header["Authorization"][0]).
				Logger()
		})),
	))

	config := cors.DefaultConfig()
	config.AllowOrigins = []string{"*"}
	config.OptionsResponseStatusCode = 200
	r.Use(cors.New(config))

	s.InitApiDoc(swaggerHost, swaggerSchemes, rpcProxy)
	r.GET(rpcProxy+"/swagger/*any", ginSwagger.WrapHandler(swaggerFiles.Handler))

	err := s.InitApiConf(apiConf)
	if err != nil {
		return err
	}

	err = s.applyApiConf(r, rpcProxy)
	if err != nil {
		return err
	}

	s.basicSrvice.InitRouter(r, rpcProxy)
	s.ordService.InitRouter(r, rpcProxy)
	s.btcdService.InitRouter(r, rpcProxy)
	parts := strings.Split(rpcUrl, ":")
	if len(parts) < 2 {
		rpcUrl += ":80"
	}
	go r.Run(rpcUrl)
	return nil
}

func (s *ServiceManager) SetLog(rcpLogFile string) {

}
