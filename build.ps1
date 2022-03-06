# Build for linux deployment

$Env:GOOS = "linux"
# $Env:GOOS = "windows"
$Env:GOARCH = "amd64"

# Set GO Path
$cwd = $pwd.Path
cd ../..
$Env:GOPATH = $pwd.Path
cd $cwd
echo "GOPATH set"

# build the front end
npm run build
cp .\public\ .\www\jsCanvas -Force -Recurse
echo "Copied public/ to www/"

go build .\main.go; mv -Force .\main Z:\mediaserver
