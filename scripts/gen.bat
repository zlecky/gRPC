@echo off
setlocal EnableExtensions

set "PROTOC=%LOCALAPPDATA%\protoc\bin\protoc.exe"
set "PATH=%LOCALAPPDATA%\protoc\bin;%GOPATH%\bin;%PATH%"
if defined GOPATH (
  set "PATH=%GOPATH%\bin;%PATH%"
) else (
  for /f "delims=" %%i in ('go env GOPATH') do set "PATH=%%i\bin;%PATH%"
)

if not exist "%PROTOC%" (
  echo protoc not found at %PROTOC%
  exit /b 1
)

if not exist gen\user\v1 mkdir gen\user\v1

"%PROTOC%" -I api ^
  --go_out=gen --go_opt=paths=source_relative ^
  --go-grpc_out=gen --go-grpc_opt=paths=source_relative ^
  user/v1/user.proto

if errorlevel 1 exit /b 1

echo protobuf generated under gen\user\v1
endlocal
