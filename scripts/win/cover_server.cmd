:: Script for creating an application cover

cd ../../cmd/server/server

:: Del old data
IF EXIST cover_server.out DEL cover_server.out /Q /S
IF EXIST cover_server.html DEL cover_server.html /Q /S

:: Code cover
go test -coverprofile cover_server.out
go tool cover -html=cover_server.out -o cover_server.html

:: Run report
start cover_server.html