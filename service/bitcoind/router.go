package bitcoind

import (
	"github.com/gin-gonic/gin"
)

type Service struct {
}

func NewService() *Service {
	return &Service{}
}

func (s *Service) InitRouter(r *gin.Engine, basePath string) {
	group := r.Group(basePath + "/btc")
	//broadcast raw tx => blockstream api: POST /tx
	group.POST("/tx", s.sendRawTx)
	//get raw block => blockstream api: GET /block/:hash
	group.POST("/getrawblock", s.getRawBlock)
}
