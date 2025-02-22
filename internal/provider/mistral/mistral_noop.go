package mistral

import (
	arceus "arceus/api"
	"fmt"
)

type Noop struct{}

func (n *Noop) GenerateText(model string, messages []*arceus.Message) (*arceus.Message, error) {
	return nil, fmt.Errorf("mistral ai is not available now")
}

func (n *Noop) GetAvailableModels() ([]string, error) {
	return nil, fmt.Errorf("mistral ai is not available now")
}
