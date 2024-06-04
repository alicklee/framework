package datapack

import (
	"io"
	"net"
	"testing"

	"github.com/CloudcadeSF/shop-heroes-legends-common/log"
	message2 "github.com/CloudcadeSF/Framework/impl/message"
)

func TestDataPack(t *testing.T) {
	/*
		启动一个服务器
	*/

	//1. 创建socketTCP
	listener, err := net.Listen("tcp", "127.0.0.1:1111")
	if err != nil {
		log.Infoln("TCP listener err", err)
		return
	}

	go func() {
		for {
			conn, err := listener.Accept()
			if err != nil {
				log.Infoln("accept err", err)
			}
			go func(conn net.Conn) {
				//定义一个
				dp := NewDataPack()
				//--------------------开始拆包--------------------//
				for {
					//第一次从conn把包的head读取出来
					headData := make([]byte, dp.GetHeadLen())
					_, err := io.ReadFull(conn, headData)
					if err != nil {
						break
					}
					msgHead, err := dp.Unpack(headData)
					if err != nil {
						log.Infoln(err)
						return
					}
					if msgHead.GetMsgLen() > 0 {
						//msg有数据开始第二次读取
						// 第二次从包里把head里面的dataLen读取出来
						msg := msgHead.(*message2.Message)
						msg.Data = make([]byte, msg.GetMsgLen())

						//根绝dataLen再次从conn的数据流中读取
						_, err := io.ReadFull(conn, msg.Data)
						if err != nil {
							log.Infoln("read msg data err", err)
						}

						log.Infoln("Recv MsgID: ", msg.Id, " Msg length: ", msg.DataLen, " Msg Data: ", string(msg.Data))
					}
				}
			}(conn)

		}
	}()

	/*
		模拟一个客户端
	*/
	conn, err := net.Dial("tcp", "127.0.0.1:1111")
	if err != nil {
		log.Infoln("client start err", err)
		return
	}

	dp := NewDataPack()

	//封装2个message的包，粘包发送过去
	msg1 := &message2.Message{
		Id:      1,
		DataLen: 5,
		Data:    []byte{'h', 'e', 'l', 'l', 'o'},
	}
	sendData1, _ := dp.Pack(msg1)
	msg2 := &message2.Message{
		Id:      2,
		DataLen: 5,
		Data:    []byte{'c', 'l', 'o', 'u', 'd'},
	}
	sendData2, _ := dp.Pack(msg2)

	sendData1 = append(sendData1, sendData2...)
	conn.Write(sendData1)

	select {}
}
