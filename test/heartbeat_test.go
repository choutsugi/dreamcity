package test

import (
	"bytes"
	"encoding/binary"
	"fmt"
	"log"
	"net"
	"testing"
	"time"
)

func TestHeartBeat(t *testing.T) {

	// 创建连接
	conn, err := net.Dial("tcp4", "127.0.0.1:3553")
	if err != nil {
		panic(err)
	}
	log.Println("连接成功", time.Now().Format("2006-01-02 15:04:05.000000"))

	// 心跳协程
	go func() {
		// 心跳数据
		dataBuff := bytes.NewBuffer([]byte{})
		_ = binary.Write(dataBuff, binary.LittleEndian, uint32(0))
		hbData := dataBuff.Bytes()
		// 10次心跳
		for i := 0; i < 10; i++ {
			_, err := conn.Write(hbData)
			if err != nil {
				log.Printf("写数据失败, err: %+v\n", err)
				fmt.Println(err)
				break
			}
			log.Println("发送心跳成功", time.Now().Format("2006-01-02 15:04:05.000000"))
			time.Sleep(time.Second)
		}
	}()

	// 读协程
	for {
		buf := make([]byte, 512)
		n, err := conn.Read(buf)
		if err != nil {
			log.Printf("读数据失败, err: %+v\n", err)
			break
		}
		log.Printf("读数据成功, 长度=%d, 数据=%+v\n", n, buf[:n])
	}

	// 断开连接
	fmt.Println("连接断开", time.Now().Format("2006-01-02 15:04:05.000000"))
}
