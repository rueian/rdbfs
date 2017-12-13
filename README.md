# build
1. brew install glide
2. install [FUSE](https://osxfuse.github.io/)
3. glide install
4. make clean ; make

# mount
1. mkdir ~/mnt
2. ./build/rdbfs ~/mnt
3. Open other terminal, ls ~/mnt, cat ~/mnt/file.txt
4. umount ~/mnt
