# go-tools-archive-to-db

List files in directory and write metadata to database

GOOS=darwin GOARCH=arm64 CGO_ENABLED=1 go build -o build/go-tools-archive-to-db

### from mac to win

brew install mingw-w64
x86_64-w64-mingw32-gcc --version
which x86_64-w64-mingw32-gcc
export CC_FOR_TARGET=/opt/homebrew/bin/x86_64-w64-mingw32-gcc
GOOS=windows GOARCH=amd64 CGO_ENABLED=1 CC=x86_64-w64-mingw32-gcc go build -o build/go-tools-archive-to-db.exe

set GOOS=windows ; set GOARCH=amd64 ; set CGO_ENABLED=1 ; go build -o build/go-tools-archive-to-db.exe

# todos
