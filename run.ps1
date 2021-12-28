$Command = 'cmd /c start powershell -NoExit -Command {
    $host.UI.RawUI.WindowTitle = "Server: 9000 - Chitty-Chat";
    $host.UI.RawUI.BackgroundColor = "black";
    $host.UI.RawUI.ForegroundColor = "white";
    Clear-Host;
    cd server; 
    go run .\server.go 9000;
}'

invoke-expression -Command $Command;

$Command = 'cmd /c start powershell -NoExit -Command {
    $host.UI.RawUI.WindowTitle = "Server: 9001 - Chitty-Chat";
    $host.UI.RawUI.BackgroundColor = "black";
    $host.UI.RawUI.ForegroundColor = "white";
    Clear-Host;
    cd server; 
    go run .\server.go 9001;
}'

invoke-expression -Command $Command;

$Command = 'cmd /c start powershell -NoExit -Command {
    $host.UI.RawUI.WindowTitle = "Server: 9002 - Chitty-Chat";
    $host.UI.RawUI.BackgroundColor = "black";
    $host.UI.RawUI.ForegroundColor = "white";
    Clear-Host;
    cd server; 
    go run .\server.go 9002;
}'

invoke-expression -Command $Command;

$Command = 'cmd /c start powershell -NoExit -Command {
    $host.UI.RawUI.WindowTitle = "Client: Ahmed (5000)- Chitty-Chat";
    $host.UI.RawUI.BackgroundColor = "black";
    $host.UI.RawUI.ForegroundColor = "white";
    Clear-Host;
    cd client;
    go run .\client.go 5000 ahmed;
}'

invoke-expression -Command $Command;

$Command = 'cmd /c start powershell -NoExit -Command {
    $host.UI.RawUI.WindowTitle = "Client: Ali (5001) - Chitty-Chat";
    $host.UI.RawUI.BackgroundColor = "black";
    $host.UI.RawUI.ForegroundColor = "white";
    Clear-Host;
    cd client;
    go run .\client.go 5001 ali;
}'

invoke-expression -Command $Command;
