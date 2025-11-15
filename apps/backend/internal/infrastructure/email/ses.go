// apps/backend/internal/infrastructure/email/ses.go

package email

import (
	"context"
	"github.com/aws/aws-sdk-go-v2/aws"
	"github.com/aws/aws-sdk-go-v2/service/sesv2"
	"github.com/aws/aws-sdk-go-v2/service/sesv2/types"
)

// SESClient رابطی برای ارسال ایمیل است.
type SESClient struct {
	client *sesv2.Client
	sender string
}

// NewSESClient یک نمونه جدید از SESClient ایجاد می‌کند.
func NewSESClient(cfg aws.Config, senderEmail string) *SESClient {
	return &SESClient{
		client: sesv2.NewFromConfig(cfg),
		sender: senderEmail,
	}
}

// SendEmail یک ایمیل با محتوای HTML ارسال می‌کند.
func (s *SESClient) SendEmail(ctx context.Context, recipient, subject, htmlBody, textBody string) error {
	input := &sesv2.SendEmailInput{
		FromEmailAddress: &s.sender,
		Destination: &types.Destination{
			ToAddresses: []string{recipient},
		},
		Content: &types.EmailContent{
			Simple: &types.Message{
				Subject: &types.Content{
					Data:    aws.String(subject),
					Charset: aws.String("UTF-8"),
				},
				Body: &types.Body{
					Html: &types.Content{
						Data:    aws.String(htmlBody),
						Charset: aws.String("UTF-8"),
					},
					Text: &types.Content{
						Data:    aws.String(textBody),
						Charset: aws.String("UTF-8"),
					},
				},
			},
		},
	}

	_, err := s.client.SendEmail(ctx, input)
	return err
}