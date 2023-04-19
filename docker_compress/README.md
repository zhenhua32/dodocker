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
docker container rm test
# 用 docker load 报错了
docker load -i test.tar
docker import test.tar test:1.1

# 这样反而大小更大了. 而且加载不了
docker save -o test2.tar test:1
docker import test2.tar test:1.2

docker run --rm -it test:1.1 bash
# docker: Error response from daemon: failed to create shim task: OCI runtime create failed: runc create failed: unable to start container process: exec: "bash": executable file not found in $PATH: unknown.
docker run --rm -it test:1.2 bash
docker run --rm -it test:1.3 env

docker create --name test --env-file env test:1
docker import test.tar test:1.3

# 试试别的库
docker-squash -f 10 -t test:1.4 test:1
```

我想查看下目录下文件的大小, 查了一堆 powershell 文档, 没找到把文件大小可读化的, 突然看到一个

```powershell
wsl ls -lh
```

这是何种的套娃.


```
你好，这是Bing。根据我从网络上搜索到的信息，docker export、import、save、load这四个命令的区别和联系如下：

- docker export和docker save都是用来导出镜像或容器的，但是docker export导出的是容器，而docker save导出的是镜像¹²³⁴。
- docker export会丢弃镜像或容器的历史记录和元数据，只保存当前的快照状态，而docker save会保留完整的历史记录和元数据¹²³⁴。
- docker load和docker import都是用来载入镜像或容器的，但是docker load载入的是镜像包，而docker import载入的是容器包¹²³⁴。
- docker load不能对载入的镜像重命名，而docker import可以为载入的镜像指定新名称和标签¹²³⁴。

希望这些信息对你有帮助。😊

源: 与必应的对话， 2023/4/19(1) Docker之save、load、export、import命令 - 小黄鱼的圈子 - 博客园. https://www.cnblogs.com/yuhuang/p/12575331.html 访问时间 2023/4/19.
(2) Docker镜像与容器备份迁移（export、import与commit、save、load） - 腾讯云开发者社区-腾讯云. https://cloud.tencent.com/developer/article/2027894 访问时间 2023/4/19.
(3) docker export、import、save、load 区别 - CSDN博客. https://blog.csdn.net/guokexiaohao/article/details/127454970 访问时间 2023/4/19.
(4) Docker commit、Docker export/import、Docker save/load的区别和适用场景_docker .... https://blog.csdn.net/Wqr_18390921824/article/details/123083151 访问时间 2023/4/19.
```

