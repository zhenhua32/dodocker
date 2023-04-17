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

## 编译

```bash
go build .
```

## 编译后的命令行工具

编译后会生成一个 dodocker 的命令行工具, 一些常见的用法如下:

```bash
sudo ./dodocker -h
sudo ./dodocker run --ti --cwd ./ -v ./aa:/aa bash
sudo ./dodocker commit --cwd ./ myimage
```

**--cwd** 是当前工作目录, 应该有 busybox.tar 等文件, 用于解压镜像.

# 好久没更新了, 也不知道当时看到哪里了

https://github.com/orgs/community/discussions/50878

https://docs.github.com/zh/authentication/keeping-your-account-and-data-secure/githubs-ssh-key-fingerprints

https://github.blog/2023-03-23-we-updated-our-rsa-ssh-host-key/

没想到 github 更新了 RSA key, 我说怎么一直克隆不了仓库.

更新步骤如下:

```bash
# 第一步
ssh-keygen -R github.com

# 第二步
curl -L https://api.github.com/meta | jq -r '.ssh_keys | .[]' | sed -e 's/^/github.com /' >> ~/.ssh/known_hosts

# 替换第二步, 手动添加到 ~/.ssh/known_hosts 中
github.com ssh-rsa AAAAB3NzaC1yc2EAAAADAQABAAABgQCj7ndNxQowgcQnjshcLrqPEiiphnt+VTTvDP6mHBL9j1aNUkY4Ue1gvwnGLVlOhGeYrnZaMgRK6+PKCUXaDbC7qtbW8gIkhL7aGCsOr/C56SJMy/BCZfxd1nWzAOxSDPgVsmerOBYfNqltV9/hWCqBywINIR+5dIg6JTJ72pcEpEjcYgXkE2YEFXV1JHnsKgbLWNlhScqb2UmyRkQyytRLtL+38TGxkxCflmO+5Z8CSSNY7GidjMIZ7Q4zMjA2n1nGrlTDkzwDCsw+wqFPGQA179cnfGWOWRVruj16z6XyvxvjJwbz0wQZ75XK5tKSb7FNyeIEs4TT4jk+S4dhPeAUC5y+bDYirYgM4GC7uEnztnZyaVWQ7B381AK4Qdrwt51ZqExKbQpTUNn+EjqoTwvqNj4kqx5QUCI0ThS/YkOxJCXmPUWZbhjpCg56i+2aB6CmK2JGhn57K5mj0MNdBXA4/WnwH6XoPWJzK5Nyu2zB3nAZp+S5hpQs+p1vN1/wsjk=
```

然后就这样吧, 时间到了, 该睡觉了.
