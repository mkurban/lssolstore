#!/bin/bash

protoc --proto_path=.																\
	--go_out=.																							\
	--go_opt=Mldproto/ldproto.proto=github.com/mkurban/lssstore/lss	\
	--go_opt=module=github.com/mkurban/lssstore							\
	ldproto/ldproto.proto
