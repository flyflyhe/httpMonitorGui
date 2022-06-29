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
	"io"
	"sync"
)

var monitorQueue *MonitorQueue
var once sync.Once

type MonitorQueue struct {
	Queue   chan *httpMonitorRpc.MonitorResponse
	Running bool
	m       sync.Mutex
}

func GetMonitorQueue() *MonitorQueue {
	once.Do(func() {
		monitorQueue = &MonitorQueue{Queue: make(chan *httpMonitorRpc.MonitorResponse, 100)}
	})
	return monitorQueue
}

var grpcConnPool = sync.Pool{
	New: func() any {
		conn, err := GetRpcConn()
		if err != nil {
			log.Error().Caller().Str("pool create grpc conn failed", err.Error())
			return nil
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

func ListProxy() ([]string, error) {
	poolConn := grpcConnPool.Get()
	conn, ok := poolConn.(*grpc.ClientConn)
	defer grpcConnPool.Put(poolConn)
	if !ok {
		return nil, errors.New("from pool get conn failed")
	}

	rpcClient := httpMonitorRpc.NewUrlServiceClient(conn)

	if res, err := rpcClient.GetAllProxy(context.Background(), &empty.Empty{}); err != nil {
		return nil, err
	} else {
		return res.ProxyList, err
	}
}

func SetProxy(proxy string) error {
	poolConn := grpcConnPool.Get()
	conn, ok := poolConn.(*grpc.ClientConn)
	defer grpcConnPool.Put(poolConn)
	if !ok {
		return errors.New("from pool get conn failed")
	}

	rpcClient := httpMonitorRpc.NewUrlServiceClient(conn)

	_, err := rpcClient.SetProxy(context.Background(), &httpMonitorRpc.ProxyRequest{Proxy: proxy})
	return err
}

func DeleteProxy(proxy string) error {
	poolConn := grpcConnPool.Get()
	conn, ok := poolConn.(*grpc.ClientConn)
	defer grpcConnPool.Put(poolConn)
	if !ok {
		return errors.New("from pool get conn failed")
	}

	rpcClient := httpMonitorRpc.NewUrlServiceClient(conn)

	_, err := rpcClient.DeleteProxy(context.Background(), &httpMonitorRpc.ProxyRequest{Proxy: proxy})
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

func StartMonitor(monitorQueue *MonitorQueue) error {
	poolConn := grpcConnPool.Get()
	conn, ok := poolConn.(*grpc.ClientConn)
	defer grpcConnPool.Put(poolConn)
	if !ok {
		return errors.New("from pool get conn failed")
	}

	rpcClient := httpMonitorRpc.NewMonitorServerClient(conn)
	if stream, err := rpcClient.Start(context.Background(), &httpMonitorRpc.MonitorRequest{}); err != nil {
		return err
	} else {
		go func() {
			monitorQueue.m.Lock()
			monitorQueue.Running = true
			defer func() {
				monitorQueue.m.Unlock()
				monitorQueue.Running = false
			}()
			for {
				//Recv() 方法接收服务端消息，默认每次Recv()最大消息长度为`1024*1024*4`bytes(4M)
				res, err := stream.Recv()

				monitorQueue.Queue <- res
				// 判断消息流是否已经结束
				if err == io.EOF {
					break
				}
				if err != nil {
					break
				}
			}
		}()

	}

	return nil
}

func StopMonitor(monitorQueue *MonitorQueue) error {
	poolConn := grpcConnPool.Get()
	conn, ok := poolConn.(*grpc.ClientConn)
	defer grpcConnPool.Put(poolConn)
	if !ok {
		return errors.New("from pool get conn failed")
	}

	rpcClient := httpMonitorRpc.NewMonitorServerClient(conn)
	_, err := rpcClient.Stop(context.Background(), &empty.Empty{})
	return err
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
