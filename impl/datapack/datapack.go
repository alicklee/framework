package datapack

import (
	"bytes"
	"encoding/binary"
	"fmt"

	"github.com/CloudcadeSF/Framework/iface/message"
	message2 "github.com/CloudcadeSF/Framework/impl/message"
	"github.com/CloudcadeSF/Framework/utils"
	"github.com/CloudcadeSF/shop-heroes-legends-common/log"
)

/*
*
包体的结构
*/
type DataPack struct {
}

var SimpleDataPacker = NewDataPack()

// 拆包 封包的实例
func NewDataPack() *DataPack {
	return &DataPack{}
}

const (
	//包长度
	ByteLenData = 4
	//协议
	ByteLenProto = 4
	//token长度
	ByteLenToken = 64
)

// 获取包的头的长度
func (d *DataPack) GetHeadLen() int {
	return ByteLenData + ByteLenProto + ByteLenToken
}

// 封包
func (d *DataPack) Pack(msg message.IMessage) ([]byte, error) {
	//创建一个buffer的缓冲
	dataBuff := bytes.NewBuffer([]byte{})

	//将dataLen写入dataBuff中
	if err := binary.Write(dataBuff, binary.BigEndian, msg.GetMsgLen()); err != nil {
		log.Infoln("write message length error", err)
		return nil, err
	}

	//将message的id写入dataBuff中
	if err := binary.Write(dataBuff, binary.BigEndian, msg.GetMsgId()); err != nil {
		log.Infoln("write message ID error", err)
		return nil, err
	}

	//将消息内容写入dataBuff中
	if err := binary.Write(dataBuff, binary.BigEndian, msg.GetData()); err != nil {
		log.Infoln("write message data error", err)
		return nil, err
	}
	return dataBuff.Bytes(), nil
}

/*
Pack with token
*/
func (d *DataPack) PackWithToken(msg message.IMessage, token string) ([]byte, error) {
	//创建一个buffer的缓冲
	dataBuff := bytes.NewBuffer([]byte{})

	//将dataLen写入dataBuff中
	if err := binary.Write(dataBuff, binary.BigEndian, msg.GetMsgLen()); err != nil {
		log.Infoln("write message length error", err)
		return nil, err
	}

	//将message的id写入dataBuff中
	if err := binary.Write(dataBuff, binary.BigEndian, msg.GetMsgId()); err != nil {
		log.Infoln("write message ID error", err)
		return nil, err
	}

	//将token写入dataBuff中
	if err := binary.Write(dataBuff, binary.BigEndian, []byte(token)); err != nil {
		log.Infoln("write token ID error", err)
		return nil, err
	}

	//将消息内容写入dataBuff中
	if err := binary.Write(dataBuff, binary.BigEndian, msg.GetData()); err != nil {
		log.Infoln("write message data error", err)
		return nil, err
	}
	return dataBuff.Bytes(), nil
}

/*
*
拆包
先把包头的长度读取出来，再根据长度读取后续的内容
*/
func (d *DataPack) Unpack(binaryData []byte) (message.IMessage, error) {
	//创建一个从输入二进制数据的ioReader
	dataBuff := bytes.NewReader(binaryData)

	msg := &message2.Message{}
	//读取length
	if err := binary.Read(dataBuff, binary.BigEndian, &msg.DataLen); err != nil {
		return nil, err
	}

	//读取id
	if err := binary.Read(dataBuff, binary.BigEndian, &msg.Id); err != nil {
		return nil, err
	}

	//读取Token
	token := make([]byte, ByteLenToken)
	if err := binary.Read(dataBuff, binary.BigEndian, token); err != nil {
		return nil, err
	}
	msg.SetToken(string(token[:]))
	//判断包体长度是不是超过最大的限制
	if utils.GlobalObject.MaxPackageSize > 0 && utils.GlobalObject.MaxPackageSize < msg.DataLen {
		return nil, fmt.Errorf("message length too long %v %v", msg.DataLen, utils.GlobalObject.MaxPackageSize)
	}
	if len(binaryData) == d.GetHeadLen()+int(msg.DataLen) {
		data := make([]byte, msg.DataLen)
		if err := binary.Read(dataBuff, binary.BigEndian, data); err == nil {
			msg.Data = data
		} else {
			return nil, err
		}
	}
	return msg, nil
}

/*
unpack with token
*/
func (d *DataPack) UnpackWithOutToken(binaryData []byte) (message.IMessage, error) {
	//创建一个从输入二进制数据的ioReader
	dataBuff := bytes.NewReader(binaryData)

	msg := &message2.Message{}
	//读取length
	if err := binary.Read(dataBuff, binary.BigEndian, &msg.DataLen); err != nil {
		return nil, err
	}

	//读取id
	if err := binary.Read(dataBuff, binary.BigEndian, &msg.Id); err != nil {
		return nil, err
	}

	//判断包体长度是不是超过最大的限制
	if utils.GlobalObject.MaxPackageSize > 0 && utils.GlobalObject.MaxPackageSize < msg.DataLen {
		return nil, fmt.Errorf("message length too long %v %v", msg.DataLen, utils.GlobalObject.MaxPackageSize)
	}
	headLen := ByteLenData + ByteLenProto
	if len(binaryData) == headLen+int(msg.DataLen) {
		data := make([]byte, msg.DataLen)
		if err := binary.Read(dataBuff, binary.BigEndian, data); err == nil {
			msg.Data = data
		} else {
			return nil, err
		}
	}
	return msg, nil
}
