
使用 mount 可以挂载 overlay2.


```bash
mkdir ./dir
mkdir ./dir1
mkdir ./dir2
mkdir ./dir3
sudo mount -t overlay overlay -o lowerdir=./dir1,upperdir=./dir2,workdir=./dir3 ./dir

sudo umount ./dir
```

