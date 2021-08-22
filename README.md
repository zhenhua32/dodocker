# dodocker

原始代码来自 `https://github.com/xianlubird/mydocker` 和 <自己动手写docker> 书籍.

## 环境

使用的 Go 版本是 1.16, 系统是

```bash
$ uname -a
Linux tx 4.15.0-143-generic #147-Ubuntu SMP Wed Apr 14 16:10:11 UTC 2021 x86_64 x86_64 x86_64 GNU/Linux
```

## 提取 busybox 文件系统的 rootfs

```bash
docker pull busybox
docker run -d busybox top -b
docker export -o busybox.tar <容器ID>
mkdir -p ./busybox
tar xvf busybox.tar -C busybox/
```
