package main

import (
	"fmt"
	"github.com/sagoo-cloud/sagooiot/extend/model"
	"net/rpc"
	"strings"

	gplugin "github.com/hashicorp/go-plugin"
	pluginModule "github.com/sagoo-cloud/sagooiot/extend/module"
)

// ProtocolTgn52 实现
type ProtocolTgn52 struct{}

func (ProtocolTgn52) Info() model.ModuleInfo {
	var res = model.ModuleInfo{}
	res.Name = "tgn52"
	res.Title = "TG-N5 v2设备协议"
	res.Author = "Microrain"
	res.Intro = "对TG-N5插座设备进行数据采集v2"
	res.Version = "0.01"
	return res
}

func (ProtocolTgn52) Encode(args interface{}) (string, error) {
	fmt.Println("接收到参数：", args)
	return "", nil
}

func (ProtocolTgn52) Decode(data []byte, dataIdent string) (string, error) {
	tmpData := strings.Split(string(data), ";")
	var rd = DeviceData{}
	l := len(tmpData)
	if l > 7 {
		rd.HeadStr = tmpData[0]
		rd.DeviceID = tmpData[1]
		rd.Signal = tmpData[2]
		rd.Battery = tmpData[3]
		rd.Temperature = tmpData[4]
		rd.Humidity = tmpData[5]
		rd.Cycle = tmpData[6]
		//处理续传数据
		for i := 7; i < l; i++ {
			rd.Update = append(rd.Update, tmpData[i])
		}
	}
	res := pluginModule.OutJsonRes(0, "", rd)
	if rd.IsEmpty() {
		res = pluginModule.OutJsonRes(1, "数据为空，或数据结构不对", nil)
	}
	return res, nil
}

// Tgn52Plugin 插件接口实现
// 这有两种方法：服务器必须为此插件返回RPC服务器类型。我们为此构建了一个RPCServer。
// 客户端必须返回我们的接口的实现通过RPC客户端。我们为此返回RPC。
type Tgn52Plugin struct{}

// Server 此方法由插件进程延迟调
func (Tgn52Plugin) Server(*gplugin.MuxBroker) (interface{}, error) {
	return &pluginModule.ProtocolRPCServer{Impl: new(ProtocolTgn52)}, nil
}

// Client 此方法由宿主进程调用
func (Tgn52Plugin) Client(b *gplugin.MuxBroker, c *rpc.Client) (interface{}, error) {
	return &pluginModule.ProtocolRPC{Client: c}, nil
}

func main() {
	//调用plugin.Serve()启动侦听，并提供服务
	//ServeConfig 握手配置，插件进程和宿主机进程，都需要保持一致
	gplugin.Serve(&gplugin.ServeConfig{
		HandshakeConfig: pluginModule.HandshakeConfig,
		Plugins:         pluginMap,
	})
}

// 插件进程必须指定Impl，此处赋值为greeter对象
var pluginMap = map[string]gplugin.Plugin{
	"tgn52": new(Tgn52Plugin),
}
