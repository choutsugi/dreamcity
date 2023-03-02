protoc:
	cd shared/pb && protoc --gofast_out=../../../ code.proto
	cd shared/pb && protoc --gofast_out=../../../ login.proto
	cd shared/pb && protoc --gofast_out=../../../ world.proto