package mock

import (
	"context"
	"crypto/tls"
	"crypto/x509"
	"errors"
	"log"
	"math"
	"net/http"
	"net/url"
	"reflect"
	"slices"
	"strings"
	"sync"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/goccy/go-yaml"
	papiv2 "github.com/gotify/plugin-api/v2"
	"github.com/gotify/plugin-api/v2/generated/protobuf"
	"github.com/gotify/plugin-api/v2/transport"
	"google.golang.org/grpc"
	"google.golang.org/grpc/keepalive"
	"google.golang.org/protobuf/types/known/emptypb"
)

var (
	httpTimeout = 10 * time.Second
	pingRate    = 4 * time.Second
)

// ModulePath is for convenient access of the module path of this mock plugin
const ModulePath = "github.com/gotify/server/v2/plugin/testing/mock"

// Name is for convenient access of the module path of the name of this mock plugin
const Name = "mock plugin"

// PluginServer is a mock plugin server.
type PluginServer struct {
	plugin        *Plugin
	webhookServer http.HandlerFunc
	rpcServer     *grpc.Server
	http.Server
}

func NewPluginServer(cliArgs []string) (*PluginServer, error) {

	cli, err := papiv2.ParsePluginCli(cliArgs)
	if err != nil {
		log.Fatalf("Failed to parse CLI flags: %v", err)
	}
	defer cli.Close()

	rootCAs := x509.NewCertPool()
	certificateChain, err := cli.Kex(ModulePath, rootCAs)
	if err != nil {
		return nil, err
	}

	tlsConfig := &tls.Config{
		Certificates: certificateChain,
		RootCAs:      rootCAs,
		ClientAuth:   tls.RequireAndVerifyClientCert,
		ClientCAs:    rootCAs,
	}

	rpcServer := grpc.NewServer(grpc.KeepaliveEnforcementPolicy(keepalive.EnforcementPolicy{
		MinTime:             httpTimeout,
		PermitWithoutStream: true,
	}), grpc.ConnectionTimeout(httpTimeout))
	if !cli.Debug {
		gin.SetMode(gin.ReleaseMode)
	}

	self := &PluginServer{
		plugin: &Plugin{
			shutdown:     make(chan struct{}),
			shutdownOnce: &sync.Once{},
		},
		rpcServer: rpcServer,
	}

	protobuf.RegisterPluginServer(rpcServer, self.plugin)
	protobuf.RegisterDisplayerServer(rpcServer, self.plugin)
	protobuf.RegisterConfigurerServer(rpcServer, self.plugin)

	protocols := new(http.Protocols)
	protocols.SetHTTP1(true)
	protocols.SetHTTP2(true)
	self.Server = http.Server{
		Handler:           self,
		TLSConfig:         tlsConfig,
		Protocols:         protocols,
		ReadTimeout:       httpTimeout,
		ReadHeaderTimeout: httpTimeout,
		WriteTimeout:      httpTimeout,
	}

	return self, nil
}

func (h *PluginServer) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	if r.TLS == nil {
		http.Error(w, "Must use TLS", http.StatusUpgradeRequired)
		return
	}

	pluginRpcHostName := transport.BuildPluginTLSName(transport.PurposePluginRPC, ModulePath)

	if r.TLS.ServerName == pluginRpcHostName {
		if r.ProtoMajor != 2 {
			http.Error(w, "Must use HTTP/2", http.StatusHTTPVersionNotSupported)
			return
		}
		if !strings.HasPrefix(r.Header.Get("Content-Type"), "application/grpc") {
			http.Error(w, "Must use application/grpc content type", http.StatusUnsupportedMediaType)
			return
		}
		h.rpcServer.ServeHTTP(w, r)

		return
	}

	pluginWebhookHostName := transport.BuildPluginTLSName(transport.PurposePluginWebhook, ModulePath)
	if r.TLS.ServerName == pluginWebhookHostName {
		h.webhookServer(w, r)
		return
	}

	http.Error(w, "Virtual host not found", http.StatusNotFound)
}

type Plugin struct {
	shutdown     chan struct{}
	shutdownOnce *sync.Once
	protobuf.UnimplementedPluginServer
	protobuf.UnimplementedDisplayerServer
	protobuf.UnimplementedConfigurerServer
}

