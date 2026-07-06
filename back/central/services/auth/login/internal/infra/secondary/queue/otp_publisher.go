package queue

import (
	"context"
	"encoding/json"

	"github.com/secamc93/probability/back/central/services/auth/login/internal/domain"
	"github.com/secamc93/probability/back/central/shared/log"
	"github.com/secamc93/probability/back/central/shared/rabbitmq"
)

type otpPublisher struct {
	queue  rabbitmq.IQueue
	logger log.ILogger
}

type passwordResetOTPMessage struct {
	Phone          string `json:"phone"`
	Code           string `json:"code"`
	UserName       string `json:"user_name"`
	ExpiresMinutes int    `json:"expires_minutes"`
}

func New(queue rabbitmq.IQueue, logger log.ILogger) domain.IOTPEventPublisher {
	return &otpPublisher{queue: queue, logger: logger}
}

func (p *otpPublisher) PublishPasswordResetOTP(ctx context.Context, event domain.PasswordResetOTPEvent) error {
	if p.queue == nil {
		p.logger.Error(ctx).Msg("[OTPPublisher] RabbitMQ no disponible - no se puede enviar OTP por WhatsApp")
		return nil
	}

	payload := passwordResetOTPMessage{
		Phone:          event.Phone,
		Code:           event.Code,
		UserName:       event.UserName,
		ExpiresMinutes: event.ExpiresMinutes,
	}

	body, err := json.Marshal(payload)
	if err != nil {
		return err
	}

	return p.queue.Publish(ctx, rabbitmq.QueueAuthPasswordResetOTP, body)
}
