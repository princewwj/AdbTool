# Telnet
Control Android Interface

1.Complie：
go build -ldflags="-s -w" -o telnet.exe main.go

2.Call Format:
telnet.exe "-S" "LocalIP" "RemoteIP" "Port" "UserName" "Password" "Command" "-t" "Timeout(s)"  

3.For Example:
telnet.exe "-S" "192.168.1.100" "192.168.1.1" "23" "root" "password" "ls /tmp" "-t" "38"
