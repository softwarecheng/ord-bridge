package service

import (
	"net"
	"net/http"
	"strings"
	"time"

	"github.com/didip/tollbooth/v7"
	"github.com/didip/tollbooth/v7/limiter"
	"github.com/gin-gonic/gin"
	"github.com/softwarecheng/ord-bridge/common/log"
	serverCommon "github.com/softwarecheng/ord-bridge/service/common"
	"gopkg.in/yaml.v2"
)

func (s *ServiceManager) InitApiConf(cfgData any) error {
	if cfgData == nil {
		return nil
	}
	readApiAuthConf := func() error {
		s.apiConfMutex.Lock()
		defer s.apiConfMutex.Unlock()

		raw, err := yaml.Marshal(cfgData)
		if err != nil {
			return err
		}
		s.apiConf = &serverCommon.API{}
		err = yaml.Unmarshal(raw, s.apiConf)
		if err != nil {
			return err
		}
		s.initApiConf = true
		return nil
	}

	err := readApiAuthConf()
	if err != nil {
		return err
	}
	go func() {
		ticker := time.NewTicker(30 * time.Second)
		defer ticker.Stop()
		for range ticker.C {
			err := readApiAuthConf()
			if err != nil {
				log.Log.Errorf("rpc.readApiAuthConf-> readApiAuthConf error: %v", err)
			}
		}
	}()
	return nil
}

func (s *ServiceManager) applyApiConf(r *gin.Engine, basePath string) error {
	localIpList := make([]string, 0)
	if len(localIpList) == 0 {
		addrs, err := net.InterfaceAddrs()
		if err != nil {
			return err
		}
		for _, addr := range addrs {
			ipNet, ok := addr.(*net.IPNet)
			if ok && ipNet.IP.To4() != nil {
				localIpList = append(localIpList, ipNet.IP.String())
			}
		}
		localIpList = append(localIpList, "localhost")
	}

	r.Use(func(c *gin.Context) {
		if !s.initApiConf {
			c.Next()
			return
		}
		for _, ip := range localIpList {
			if strings.Contains(c.Request.Host, ip) {
				c.Next()
				return
			}
		}

		s.apiConfMutex.Lock()
		defer s.apiConfMutex.Unlock()
		for _, apiUrl := range s.apiConf.NoLimitAPIList {
			if basePath+apiUrl == c.Request.URL.Path {
				c.Next()
				return
			}
		}

		apiKey := c.GetHeader("Authorization")
		var hostConf *serverCommon.APIKeyList
		for _, apiUrl := range s.apiConf.APIKeyList {
			if apiKey == apiUrl.APIKey {
				hostConf = &apiUrl
				break
			}
		}
		if hostConf == nil {
			c.JSON(http.StatusUnauthorized, gin.H{"error": "Invalid API Key"})
			c.Abort()
			return
		}
		if hostConf.RateLimit.PerSecond <= 0 || hostConf.RateLimit.PerDay <= 0 {
			c.JSON(http.StatusTooManyRequests, gin.H{"error": "Rate limit not configured"})
			c.Abort()
			return
		}
		authorization := strings.Trim(c.GetHeader("Authorization"), " ")
		if authorization == "" || authorization != hostConf.UserName {
			c.JSON(http.StatusUnauthorized, gin.H{"error": "Invalid API Key"})
			c.Abort()
			return
		}

		var rateLimit *RateLimit
		v, ok := s.apiLimitMap.Load(apiKey)
		if !ok {
			lmt := tollbooth.NewLimiter(float64(hostConf.RateLimit.PerSecond), &limiter.ExpirableOptions{DefaultExpirationTTL: time.Hour})
			lmt.SetOnLimitReached(func(w http.ResponseWriter, r *http.Request) {
				c.JSON(http.StatusTooManyRequests, gin.H{"error": "Rate limit exceeded"})
				c.Abort()
			})
			rateLimit := &RateLimit{limit: lmt, reqCount: 0}
			s.apiLimitMap.Store(apiKey, rateLimit)
		} else {
			rateLimit = v.(*RateLimit)
		}

		now := time.Now()
		today := time.Date(now.Year(), now.Month(), now.Day(), 0, 0, 0, 0, time.Local)
		if now.Before(today.AddDate(0, 0, 1)) {
			rateLimit.reqCount++
			if rateLimit.reqCount > hostConf.RateLimit.PerDay {
				c.JSON(http.StatusTooManyRequests, gin.H{"error": "Rate limit exceeded"})
				c.Abort()
				return
			}
		} else {
			rateLimit.reqCount = 1
		}

		httpError := tollbooth.LimitByRequest(rateLimit.limit, c.Writer, c.Request)
		if httpError != nil {
			c.JSON(http.StatusTooManyRequests, gin.H{"error": "Rate limit exceeded"})
			c.Abort()
			return
		}
		c.Next()
	})

	return nil
}
