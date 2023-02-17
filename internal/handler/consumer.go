// Package handler
package handler

import (
	"context"

	"GolangBookingApp/internal/appctx"
	"GolangBookingApp/internal/consts"
	uContract "GolangBookingApp/internal/ucase/contract"
	"GolangBookingApp/pkg/awssqs"
)

// SQSConsumerHandler sqs consumer message processor handler
func SQSConsumerHandler(msgHandler uContract.MessageProcessor) awssqs.MessageProcessorFunc {
	return func(decoder *awssqs.MessageDecoder) error {
		return msgHandler.Serve(context.Background(), &appctx.ConsumerData{
			Body:        []byte(*decoder.Body),
			Key:         []byte(*decoder.MessageId),
			ServiceType: consts.ServiceTypeConsumer,
		})
	}
}
