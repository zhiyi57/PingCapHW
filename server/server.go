// Copyright 2019 PingCAP, Inc.
//
// Licensed under the Apache License, Version 2.0 (the "License");
// you may not use this file except in compliance with the License.
// You may obtain a copy of the License at
//
//     http://www.apache.org/licenses/LICENSE-2.0
//
// Unless required by applicable law or agreed to in writing, software
// distributed under the License is distributed on an "AS IS" BASIS,
// See the License for the specific language governing permissions and
// limitations under the License.

package main

import (
	"context"
	"flag"
	"fmt"
	"log"
	"net"
	pb "github.com/zhiyi57/PingCapHW/proto/kvrawpb"
	"google.golang.org/grpc"
	"github.com/pingcap/tidb/store/tikv"
	"github.com/pingcap/tidb/config"
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
	pairs := [] *pb.KvPair{}

	for i := 0; i < len(keys) && i < len(values); i++ {
		tmp := pb.KvPair{Key:keys[i],Value:values[i]}
		pairs = append(pairs, &tmp)
	}
	return &pb.RawScanResponse{Kvs:pairs}, nil
}

func (t *tikvServer) Close() error {
	err := t.cli.Close()
	if err != nil {
		return err
	}
	return nil
}


func NewServer() *tikvServer {
	pdCli, err := tikv.NewRawKVClient([]string{"localhost:2379"}, config.Security{})
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