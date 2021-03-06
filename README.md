# RDBFS
Using FUSE interface to build a filesystem on Relational Database.

# Features
Full FUSE implementation:

All normal file system tools are supported. (mv, rm, cp, gzip, etc…)

Support MySQL and PostgreSQL as backend database.
Take full advantage of relational database, for examples:
* Multiple devices can mount the filesystem at the same time via the Internet.
* Point in Time Recovery by saving snapshots and write ahead logs of database.
* Scale up read throughput by adding more replication slaves.
* Data redundancy with replication slaves for High Availability.


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
