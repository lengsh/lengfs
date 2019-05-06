# lengfs
mini simple distributed web file system

======
run: lengfs -p=80 -i=0 -s=localhost:8080

// static/lengfs/Node/Date/domain/filename.xyz
//    |-----|     |     |    |          |
// Parent  Pnode Inode  |  user-domain  |
//                 Create-date        fileName
//


git clone https://github.com/lengsh/lengfs.git
git add xx.go
git commit -m "init"
git push origin master
