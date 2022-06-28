package rpc

import (
	"context"
	"crypto/tls"
	"crypto/x509"
	"errors"
	"fmt"
	"github.com/flyflyhe/httpMonitor/config"
	httpMonitorRpc "github.com/flyflyhe/httpMonitor/rpc"
	"github.com/flyflyhe/httpMonitor/services"
	"github.com/golang/protobuf/ptypes/empty"
	"github.com/rs/zerolog/log"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials"
	"sync"
)

var I = 0

var grpcConnPool = sync.Pool{
	New: func() any {
		conn, err := GetRpcConn()
		if err != nil {
			log.Debug().Caller().Str("pool create grpc conn failed", err.Error())
		}
		return conn
	},
}

const address = "localhost:50051"

func Start() {
	services.Start(address)
}

func ListUrl() ([]string, error) {
	poolConn := grpcConnPool.Get()
	conn, ok := poolConn.(*grpc.ClientConn)
	defer grpcConnPool.Put(poolConn)

	if !ok {
		return nil, errors.New("from pool get conn failed")
	}

	rpcClient := httpMonitorRpc.NewUrlServiceClient(conn)

	if res, err := rpcClient.GetAll(context.Background(), &empty.Empty{}); err != nil {
		return nil, err
	} else {
		return res.Urls, nil
	}
}

func SetUrl(url string, interval int32) error {
	poolConn := grpcConnPool.Get()
	conn, ok := poolConn.(*grpc.ClientConn)
	defer grpcConnPool.Put(poolConn)
	if !ok {
		return errors.New("from pool get conn failed")
	}

	rpcClient := httpMonitorRpc.NewUrlServiceClient(conn)

	_, err := rpcClient.SetUrl(context.Background(), &httpMonitorRpc.UrlRequest{Url: url, Interval: interval})
	return err
}

func DeleteUrl(url string) error {
	poolConn := grpcConnPool.Get()
	conn, ok := poolConn.(*grpc.ClientConn)
	defer grpcConnPool.Put(poolConn)
	if !ok {
		return errors.New("from pool get conn failed")
	}

	rpcClient := httpMonitorRpc.NewUrlServiceClient(conn)

	_, err := rpcClient.DeleteUrl(context.Background(), &httpMonitorRpc.UrlRequest{Url: url})
	return err
}

func GetRpcConn() (*grpc.ClientConn, error) {
	tlsCredentials, err := loadClientTLSCredentials()
	if err != nil {
		log.Error().Caller().Msg("cannot load TLS credentials: " + err.Error())
		return nil, err
	}
	if err != nil {
		log.Error().Caller().Msg("credentials.NewClientTLSFromFile err: " + err.Error())
		return nil, err
	}
	conn, err := grpc.Dial(address, grpc.WithTransportCredentials(tlsCredentials))
	if err != nil {
		log.Error().Caller().Msg("did not connect: " + err.Error())
		return nil, err
	}

	return conn, nil
}

func loadClientTLSCredentials() (credentials.TransportCredentials, error) {

	certPool := x509.NewCertPool()
	if !certPool.AppendCertsFromPEM(config.GetRoot()) {
		return nil, fmt.Errorf("failed to add server CA's certificate")
	}

	// Load client's certificate and private key
	clientCert, err := tls.X509KeyPair(config.GetClientCertChain(), config.GetClientPrivateKey())
	if err != nil {
		return nil, err
	}

	// Create the credentials and return it
	tlsConfig := &tls.Config{
		ServerName:   "test.com", //生成的证书通用名称 必须一致
		Certificates: []tls.Certificate{clientCert},
		RootCAs:      certPool,
	}

	return credentials.NewTLS(tlsConfig), nil
}
