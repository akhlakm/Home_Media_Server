# Run the following after git checkout

# Install the node packages for the frontend
# npm install
# Ignore the node_modules directory for Dropbox
# Set-Content -Stream com.dropbox.ignored -Value 1 -Path node_modules

$Env:GOARCH = "amd64"

# Build for linux deployment
# $Env:GOOS = "windows"
$Env:GOOS = "linux"

# Set GO Path
$cwd = $pwd.Path
cd ../..
$Env:GOPATH = $pwd.Path
cd $cwd
echo "GOPATH set"

# build the front end
npm run build
Remove-Item .\www\jsCanvas\ -Recurse -Force -Confirm:$false
cp .\public\ .\www\jsCanvas -Force -Recurse
echo "Copied public/ to www/"

go build .\main.go; mv -Force .\main Z:\mediaserver

# run with ./mediaserver -serve -walk
