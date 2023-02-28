protoc:
	cd shared/pb && protoc --gofast_out=. *.proto