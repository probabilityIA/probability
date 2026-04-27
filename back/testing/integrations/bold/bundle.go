package bold

import (
	"crypto/rand"
	"encoding/hex"
	"fmt"
	"net/http"
	"sync"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/secamc93/probability/back/testing/shared/log"
)

type Bold struct {
	logger        log.ILogger
	port          string
	webhookTarget string
	links         map[string]*linkState
	mu            sync.RWMutex
}

type linkState struct {
	ID          string
	Reference   string
	Amount      float64
	Currency    string
	CreatedAt   time.Time
	Status      string
	PaymentID   string
	WebhookSent bool
}

func New(logger log.ILogger, port, webhookTarget string) *Bold {
	return &Bold{
		logger:        logger,
		port:          port,
		webhookTarget: webhookTarget,
		links:         make(map[string]*linkState),
	}
}

func (b *Bold) Start() error {
	gin.SetMode(gin.ReleaseMode)
	router := gin.New()
	router.Use(gin.Recovery())

	router.Use(func(c *gin.Context) {
		start := time.Now()
		c.Next()
		b.logger.Info().Msgf("[bold-mock] %s %s -> %d (%v)",
			c.Request.Method, c.Request.URL.Path, c.Writer.Status(), time.Since(start))
	})

	router.GET("/health", func(c *gin.Context) {
		c.JSON(http.StatusOK, gin.H{"status": "ok", "service": "bold-mock"})
	})

	router.POST("/online/link/v1", b.handleCreateLink)
	router.GET("/online/link/v1/:id", b.handleGetLink)

	router.POST("/admin/simulate/sale-approved/:id", b.handleSimulateApproved)
	router.POST("/admin/simulate/sale-rejected/:id", b.handleSimulateRejected)
	router.POST("/admin/simulate/void-approved/:id", b.handleSimulateVoidApproved)
	router.GET("/admin/links", b.handleListLinks)

	addr := fmt.Sprintf(":%s", b.port)
	b.logger.Info().Msgf("Bold mock server listening on %s (webhook target: %s)", addr, b.webhookTarget)
	return router.Run(addr)
}

func randHex(n int) string {
	buf := make([]byte, n)
	_, _ = rand.Read(buf)
	return hex.EncodeToString(buf)
}