func (s *Plugin) GetPluginInfo(ctx context.Context, req *emptypb.Empty) (*protobuf.Info, error) {
	return &protobuf.Info{
		Name:       Name,
		ModulePath: ModulePath,
	}, nil
}

func (s *Plugin) Display(ctx context.Context, req *protobuf.DisplayRequest) (*protobuf.DisplayResponse, error) {
	instance, err := s.shim.getInstanceByUserId(req.User.Id)
	if err != nil {
		return nil, err
	}
	if displayer, ok := instance.(papiv1.Displayer); ok {
		location, err := url.Parse(req.Location)
		if err != nil {
			return nil, err
		}
		return &protobuf.DisplayResponse{
			Response: &protobuf.DisplayResponse_Markdown{
				Markdown: displayer.GetDisplay(location),
			},
		}, nil
	}
	return nil, errors.New("instance does not implement displayer")
}

func (s *Plugin) DefaultConfig(ctx context.Context, req *protobuf.DefaultConfigRequest) (*protobuf.Config, error) {
	instance, err := s.shim.getInstanceByUserId(req.User.Id)
	if err != nil {
		return nil, err
	}
	if configurer, ok := instance.(papiv1.Configurer); ok {
		defaultConfig := configurer.DefaultConfig()
		bytes, err := yaml.Marshal(defaultConfig)
		if err != nil {
			return nil, err
		}
		return &protobuf.Config{
			Config: string(bytes),
		}, nil
	}
	return nil, errors.New("instance does not implement configurer")
}

func (s *Plugin) ValidateAndSetConfig(ctx context.Context, req *protobuf.ValidateAndSetConfigRequest) (*protobuf.ValidateAndSetConfigResponse, error) {
	instance, err := s.shim.getInstanceByUserId(req.User.Id)
	if err != nil {
		return nil, err
	}
	if configurer, ok := instance.(papiv1.Configurer); ok {
		currentConfig := configurer.DefaultConfig()
		if req.Config != nil {
			if reflect.TypeOf(currentConfig).Kind() == reflect.Pointer {
				yaml.Unmarshal([]byte(req.Config.Config), currentConfig)
			} else {
				yaml.Unmarshal([]byte(req.Config.Config), &currentConfig)
			}
		}
		if err := configurer.ValidateAndSetConfig(currentConfig); err != nil {
			return &protobuf.ValidateAndSetConfigResponse{
				Response: &protobuf.ValidateAndSetConfigResponse_ValidationError{
					ValidationError: &protobuf.Error{
						Message: err.Error(),
					},
				},
			}, nil
		}
		return &protobuf.ValidateAndSetConfigResponse{
			Response: &protobuf.ValidateAndSetConfigResponse_Success{
				Success: new(emptypb.Empty),
			},
		}, nil
	}
	return nil, errors.New("instance does not implement configurer")
}

func (s *Plugin) GracefulShutdown(ctx context.Context, req *emptypb.Empty) (*emptypb.Empty, error) {
	s.shutdownOnce.Do(func() {
		close(s.shutdown)
	})
	return &emptypb.Empty{}, nil
}

