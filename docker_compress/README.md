试试看怎么压缩 docker 镜像.

https://zhuanlan.zhihu.com/p/161685245


# 命令记录

```bash
docker build --no-cache -t test:1 -f Dockerfile1 .
docker build --no-cache -t test:2 -f Dockerfile2 .


# 需要初始化一个容器
docker create --name test test:1
docker container list -a
docker export test -o test.tar
```

我想查看下目录下文件的大小, 查了一堆 powershell 文档, 没找到把文件大小可读化的, 突然看到一个

```powershell
wsl ls -lh
```

这是何种的套娃.
