# 关于

监控平台

# 测试用户

```
dengyang.pan@aliyun.com:4fc27a29
```

# Config

```
ln -s test working
ln -s production working
```

> 切换环境

```
cd env/working
go run gen_key.go
```

> 在当前环境生成密钥对

```
vi env/working/config.tml
```

> 修改当前环境的数据库等配置信息

# Compile

```
make clean
make
```

# Build

```
make clean
make
make dist
```

# Run

```
./accountd --ip 127.0.0.1 --port 9005 &
go run main.go --ip 127.0.0.1 --port 10001
go run main.go --ip 192.168.50.128 --port 10001
```
