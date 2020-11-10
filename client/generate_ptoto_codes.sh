if [ -d "./generated" ]
then
	rm -rf ./generated/*
else
	mkdir -p ./generated
fi

protoc -I protos/ protos/*.proto --go_out=plugins=grpc:./
