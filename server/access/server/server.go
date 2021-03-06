package server

import (
	"github.com/golang/glog"
	"github.com/golang/protobuf/proto"
	"github.com/oikomi/FishChatServer2/common/ecode"
	"github.com/oikomi/FishChatServer2/libnet"
	"github.com/oikomi/FishChatServer2/protocol/external"
	"github.com/oikomi/FishChatServer2/server/access/client"
	"github.com/oikomi/FishChatServer2/server/access/conf"
	"github.com/oikomi/FishChatServer2/server/access/rpc"
	"github.com/oikomi/FishChatServer2/service_discovery/etcd"
)

type Server struct {
	Server    *libnet.Server
	RPCServer *rpc.RPCServer
}

func New() (s *Server) {
	s = &Server{}
	go rpc.RPCServerInit()
	return
}

func (s *Server) sessionLoop(client *client.Client) {
	for {
		reqData, err := client.Session.Receive()
		if err != nil {
			glog.Error(err)
		}
		if reqData != nil {
			baseCMD := &external.Base{}
			if err = proto.Unmarshal(reqData, baseCMD); err != nil {
				if err = client.Session.Send(&external.Error{
					Cmd:     external.ErrServerCMD,
					ErrCode: ecode.ServerErr.Uint32(),
					ErrStr:  ecode.ServerErr.String(),
				}); err != nil {
					glog.Error(err)
				}
				continue
			}
			if err = client.Parse(baseCMD.Cmd, reqData); err != nil {
				glog.Error(err)
				continue
			}
		}
	}
}

func (s *Server) Loop(rpcClient *rpc.RPCClient) {
	for {
		session, err := s.Server.Accept()
		if err != nil {
			glog.Error(err)
		}
		glog.Info("a new client ", session.ID())
		c := client.New(session, rpcClient)
		go s.sessionLoop(c)
	}
}

func (s *Server) SDHeart() {
	work := etcd.NewWorker(conf.Conf.Etcd.Name, conf.Conf.Server.Addr, conf.Conf.Etcd.Root, conf.Conf.Etcd.Addrs)
	go work.HeartBeat()
}
