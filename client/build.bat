@echo off
echo === [ Building ] ===
go build -ldflags="-s -w -H windowsgui" main.go
echo.
echo === [ Compressing ] ===
echo.
.\buildtool\upx.exe -4 main.exe
echo.
echo === [ Moving Artifacts ] ===
echo.
move /-Y .\main.exe .\artifacts\main-%date:~0,4%%date:~5,2%%date:~8,2%%time:~0,2%%time:~3,2%%time:~6,2%.exe
echo.
echo. === [ Finished ] ===
pause