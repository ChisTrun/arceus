package server

import (
	dbe "arceus/pkg/database/pkg/ent"
	"arceus/pkg/ent"
	"arceus/pkg/ent/migrate"
	mykit "arceus/pkg/mykit/pkg/api"
	"context"

	"github.com/grpc-ecosystem/grpc-gateway/v2/runtime"
	"go.uber.org/zap"
	"google.golang.org/grpc/reflection"
	"google.golang.org/protobuf/encoding/protojson"

	pb0 "arceus/api"
	"arceus/internal/feature"
	"arceus/internal/provider"
	"arceus/internal/provider/gemini"
	"arceus/internal/provider/mistral"
	"arceus/internal/provider/openai"
	"arceus/internal/repository"
	"arceus/internal/server/arceus"
	config "arceus/pkg/config"
)

// Serve ...
func Serve(cfg *config.Config) {
	service := newService(cfg, []mykit.Option{}...)
	logger := service.Logger()

	server := service.Server()

	drv, err := dbe.Open("mysql_rum", cfg.GetDatabase())
	ent := ent.NewClient(ent.Driver(drv))
	defer func() {
		if err := ent.Close(); err != nil {
			logger.Fatal("can not close ent client", zap.Error(err))
		}
	}()
	if err != nil {
		logger.Fatal("can not open ent client", zap.Error(err))
	}
	if err = ent.Schema.Create(context.Background(), migrate.WithDropIndex(true)); err != nil {
		logger.Fatal("can not init my database", zap.Error(err))
	}

	repo := repository.New(ent)

	mistralProvider := mistral.New(cfg)
	openaiProvider := openai.New(cfg)
	geminiProvider := gemini.New(cfg)

	feature := feature.New(repo, []provider.Provider{mistralProvider, openaiProvider, geminiProvider})

	arceusServer := arceus.NewServer(feature)

	grpcGatewayMux := runtime.NewServeMux(
		runtime.WithMarshalerOption(runtime.MIMEWildcard, &runtime.JSONPb{
			MarshalOptions: protojson.MarshalOptions{
				UseProtoNames:   true,
				EmitUnpopulated: true,
				UseEnumNumbers:  false,
			},
		}),
	)
	service.HttpServeMux().Handle("/api/", grpcGatewayMux)

	err = pb0.RegisterArceusHandlerServer(context.Background(), grpcGatewayMux, arceusServer)
	if err != nil {
		logger.Fatal("can not register http sibel server", zap.Error(err))
	}

	pb0.RegisterArceusServer(server, arceusServer)
	// Register reflection service on gRPC server.
	// Please remove if you it's not necessary for your service
	reflection.Register(server)

	service.Serve()
}
