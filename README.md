# RDBFS
Corrently support darwin kernel (macOS) with PostgreSQL

# build
1. install [FUSE](https://osxfuse.github.io/)
2. run `/Library/Filesystems/osxfuse.fs/Contents/Resources/load_osxfuse`
3. brew install go
4. add the following to your .bash_profile
```sh
# GOLANG
export GOPATH="/Users/USERNAME/Go"
export GOROOT="/usr/local/opt/go/libexec"
export PATH="$GOPATH/bin:$GOROOT/bin:$PATH"
```
5. brew install glide

# start postgres docker
1. docker run -d -p 5432:5432 postgres:10-alpine

# clone this project
```sh
mkdir -p ~/Go/src/github.com/rueian
cd ~/Go/src/github.com/rueian
git clone git@github.com:rueian/rdbfs.git
glide install

```

# mount
1. mkdir ~/mnt
2. go run main.go --db-driver "postgres" --db-url "user=postgres sslmode=disable password=mypassword dbname=postgres" start ~/mnt
3. Open other terminal, ls ~/mnt

# unmount
umount ~/mnt
