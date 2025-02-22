package arceus

import (
	"context"

	arceus "arceus/api"
)

func (s *arceusServer) GenerateText(ctx context.Context, request *arceus.GenerateTextRequest) (*arceus.GenerateTextResponse, error) {
	if err := request.Validate(); err != nil {
		return nil, err
	}

	return s.feature.Llm.GenerateText(ctx, request)
}
