package main

import (
	"bytes"
	pb "dreamcity/shared/pb/login"
	"encoding/binary"
	"github.com/gogo/protobuf/proto"
	"io"
	"log"
	"net"
	"time"
)

const (
	msgHeadLen        = 4 // 消息头长度=消息序列ID长度+消息路由ID长度+消息体长度
	msgHeadSeqIdLen   = 2 // 消息序列ID长度
	msgHeadRouteIdLen = 2 // 消息路由ID长度
)

const (
	routeIdLogin uint16 = 1 // 登录路由ID
)

func main() {
	conn, err := net.Dial("tcp4", "127.0.0.1:3553")
	if err != nil {
		panic(err)
	}
	log.Println("连接成功")

	// 心跳
	go func() {
		dataBuff := bytes.NewBuffer([]byte{})
		_ = binary.Write(dataBuff, binary.LittleEndian, uint32(0))
		hbData := dataBuff.Bytes()
		ticker := time.NewTicker(1 * time.Second)
		defer ticker.Stop()
		for range ticker.C {
			_, err := conn.Write(hbData)
			if err != nil {
				log.Fatal(err)
			}
		}
	}()

	// 读协程
	go func() {
		for {
			// 解析消息头
			var (
				msgHead    []byte
				msgLen     uint32
				msgSeqID   uint16
				msgRouteID uint16
			)
			msgHead = make([]byte, msgHeadLen+msgHeadSeqIdLen+msgHeadRouteIdLen)
			_, _ = io.ReadFull(conn, msgHead)
			reader := bytes.NewReader(msgHead)
			_ = binary.Read(reader, binary.LittleEndian, &msgLen)
			_ = binary.Read(reader, binary.LittleEndian, &msgSeqID)
			_ = binary.Read(reader, binary.LittleEndian, &msgRouteID)
			// 解析消息体
			msgData := make([]byte, msgLen-2-2)
			_, _ = io.ReadFull(conn, msgData)
			// 匹配消息
			switch msgRouteID {
			case routeIdLogin:
				res := &pb.LoginRes{}
				if err := proto.Unmarshal(msgData, res); err != nil {
					log.Printf("failed to unmarshal pb.LoginRes with err, %+v\n", err)
					continue
				}
				log.Printf("登录响应：序列号=%d, 结果={%v}\n", msgSeqID, res)
			default:
				log.Println("未知消息")
			}
		}
	}()

	// 登录
	{
		req := pb.LoginReq{Token: "token"}
		msg, err := proto.Marshal(&req)
		if err != nil {
			panic(err)
		}

		dataBuff := bytes.NewBuffer([]byte{})
		_ = binary.Write(dataBuff, binary.LittleEndian, uint32(len(msg)+msgHeadSeqIdLen+msgHeadRouteIdLen))
		_ = binary.Write(dataBuff, binary.LittleEndian, 1)
		_ = binary.Write(dataBuff, binary.LittleEndian, routeIdLogin)
		_ = binary.Write(dataBuff, binary.LittleEndian, msg)
		pack := dataBuff.Bytes()

		_, err = conn.Write(pack)
		if err != nil {
			panic(err)
		}
		log.Println("发送登录请求")
	}

	select {}
}
