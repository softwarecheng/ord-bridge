package ord

import (
	"net/http"
	"strconv"

	"github.com/gin-gonic/gin"
	"github.com/softwarecheng/ord-bridge/rpc/ord"
	"github.com/softwarecheng/ord-bridge/service/common"
)

// @Summary get ordinal status
// @Description get ordinal status
// @Tags ordx.ord
// @Produce json
// @Security Bearer
// @Success 200 {object} OrdStatusResp "Successful response"
// @Failure 401 "Invalid API Key"
// @Router /ord/status [get]
func (s *Service) getOrdStatus(c *gin.Context) {
	resp := &OrdStatusResp{
		BaseResp: common.BaseResp{
			Code: 0,
			Msg:  "ok",
		},
		IsOrdx: true,
		Data:   nil,
	}
	ordinalStatusResp, err := ord.OrdRpc.GetStatus()
	if err != nil {
		resp.Code = -1
		resp.Msg = err.Error()
		c.JSON(http.StatusOK, resp)
		return
	}
	resp.Data = ordinalStatusResp
	c.JSON(http.StatusOK, resp)
}

// @Summary get ordinal content
// @Description get ordinal content
// @Tags ordx.ord
// @Produce image/*
// @Param inscriptionid path string true "inscription ID" example:"5ea0f691bc818afd978e26d308f56ac2bdadfe6c403f53f73872aaa0cef55fd1i0"
// @Security Bearer
// @Success 200 {object} []byte "Successful response"
// @Failure 401 "Invalid API Key"
// @Router /ord/content/{inscriptionid} [get]
func (s *Service) getInscriptionContent(c *gin.Context) {
	resp := &InscriptionContentResp{
		BaseResp: common.BaseResp{
			Code: 0,
			Msg:  "ok",
		},
		Data: nil,
	}
	inscriptionContent, header, err := ord.OrdRpc.GetInscriptionContent(c.Param("inscriptionid"))
	if err != nil {
		resp.Code = -1
		resp.Msg = err.Error()
		c.JSON(http.StatusOK, err.Error())
		return
	}
	for k, v := range header {
		c.Header(k, v[0])
	}
	resp.Data = inscriptionContent
	c.Data(http.StatusOK, header["Content-Type"][0], resp.Data)
}

// @Summary get ordinal preview
// @Description get ordinal preview
// @Tags ordx.ord
// @Produce image/*
// @Param inscriptionid path string true "inscription ID" example:"5ea0f691bc818afd978e26d308f56ac2bdadfe6c403f53f73872aaa0cef55fd1i0"
// @Security Bearer
// @Success 200 {object} []byte "Successful response"
// @Failure 401 "Invalid API Key"
// @Router /ord/preview/{inscriptionid} [get]
func (s *Service) getInscriptionPreview(c *gin.Context) {
	resp := &InscriptionContentResp{
		BaseResp: common.BaseResp{
			Code: 0,
			Msg:  "ok",
		},
		Data: nil,
	}
	inscriptionContent, header, err := ord.OrdRpc.GetInscriptionContent(c.Param("inscriptionid"))
	if err != nil {
		resp.Code = -1
		resp.Msg = err.Error()
		c.JSON(http.StatusOK, err.Error())
		return
	}
	for k, v := range header {
		c.Header(k, v[0])
	}
	resp.Data = inscriptionContent
	c.Data(http.StatusOK, header["Content-Type"][0], resp.Data)
}

// @Summary ordinal recursive endpoint for get block hash
// @Description ordinal recursive endpoint for get block hash
// @Tags ordx.ord.r
// @Produce json
// @Param height path uint64 true "height"
// @Security Bearer
// @Success 200 {object} RBlockHashResp "Successful response"
// @Failure 401 "Invalid API Key"
// @Router /ord/r/blockhash/{height} [get]
func (s *Service) getRBlockHash(c *gin.Context) {
	resp := &RBlockHashResp{
		BaseResp: common.BaseResp{
			Code: 0,
			Msg:  "ok",
		},
		Data: "",
	}
	height, err := strconv.ParseUint(c.Param("height"), 10, 64)
	if err != nil {
		resp.Code = -1
		resp.Msg = err.Error()
		c.JSON(http.StatusOK, resp)
		return
	}
	blockHash, err := ord.OrdRpc.GetRBlockHash(height)
	if err != nil {
		resp.Code = -1
		resp.Msg = err.Error()
		c.JSON(http.StatusOK, resp)
		return
	}
	resp.Data = blockHash
	c.JSON(http.StatusOK, resp)
}

