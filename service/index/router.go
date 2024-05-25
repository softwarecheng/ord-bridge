package index

import (
	"github.com/gin-gonic/gin"
	"github.com/softwarecheng/ord-bridge/indexer"
)

type Service struct {
	handle *Handle
}

func NewService(indexer *indexer.Indexer) *Service {
	return &Service{
		handle: NewHandle(indexer),
	}
}

func (s *Service) InitRouter(r *gin.Engine, basePath string) {
	r.GET(basePath+"/health", s.handle.getHealth)
	r.GET(basePath+"/ordstatus", s.handle.getOrdStatus)
	group := r.Group(basePath + "/inscription")
	// the address of the given inscription
	group.GET("/address/:address", s.handle.getAddressInscriptionList)
	// the geneses address of the given inscription
	group.GET("/genesesaddress/:address", s.handle.getGenesesAddressInscriptionList)
	// the content of the given inscription id
	group.GET("/id/:id", s.handle.getInscriptionWithId)
	// the content of the given inscription number
	group.GET("/number/:number", s.handle.getInscriptionWithNumber)
	// the content of the given inscription sat
	group.GET("/sat/:sat", s.handle.getInscriptionListWithSat)
}
