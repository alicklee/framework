package connection

import (
	"net"
	"time"

	"github.com/CloudcadeSF/shop-heroes-legends-common/log"
	"github.com/CloudcadeSF/shop-heroes-legends-common/network/mq"
	response3 "github.com/CloudcadeSF/shop-heroes-legends-common/response"

	"github.com/CloudcadeSF/Framework/consts"
	idatapack "github.com/CloudcadeSF/Framework/iface/IDatapack"
	imessage "github.com/CloudcadeSF/Framework/iface/message"
	"github.com/CloudcadeSF/Framework/iface/response"
	"github.com/CloudcadeSF/Framework/impl/message"
	response2 "github.com/CloudcadeSF/Framework/impl/response"
)

func NewMQConn(target *mq.SendTarget, dataPacker idatapack.IDataPack, channel *mq.Channel) *MQConnection {
	return &MQConnection{
		mqChannel:  channel,
		Target:     target,
		dataPacker: dataPacker,
	}
}

type MQConnection struct {
	pid string
	//mq接收消息的频道
	mqChannel *mq.Channel
	//
	Target     *mq.SendTarget
	dataPacker idatapack.IDataPack
}

func (m *MQConnection) Start() {

}

func (m *MQConnection) Stop() {

}

func (m *MQConnection) Kick() {

}

func (m *MQConnection) GetTCPConnection() *net.TCPConn {
	panic("implement me")
}

func (m *MQConnection) RemoteAddr() net.Addr {
	panic("implement me")
}

func (m *MQConnection) FanoutMsg(msgId int32, data []byte) error {
	panic("implement me")
}

func (m *MQConnection) GetPid() string {
	return m.pid
}

func (m *MQConnection) SetPid(pid string) {
	m.pid = pid
}

func (m *MQConnection) Request(msgId int32, data []byte, timeout time.Duration) response.IResponse {
	pack, err := m.dataPacker.Pack(message.NewMessage(msgId, data)) //m.dataPacker.Pack(message.NewMessage(msgId, data))
	if err != nil {
		log.Errorf("pack data err %s", err.Error())
		return nil
	}
	b, err := m.SimpleRequest(pack, timeout)
	if err != nil {
		if err != response3.ErrTimeout {
			log.Errorf("mq conn request err %s", err.Error())
		}
		return response2.NewResponse(0, nil, err)
	}
	//log.Infof("mqconn Request get b len %d", len(b))
	respMsg, err := m.dataPacker.UnpackWithOutToken(b)
	if err != nil {
		//log.Panic(err)
		log.Errorf("mq conn unpack data err %s,msgId %d", err.Error(), msgId)
		return response2.NewResponse(0, nil, err)
	}
	//respMsg.SetMsgData(b[datapack.ByteLenData+datapack.ByteLenProto:])
	//log.Infof("respMsg is %+v", respMsg)
	return response2.NewResponse(respMsg.GetMsgId(), respMsg.GetData(), err)
}

func (m *MQConnection) SimpleRequest(data []byte, timeout time.Duration) ([]byte, error) {
	if m.Target == nil || m.mqChannel == nil {
		return nil, ErrTargetOrChannelIsNotSet
	}
	header := make(map[string]interface{})
	if m.pid != "" {
		header[consts.HeaderServerPlayerId] = m.pid
	}
	resp, err := m.mqChannel.Request(m.Target, timeout, data, header)
	return resp, err
}

func (m *MQConnection) Close() error {
	return nil
}

func (m *MQConnection) GetConnId() uint64 {
	if m.Target == nil {
		return 0
	}
	return m.Target.GetSessionId()
}

var ErrTargetOrChannelIsNotSet = response3.NewRespErr(10401, "target or channel is not set")

func (m *MQConnection) SendProxy(data []byte) error {
	_, err := m.SimpleSendMsg("", data)
	return err
}

func (m *MQConnection) SendMsgImmediately(msgId int32, data []byte) error {
	return m.SendMsg(msgId, data)
}

func (m *MQConnection) SendMsg(mid int32, data []byte) error {
	if m.Target == nil || m.mqChannel == nil {
		return ErrTargetOrChannelIsNotSet
	}
	pack, err := m.dataPacker.Pack(message.NewMessage(mid, data))
	if err != nil {
		return err
	}

	_, err = m.SimpleSendMsg("", pack)
	return err
}

func (m *MQConnection) SendMsgWithCID(cid string, mid int32, data []byte) (string, error) {
	if m.Target == nil || m.mqChannel == nil {
		return "", ErrTargetOrChannelIsNotSet
	}
	pack, err := m.dataPacker.Pack(message.NewMessage(mid, data))
	if err != nil {
		return "", err
	}

	return m.SimpleSendMsg(cid, pack)
}

func (m *MQConnection) SendMsgList(messages ...imessage.IMessage) error {
	if m.Target == nil || m.mqChannel == nil {
		return ErrTargetOrChannelIsNotSet
	}
	for _, msg := range messages {
		if err := m.SendMsg(msg.GetMsgId(), msg.GetData()); err != nil {
			return err
		}
	}
	return nil
}

func (m *MQConnection) SimpleSendMsg(cid string, data []byte) (string, error) {
	if m.Target == nil || m.mqChannel == nil {
		return cid, ErrTargetOrChannelIsNotSet
	}
	header := make(map[string]interface{})
	if m.pid != "" {
		header[consts.HeaderServerPlayerId] = m.pid
	}

	return m.mqChannel.SendNft(cid, m.Target, data, header)
}

func (m *MQConnection) SimpleSendMsgList(datas ...[]byte) error {
	if m.Target == nil || m.mqChannel == nil {
		return ErrTargetOrChannelIsNotSet
	}
	for _, data := range datas {
		if _, err := m.SimpleSendMsg("", data); err != nil {
			return err
		}
	}
	return nil
}
