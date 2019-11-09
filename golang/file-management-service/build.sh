rm -rf ./plugin
mkdir ./plugin
go build -buildmode=plugin -o ./plugin/FileService.so Plugin.go

