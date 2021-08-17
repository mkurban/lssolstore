#!/bin/bash

protoc --proto_path=.																						\
	--go_out=.																										\
	--go_opt=Mlssproto/lssproto.proto=github.com/mkurban/lssstore/lss	\
	--go_opt=Mldproto/ldproto.proto=github.com/mkurban/lssstore/lss 	\
	--go_opt=module=github.com/mkurban/lssstore										\
	--go-grpc_out=.																								\
	--go-grpc_opt=Mldproto/ldproto.proto=github.com/mkurban/lssstore/lss	\
	--go-grpc_opt=Mlssproto/lssproto.proto=github.com/mkurban/lssstore/lss	\
	--go-grpc_opt=module=github.com/mkurban/lssstore							\
	lssproto/lssproto.proto
