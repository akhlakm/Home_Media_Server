# Build for linux deployment

$Env:GOOS = "linux"
$Env:GOARCH = "amd64"

# Set GO Path
$cwd = $pwd.Path
cd ../..
$Env:GOPATH = $pwd.Path
cd $cwd

echo "Done"

# go build .\main.go; mv -Force .\main Z:\mediaindexer