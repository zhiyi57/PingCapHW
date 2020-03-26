package main

import (
	"context"
	"flag"
	"fmt"
	"github.com/pingcap/tidb/config"
	"github.com/pingcap/tidb/store/tikv"
	pb "github.com/zhiyi57/PingCapHW/proto/kvrawpb"
	"google.golang.org/grpc"
	"log"
	"net"
)

var (
	port = flag.Int("port", 8080, "The server port")
)

type tikvServer struct {
	cli *tikv.RawKVClient
}

func (t *tikvServer) RawPut(ctx context.Context, r *pb.RawPutRequest) (*pb.RawPutResponse, error) {
	err := t.cli.Put(r.Key, r.Value)
	if err != nil {
		return nil, err
	}
	return &pb.RawPutResponse{}, nil
}

func (t *tikvServer) RawGet(ctx context.Context, r *pb.RawGetRequest) (*pb.RawGetResponse, error) {
	value, err := t.cli.Get(r.Key)
	if err != nil {
		return nil, err
	}
	return &pb.RawGetResponse{Value: value}, nil
}

func (t *tikvServer) RawDelete(ctx context.Context, r *pb.RawDeleteRequest) (*pb.RawDeleteResponse, error) {
	err := t.cli.Delete(r.Key)
	if err != nil {
		return nil, err
	}
	return nil, nil
}

func (t *tikvServer) RawScan(ctx context.Context, r *pb.RawScanRequest) (*pb.RawScanResponse, error) {
	keys, values, err := t.cli.Scan(r.StartKey, r.EndKey, int(r.Limit))
	if err != nil {
		return nil, err
	}
	pairs := []*pb.KvPair{}

	for i := 0; i < len(keys) && i < len(values); i++ {
		tmp := pb.KvPair{Key: keys[i], Value: values[i]}
		pairs = append(pairs, &tmp)
	}
	return &pb.RawScanResponse{Kvs: pairs}, nil
}

func (t *tikvServer) Close() error {
	err := t.cli.Close()
	if err != nil {
		return err
	}
	return nil
}

func NewServer() *tikvServer {
	pdCli, err := tikv.NewRawKVClient([]string{"10.0.26.75:2379"}, config.Security{})
	if err != nil {
		panic(err)
	}
	log.Println("Set up kv connection for proxy")
	return &tikvServer{cli: pdCli}
}

func main() {
	flag.Parse()
	grpcServer := grpc.NewServer()
	ts := NewServer()
	defer ts.Close()

	pb.RegisterTikvServer(grpcServer, ts)
	lis, err := net.Listen("tcp", fmt.Sprintf("localhost:%d", *port))
	if err != nil {
		log.Fatalf("failed to listen: %v", err)
	}
	log.Println("server is up and listen to port ", *port)
	grpcServer.Serve(lis)
}
