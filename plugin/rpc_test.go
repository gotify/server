package plugin

import (
	"context"
	"crypto/ed25519"
	"crypto/rand"
	"crypto/tls"
	"crypto/x509"
	"testing"
	"time"

	papiv2 "github.com/gotify/plugin-api/v2"
	"github.com/gotify/plugin-api/v2/generated/protobuf"
	"github.com/stretchr/testify/assert"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials"
	"google.golang.org/protobuf/types/known/emptypb"
)

type dummyPlugin struct {
	protobuf.UnimplementedPluginServer
}

var dummyPluginInfo = &protobuf.Info{
	Name:        "dummy",
	Version:     "test",
	Description: "dummy plugin",
	Author:      "gotify",
	License:     "MIT",
	ModulePath:  "dummy.example",
}

func (p *dummyPlugin) GetPluginInfo(ctx context.Context, req *emptypb.Empty) (*protobuf.Info, error) {
	return dummyPluginInfo, nil
}

func TestRPC(t *testing.T) {
	rpc := NewServerMux(ServerVersionInfo{
		Version:   "test",
		Commit:    "test",
		BuildDate: time.Now().Format(time.RFC3339),
	})
	defer rpc.Close()
	_, pluginPriv, err := ed25519.GenerateKey(rand.Reader)
	assert.NoError(t, err)
	pluginCsrBytes, err := x509.CreateCertificateRequest(rand.Reader, new(x509.CertificateRequest), pluginPriv)
	pluginCsr, err := x509.ParseCertificateRequest(pluginCsrBytes)
	assert.NoError(t, err)
	assert.NoError(t, pluginCsr.CheckSignature())
	pluginCert, err := rpc.SignPluginCSR(dummyPluginInfo.ModulePath, pluginCsr)
	assert.NoError(t, err)
	pluginCertParsed, err := x509.ParseCertificate(pluginCert)
	assert.NoError(t, err)
	pluginTlsConfig := rpc.tlsClient.ServerTLSConfig()
	pluginTlsConfig.Certificates = []tls.Certificate{
		{
			Certificate: [][]byte{pluginCert},
			PrivateKey:  pluginPriv,
		},
	}
	pluginListener, err := papiv2.NewListener()
	assert.NoError(t, err)

	defer pluginListener.Close()
	pluginListenerTarget := pluginListener.Addr().String()
	if pluginListener.Addr().Network() == "unix" {
		pluginListenerTarget = "unix://" + pluginListenerTarget
	}

	pluginServer := grpc.NewServer(grpc.Creds(credentials.NewTLS(pluginTlsConfig)))
	protobuf.RegisterPluginServer(pluginServer, &dummyPlugin{})
	go pluginServer.Serve(pluginListener)
	defer pluginServer.GracefulStop()

	conn, err := rpc.RegisterPlugin(pluginListenerTarget, dummyPluginInfo.ModulePath)
	assert.NoError(t, err)
	defer conn.Close()

	caCertPool := x509.NewCertPool()
	caCertPool.AddCert(rpc.CACert())
	pluginClientTlsConfig := &tls.Config{
		RootCAs:    caCertPool,
		ServerName: papiv2.ServerTLSName,
		Certificates: []tls.Certificate{
			{
				Certificate: [][]byte{pluginCert},
				PrivateKey:  pluginPriv,
				Leaf:        pluginCertParsed,
			},
		},
		ClientAuth: tls.RequireAndVerifyClientCert,
		ClientCAs:  caCertPool,
	}

	pluginClient, err := grpc.NewClient(pluginListenerTarget, grpc.WithTransportCredentials(credentials.NewTLS(pluginClientTlsConfig)))
	if err != nil {
		t.Fatal(err)
	}
	defer pluginClient.Close()
}
