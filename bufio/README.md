# Go语言bufio包源码阅读
bufio包在go语言中使用还是非常多的，这个包的作用是为了将读写操作聚合起到缓冲的作用，减少内存和磁盘（或者其他I/O设备）之间的数据交换。整个包的代码还是比较简洁的，只有两个文件 `bufio.go` 和 `scan.go` ，涉及三种接口类型，如下所示：

* [Writer](writer)
* [Reader](reader)
* [Scanner](scanner)