func (s *Plugin) RunUserInstance(req *protobuf.UserInstanceRequest, stream protobuf.Plugin_RunUserInstanceServer) error {
	if req.User.Id > math.MaxUint {
		return errors.New("user id is too large")
	}

	unlockOnce := new(sync.Once)

	s.shim.mu.Lock()

	defer unlockOnce.Do(func() {
		s.shim.mu.Unlock()
	})

	instance, alreadyRunning := s.shim.instances[req.User.Id]

	if !alreadyRunning {
		var err error
		instance, err = s.shim.compatV1.GetInstance(papiv1.UserContext{
			ID:    uint(req.User.Id),
			Name:  req.User.Name,
			Admin: req.User.Admin,
		})
		if err != nil {
			return err
		}

		// enable supported capabilities
		if _, ok := instance.(papiv1.Displayer); ok {
			if slices.Contains(req.ServerInfo.Capabilities, protobuf.Capability_DISPLAYER) {
				if err := stream.Send(&protobuf.InstanceUpdate{
					Update: &protobuf.InstanceUpdate_Capable{
						Capable: protobuf.Capability_DISPLAYER,
					},
				}); err != nil {
					return err
				}
			} else {
				return errors.New("displayer not supported by server but V1 API does not support backwards compatibility")
			}
		}
		if _, ok := instance.(papiv1.Messenger); ok {
			if slices.Contains(req.ServerInfo.Capabilities, protobuf.Capability_MESSENGER) {
				if err := stream.Send(&protobuf.InstanceUpdate{
					Update: &protobuf.InstanceUpdate_Capable{
						Capable: protobuf.Capability_MESSENGER,
					},
				}); err != nil {
					return err
				}
			} else {
				return errors.New("messenger not supported by server but V1 API does not support backwards compatibility")
			}
		}
		if _, ok := instance.(papiv1.Configurer); ok {
			if slices.Contains(req.ServerInfo.Capabilities, protobuf.Capability_CONFIGURER) {
				if err := stream.Send(&protobuf.InstanceUpdate{
					Update: &protobuf.InstanceUpdate_Capable{
						Capable: protobuf.Capability_CONFIGURER,
					},
				}); err != nil {
					return err
				}
			} else {
				return errors.New("configurer not supported by server but V1 API does not support backwards compatibility")
			}
		}
		if _, ok := instance.(papiv1.Storager); ok {
			if slices.Contains(req.ServerInfo.Capabilities, protobuf.Capability_STORAGER) {
				if err := stream.Send(&protobuf.InstanceUpdate{
					Update: &protobuf.InstanceUpdate_Capable{
						Capable: protobuf.Capability_STORAGER,
					},
				}); err != nil {
					return err
				}
			} else {
				return errors.New("storager not supported by server but V1 API does not support backwards compatibility")
			}
		}
		if _, ok := instance.(papiv1.Webhooker); ok {
			if slices.Contains(req.ServerInfo.Capabilities, protobuf.Capability_WEBHOOKER) {
				if err := stream.Send(&protobuf.InstanceUpdate{
					Update: &protobuf.InstanceUpdate_Capable{
						Capable: protobuf.Capability_WEBHOOKER,
					},
				}); err != nil {
					return err
				}
			} else {
				return errors.New("webhooker not supported by server but V1 API does not support backwards compatibility")
			}
		}

		if messenger, ok := instance.(papiv1.Messenger); ok {
			if slices.Contains(req.ServerInfo.Capabilities, protobuf.Capability_MESSENGER) {
				messenger.SetMessageHandler(&shimV1MessageHandler{
					stream: &stream,
				})
			} else {
				return errors.New("messenger not supported by server but V1 API does not support backwards compatibility")
			}
		}

		if configurer, ok := instance.(papiv1.Configurer); ok {
			if slices.Contains(req.ServerInfo.Capabilities, protobuf.Capability_CONFIGURER) {
				currentConfig := configurer.DefaultConfig()
				if req.Config != nil {
					if err := yaml.Unmarshal(req.Config, &currentConfig); err != nil {
						return err
					}
					if err := configurer.ValidateAndSetConfig(currentConfig); err != nil {
						return err
					}
				}
			} else {
				return errors.New("configurer not supported by server but V1 API does not support backwards compatibility")
			}
		}

		if storager, ok := instance.(papiv1.Storager); ok {
			storageHandler := &shimV1StorageHandler{
				mutex:          &sync.RWMutex{},
				currentStorage: req.Storage,
				stream:         &stream,
			}
			storager.SetStorageHandler(storageHandler)
		}

		if webhooker, ok := instance.(papiv1.Webhooker); ok {
			if req.WebhookBasePath != nil {
				group := s.shim.gin.Group(*req.WebhookBasePath)
				webhooker.RegisterWebhook(*req.WebhookBasePath, group)
			}
		}
	}

	if err := instance.Enable(); err != nil {
		return err
	}

	defer instance.Disable()

	s.shim.instances[req.User.Id] = instance
	unlockOnce.Do(func() {
		s.shim.mu.Unlock()
	})

	ticker := time.NewTicker(pingRate)

	defer ticker.Stop()
	for {
		select {
		case <-ticker.C:
			if err := stream.Send(&protobuf.InstanceUpdate{
				Update: &protobuf.InstanceUpdate_Ping{
					Ping: new(emptypb.Empty),
				},
			}); err != nil {
				return err
			}
		case <-s.shim.shutdown:
			return nil
		}
	}
}
