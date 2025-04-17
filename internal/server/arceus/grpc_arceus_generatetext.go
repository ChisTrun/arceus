package arceus

import (
	"context"
	"fmt"

	arceus "arceus/api"
	"arceus/pkg/logger/pkg/logging"
)

func (s *arceusServer) GenerateText(ctx context.Context, request *arceus.GenerateTextRequest) (*arceus.GenerateTextResponse, error) {
	logging.Logger(ctx).Info("Received GenerateText request")

	if err := request.Validate(); err != nil {
		logging.Logger(ctx).Error(fmt.Sprintf("Error: %v", err.Error()))
		return nil, err
	}

	return s.feature.Llm.GenerateText(ctx, request)
}
