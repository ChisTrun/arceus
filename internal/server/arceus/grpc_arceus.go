package arceus

import (
	arceus "arceus/api"
	"arceus/internal/feature"
)

func NewServer(feature *feature.Feature) arceus.ArceusServer {
	return &arceusServer{
		feature: feature,
	}
}

type arceusServer struct {
	arceus.UnimplementedArceusServer
	feature *feature.Feature
}
