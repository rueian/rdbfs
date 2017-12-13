# build
1. install [FUSE](https://osxfuse.github.io/)
2. brew install go
3. brew install glide
4. glide install

# start postgres docker
1. docker run -d -p 5432:5432 postgres:10-alpine

# mount
1. mkdir ~/mnt
2. go run main.go --db-driver "postgres" --db-url "user=postgres sslmode=disable password=mypassword dbname=postgres" start ~/mnt
3. Open other terminal, ls ~/mnt
4. umount ~/mnt