// @Summary ordinal recursive endpoint for get lastest block hash
// @Description ordinal recursive endpoint for get lastest block hash
// @Tags ordx.ord.r
// @Produce json
// @Security Bearer
// @Success 200 {object} RBlockHashResp "Successful response"
// @Failure 401 "Invalid API Key"
// @Router /ord/r/blockhash [get]
func (s *Service) getRLastestBlockHash(c *gin.Context) {
	resp := &RBlockHashResp{
		BaseResp: common.BaseResp{
			Code: 0,
			Msg:  "ok",
		},
		Data: "",
	}
	blockHash, err := ord.OrdRpc.GetRLastestBlockHash()
	if err != nil {
		resp.Code = -1
		resp.Msg = err.Error()
		c.JSON(http.StatusOK, resp)
		return
	}
	resp.Data = blockHash
	c.JSON(http.StatusOK, resp)
}

// @Summary ordinal recursive endpoint for get lastest block height
// @Description ordinal recursive endpoint for get lastest block height
// @Tags ordx.ord.r
// @Produce json
// @Security Bearer
// @Success 200 {object} RBlockHeightResp "Successful response"
// @Failure 401 "Invalid API Key"
// @Router /ord/r/blockheight [get]
func (s *Service) getRLastestBlockHeight(c *gin.Context) {
	resp := &RBlockHeightResp{
		BaseResp: common.BaseResp{
			Code: 0,
			Msg:  "ok",
		},
		Data: 0,
	}
	blockHeight, err := ord.OrdRpc.GetRLastestBlockHeight()
	if err != nil {
		resp.Code = -1
		resp.Msg = err.Error()
		c.JSON(http.StatusOK, resp)
		return
	}
	resp.Data = blockHeight
	c.JSON(http.StatusOK, resp)
}

// @Summary ordinal recursive endpoint for get block info
// @Description ordinal recursive endpoint for get block info
// @Tags ordx.ord.r
// @Produce json
// @Param query path string true "block height or block hash"
// @Security Bearer
// @Success 200 {object} RBlockInfoResp "Successful response"
// @Failure 401 "Invalid API Key"
// @Router /ord/r/blockinfo/{query} [get]
func (s *Service) getRBlockInfo(c *gin.Context) {
	resp := &RBlockInfoResp{
		BaseResp: common.BaseResp{
			Code: 0,
			Msg:  "ok",
		},
		Data: nil,
	}
	blockInfo, err := ord.OrdRpc.GetRBlockInfo(c.Param("query"))
	if err != nil {
		resp.Code = -1
		resp.Msg = err.Error()
		c.JSON(http.StatusOK, resp)
		return
	}
	resp.Data = blockInfo
	c.JSON(http.StatusOK, resp)
}

// @Summary ordinal recursive endpoint for get UNIX time stamp of latest block
// @Description ordinal recursive endpoint for get UNIX time stamp of latest block
// @Tags ordx.ord.r
// @Produce json
// @Param query path string true "block height or block hash"
// @Security Bearer
// @Success 200 {object} uint64 "Successful response"
// @Failure 401 "Invalid API Key"
// @Router /ord/r/blocktime [get]
func (s *Service) getRLatestBlockTimestamp(c *gin.Context) {
	resp := &RBlockTimestampResp{
		BaseResp: common.BaseResp{
			Code: 0,
			Msg:  "ok",
		},
		Data: 0,
	}
	blockInfo, err := ord.OrdRpc.GetRLatestBlockTimestamp()
	if err != nil {
		resp.Code = -1
		resp.Msg = err.Error()
		c.JSON(http.StatusOK, resp)
		return
	}
	resp.Data = blockInfo
	c.JSON(http.StatusOK, resp)
}

// @Summary ordinal recursive endpoint for get the first 100 children ids
// @Description ordinal recursive endpoint for get the first 100 children ids
// @Tags ordx.ord.r
// @Produce json
// @Param inscriptionid path string true "inscription ID example: 79b0e9dbfaf11e664abafbd8fec7d734bfa2d59013f25c50aaac1264f700832di0"
// @Param page path string false "page example: 0"
// @Security Bearer
// @Success 200 {object} RInscriptionIdsPaginationResp "Successful response"
// @Failure 401 "Invalid API Key"
// @Router /ord/r/children/{inscriptionid}/{page} [get]
func (s *Service) getRChildrenInscriptionIdList(c *gin.Context) {
	resp := &RInscriptionIdsPaginationResp{
		BaseResp: common.BaseResp{
			Code: 0,
			Msg:  "ok",
		},
		Data: nil,
	}
	inscriptionIdList, err := ord.OrdRpc.GetRChildrenInscriptionIdList(c.Param("inscriptionid"), c.Param("page"))
	if err != nil {
		resp.Code = -1
		resp.Msg = err.Error()
		c.JSON(http.StatusOK, resp)
		return
	}
	resp.Data = inscriptionIdList
	c.JSON(http.StatusOK, resp)
}

