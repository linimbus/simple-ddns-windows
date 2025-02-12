rsrc -manifest exe.manifest -ico main.ico
go-bindata -o icon_files.go main.ico status.ico stop.ico start.ico
go build -buildvcs=false -ldflags="-H windowsgui -w -s" -o simple-ddns-windows.exe
