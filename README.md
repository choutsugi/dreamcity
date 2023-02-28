# dreamcity

## 心跳机制

心跳检测：服务端连续3个周期未收到心跳时断开连接。

心跳包：仅包含消息头len，值为0。

```go
// 心跳数据
dataBuff := bytes.NewBuffer([]byte{})
binary.Write(dataBuff, binary.LittleEndian, uint32(0))
hbData := dataBuff.Bytes()
```

## 序列化

默认使用gogoprotobuf，安装：

```bash
go install github.com/gogo/protobuf/protoc-gen-gofast@latest
```

生成文件：

```makefile
protoc:
	cd shared/pb && protoc --gofast_out=. *.proto
```

