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

# SQL Code
`
 SELECT
	*,
 	ROW_NUMBER() OVER (PARTITION BY join_id ORDER BY t1.id) AS rn
 from (
 select
 	t1.id,
 	t1.size,
 	t1.name as t1name,
 	t2.name as t2name,
 	t1.path as t1path,
 	t2.path as t2path,
 	t1.path||t2.path as join_id
 from
 	timotewb t1
 	inner join timotewb t2
 		on(t1.extension = t2.extension
 			and t1.size = t2.size
 			and t1.path <> t2.path)
 where
 	t1.name not in('._.DS_Store', '.DS_Store')
 ) t1
 order by
 	t1.size desc
`


