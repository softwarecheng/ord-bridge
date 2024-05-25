package index

import (
	"net/http"
	"strconv"

	"github.com/gin-gonic/gin"
	"github.com/softwarecheng/ord-bridge/indexer"
	"github.com/softwarecheng/ord-bridge/service/common"
)

const QueryParamDefaultLimit = "100"

type Handle struct {
	model *Model
}

func NewHandle(indexer *indexer.Indexer) *Handle {
	return &Handle{
		model: NewModel(indexer),
	}
}

// @Summary Health Check
// @Description Check the health status of the service
// @Tags index
// @Produce json
// @Success 200 {object} HealthResp "Successful response"
// @Router /health [get]
func (s *Handle) getHealth(c *gin.Context) {
	c.JSON(200, &HealthResp{
		Status: "ok",
	})
}

// @Summary index status
// @Description index status
// @Tags index
// @Produce json
// @Security Bearer
// @Success 200 {object} IndexStatusResp "Successful response"
// @Failure 401 "Invalid API Key"
// @Router /ordstatus [get]
func (s *Handle) getOrdStatus(c *gin.Context) {
	indexer := s.model.indexer
	status := indexer.GetStatus()
	genesesAddressCount, err := indexer.GetGenesesAddressCount()
	if err != nil {
		genesesAddressCount = 0
	}
	addressCount, err := indexer.GetAddressCount()
	if err != nil {
		addressCount = 0
	}
	resp := &IndexStatusResp{
		IndexVersion:                  status.IndexVersion,
		DbVersion:                     status.DbVersion,
		SyncInscriptionHeight:         status.SyncInscriptionHeight,
		SyncTransferInscriptionHeight: status.SyncTransferInscriptionHeight,
		BlessedInscriptions:           status.BlessedInscriptions,
		CursedInscriptions:            status.CursedInscriptions,
		AddressCount:                  addressCount,
		GenesesAddressCount:           genesesAddressCount,
	}
	c.JSON(http.StatusOK, resp)
}

// @Summary Get a list of Inscription for a specific address
// @Description Get a list of Inscription for a specific address
// @Tags index.inscription
// @Produce json
// @Param address path string true "Address"
// @Param start query int false "Start index for pagination"
// @Param limit query int false "Limit for pagination"
// @Security Bearer
// @Success 200 {object} OrdInscriptionListResp "Successful response"
// @Failure 401 "Invalid API Key"
// @Router /inscription/address/{address} [get]
func (s *Handle) getAddressInscriptionList(c *gin.Context) {
	resp := &OrdInscriptionListResp{
		BaseResp: common.BaseResp{
			Code: 0,
			Msg:  "ok",
		},
		Data: nil,
	}

	address := c.Param("address")
	start, err := strconv.Atoi(c.DefaultQuery("start", "0"))
	if err != nil {
		start = 0
	}
	limit, err := strconv.Atoi(c.DefaultQuery("limit", QueryParamDefaultLimit))
	if err != nil {
		limit = 100
	}

	ordInscriptionList, total, err := s.model.GetInscriptionAssetList(address, false, start, limit)
	if err != nil {
		resp.Code = -1
		resp.Msg = err.Error()
		c.JSON(http.StatusOK, resp)
		return
	}

	resp.Data = &OrdInscriptionListData{
		ListResp: common.ListResp{
			Total: total,
			Start: start,
		},
		Detail: ordInscriptionList,
	}
	c.JSON(http.StatusOK, resp)
}

