package bitcoind

import (
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/softwarecheng/ord-bridge/rpc/bitcoind"
	"github.com/softwarecheng/ord-bridge/service/common"
)

// @Summary send Raw Transaction
// @Description send Raw Transaction
// @Tags ordx.btc
// @Produce json
// @Param signedTxHex body string true "Signed transaction hex"
// @Param maxfeerate body number false "Reject transactions whose fee rate is higher than the specified value, expressed in BTC/kB.default:"0.01"
// @Security Bearer
// @Success 200 {object} SendRawTxResp "Successful response"
// @Failure 401 "Invalid API Key"
// @Router /btc/tx [post]
func (s *Service) sendRawTx(c *gin.Context) {
	resp := &SendRawTxResp{
		BaseResp: common.BaseResp{
			Code: 0,
			Msg:  "ok",
		},
		Data: "",
	}
	var req SendRawTxReq
	if err := c.ShouldBindJSON(&req); err != nil {
		resp.Code = -1
		resp.Msg = err.Error()
		c.JSON(http.StatusOK, resp)
		return
	}

	txid, err := bitcoind.BitcoinRpc.SendRawTransaction(req.SignedTxHex, req.Maxfeerate)
	if err != nil {
		resp.Code = -1
		resp.Msg = err.Error()
		c.JSON(http.StatusOK, resp)
		return
	}

	resp.Data = txid
	c.JSON(http.StatusOK, resp)
}

// @Summary get raw block with blockhash
// @Description get raw block with blockhash
// @Tags ordx.btc
// @Produce json
// @Param blockHash path string true "blockHash"
// @Security Bearer
// @Success 200 {object} RawBlockResp "Successful response"
// @Failure 401 "Invalid API Key"
// @Router /btc/getrawblock [get]
func (s *Service) getRawBlock(c *gin.Context) {
	resp := &RawBlockResp{
		BaseResp: common.BaseResp{
			Code: 0,
			Msg:  "ok",
		},
		Data: "",
	}
	blockHash := c.Param("blockHash")
	data, err := bitcoind.BitcoinRpc.GetRawBlock(blockHash)
	if err != nil {
		resp.Code = -1
		resp.Msg = err.Error()
		c.JSON(http.StatusOK, resp)
		return
	}

	resp.Data = data
	c.JSON(http.StatusOK, resp)
}