// @Summary ordinal recursive endpoint for get inscription info
// @Description ordinal recursive endpoint for get inscription info
// @Tags ordx.ord.r
// @Produce json
// @Param inscriptionid path string true "inscription ID example: 79b0e9dbfaf11e664abafbd8fec7d734bfa2d59013f25c50aaac1264f700832di0"
// @Security Bearer
// @Success 200 {object} RInscriptionResp "Successful response"
// @Failure 401 "Invalid API Key"
// @Router /ord/r/inscription/{inscriptionid} [get]
func (s *Service) getRInscriptionInfo(c *gin.Context) {
	resp := &RInscriptionResp{
		BaseResp: common.BaseResp{
			Code: 0,
			Msg:  "ok",
		},
		Data: nil,
	}
	inscription, err := ord.OrdRpc.GetRInscriptionInfo(c.Param("inscriptionid"))
	if err != nil {
		resp.Code = -1
		resp.Msg = err.Error()
		c.JSON(http.StatusOK, resp)
		return
	}
	resp.Data = inscription
	c.JSON(http.StatusOK, resp)
}

// @Summary ordinal recursive endpoint for get hex-encoded CBOR metadata of an inscription
// @Description ordinal recursive endpoint for get hex-encoded CBOR metadata of an inscription
// @Tags ordx.ord.r
// @Produce json
// @Param inscriptionid path string true "inscription ID example: a4b6fccd00222e79ec0307d52fe9f8bfa3713cd0c170f95065f5d859e0c6a0f5i0"
// @Security Bearer
// @Success 200 {object} RMetadataResp "Successful response"
// @Failure 401 "Invalid API Key"
// @Router /ord/r/metadata/{inscriptionid} [get]
func (s *Service) getRMetadata(c *gin.Context) {
	resp := &RMetadataResp{
		BaseResp: common.BaseResp{
			Code: 0,
			Msg:  "ok",
		},
		Data: "",
	}
	metaData, err := ord.OrdRpc.GetRMetaData(c.Param("inscriptionid"))
	if err != nil {
		resp.Code = -1
		resp.Msg = err.Error()
		c.JSON(http.StatusOK, resp)
		return
	}
	resp.Data = metaData
	c.JSON(http.StatusOK, resp)
}

// @Summary ordinal recursive endpoint for get the first 100 inscription ids on a sat
// @Description ordinal recursive endpoint for get the first 100 inscription ids on a sat
// @Tags ordx.ord.r
// @Produce json
// @Param satnumber path string true "sat number example: 1165647477496168"
// @Param page path string false "page example: 0"
// @Security Bearer
// @Success 200 {object} RInscriptionIdsPaginationResp "Successful response"
// @Failure 401 "Invalid API Key"
// @Router /ord/r/sat/{satnumber}/{page} [get]
func (s *Service) getRSatInscriptionIdList(c *gin.Context) {
	resp := &RInscriptionIdsPaginationResp{
		BaseResp: common.BaseResp{
			Code: 0,
			Msg:  "ok",
		},
		Data: nil,
	}
	inscriptionIdList, err := ord.OrdRpc.GetRSatInscriptionIdList(c.Param("satnumber"), c.Param("page"))
	if err != nil {
		resp.Code = -1
		resp.Msg = err.Error()
		c.JSON(http.StatusOK, resp)
		return
	}
	resp.Data = inscriptionIdList
	c.JSON(http.StatusOK, resp)
}

// @Summary ordinal recursive endpoint for get the inscription id at <INDEX> of all inscriptions on a sat
// @Description ordinal recursive endpoint for get the inscription id at <INDEX> of all inscriptions on a sat
// @Tags ordx.ord.r
// @Produce json
// @Param satnumber path string true "sat number example: 1165647477496168"
// @Param index path string false "page example: -1"
// @Security Bearer
// @Success 200 {object} RInscriptionIdResp "Successful response"
// @Failure 401 "Invalid API Key"
// @Router /ord/r/sat/{satnumber}/at/{index} [get]
func (s *Service) getRSatInscriptionId(c *gin.Context) {
	resp := &RInscriptionIdResp{
		BaseResp: common.BaseResp{
			Code: 0,
			Msg:  "ok",
		},
		Data: nil,
	}
	inscriptionID, err := ord.OrdRpc.GetRSatInscriptionId(c.Param("satnumber"), c.Param("index"))
	if err != nil {
		resp.Code = -1
		resp.Msg = err.Error()
		c.JSON(http.StatusOK, resp)
		return
	}
	resp.Data = inscriptionID
	c.JSON(http.StatusOK, resp)
}