// @Summary Get a list of Inscription for a specific geneses address
// @Description Get a list of Inscription for a specific geneses address
// @Tags index.inscription
// @Produce json
// @Param address path string true "geneses address"
// @Param start query int false "Start index for pagination"
// @Param limit query int false "Limit for pagination"
// @Security Bearer
// @Success 200 {object} OrdInscriptionListResp "Successful response"
// @Failure 401 "Invalid API Key"
// @Router /inscription/genesesaddress/{address} [get]
func (s *Handle) getGenesesAddressInscriptionList(c *gin.Context) {
	resp := &OrdInscriptionListResp{
		BaseResp: common.BaseResp{
			Code: 0,
			Msg:  "ok",
		},
		Data: nil,
	}

	address := c.Param("address")
	start, err := strconv.Atoi(c.DefaultQuery("start", "0"))
	if err != nil {
		start = 0
	}
	limit, err := strconv.Atoi(c.DefaultQuery("limit", QueryParamDefaultLimit))
	if err != nil {
		limit = 100
	}

	ordInscriptionList, total, err := s.model.GetInscriptionAssetList(address, true, start, limit)
	if err != nil {
		resp.Code = -1
		resp.Msg = err.Error()
		c.JSON(http.StatusOK, resp)
		return
	}

	resp.Data = &OrdInscriptionListData{
		ListResp: common.ListResp{
			Total: total,
			Start: start,
		},
		Detail: ordInscriptionList,
	}
	c.JSON(http.StatusOK, resp)
}

// @Summary Get a list of Inscription with sat
// @Description Get a list of Inscription with sat
// @Tags index.inscription
// @Produce json
// @Param sat path string true "sat"
// @Param start query int false "Start index for pagination"
// @Param limit query int false "Limit for pagination"
// @Security Bearer
// @Success 200 {object} OrdInscriptionListResp "Successful response"
// @Failure 401 "Invalid API Key"
// @Router /inscription/sat/{sat} [get]
func (s *Handle) getInscriptionListWithSat(c *gin.Context) {
	resp := &OrdInscriptionListResp{
		BaseResp: common.BaseResp{
			Code: 0,
			Msg:  "ok",
		},
		Data: nil,
	}

	query := c.Param("sat")
	start, err := strconv.Atoi(c.DefaultQuery("start", "0"))
	if err != nil {
		start = 0
	}
	limit, err := strconv.Atoi(c.DefaultQuery("limit", QueryParamDefaultLimit))
	if err != nil {
		limit = 100
	}
	inscriptionList, total, err := s.model.GetOrdInscriptionWithSat(query, start, limit)
	if err != nil {
		resp.Code = -1
		resp.Msg = err.Error()
		c.JSON(http.StatusOK, resp)
		return
	}

	resp.Data = &OrdInscriptionListData{
		ListResp: common.ListResp{
			Total: total,
			Start: start,
		},
		Detail: inscriptionList,
	}
	c.JSON(http.StatusOK, resp)
}

// @Summary Get inscription with inscription id
// @Description Get inscription with inscription id
// @Tags index.inscription
// @Produce json
// @Param id path string true "inscription id"
// @Security Bearer
// @Success 200 {object} OrdInscriptionResp "Successful response"
// @Failure 401 "Invalid API Key"
// @Router /inscription/id/{id} [get]
func (s *Handle) getInscriptionWithId(c *gin.Context) {
	resp := &OrdInscriptionResp{
		BaseResp: common.BaseResp{
			Code: 0,
			Msg:  "ok",
		},
		Data: nil,
	}

	query := c.Param("id")
	inscription, err := s.model.GetOrdInscriptionWithId(query)
	if err != nil {
		resp.Code = -1
		resp.Msg = err.Error()
		c.JSON(http.StatusOK, resp)
		return
	}
	resp.Data = inscription
	c.JSON(http.StatusOK, resp)
}

// @Summary Get inscription with inscription number
// @Description Get inscription with inscription number
// @Tags index.inscription
// @Produce json
// @Param number path string true "inscription Number"
// @Security Bearer
// @Success 200 {object} OrdInscriptionResp "Successful response"
// @Failure 401 "Invalid API Key"
// @Router /inscription/number/{number} [get]
func (s *Handle) getInscriptionWithNumber(c *gin.Context) {
	resp := &OrdInscriptionResp{
		BaseResp: common.BaseResp{
			Code: 0,
			Msg:  "ok",
		},
		Data: nil,
	}

	query := c.Param("number")
	inscription, err := s.model.GetOrdInscriptionWithNumer(query)
	if err != nil {
		resp.Code = -1
		resp.Msg = err.Error()
		c.JSON(http.StatusOK, resp)
		return
	}
	resp.Data = inscription
	c.JSON(http.StatusOK, resp)
}
