@echo off

protoc --go_out=. --go_opt=paths=source_relative *.proto

echo Protobuf generation completed.