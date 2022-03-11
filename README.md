# Gokins

一款基于Go语言实现的快速部署 Jenkins 服务的命令行工具。

## 安装

### Linux

```bash
sudo wget https://github.com/hxulin/gokins/releases/download/v2.0.0/gokins-2.0.0-linux-amd64 -O /usr/local/bin/gokins
sudo chmod +x /usr/local/bin/gokins
```

### OS X

```bash
sudo curl -Lo /usr/local/bin/gokins https://github.com/hxulin/gokins/releases/download/v2.0.0/gokins-2.0.0-darwin-amd64
sudo chmod +x /usr/local/bin/gokins
```

### Windows

可通过浏览器等工具下载获得可执行文件 `gokins.exe`，建议存放到系统环境变量目录。

如：`C:\Windows`

## 快速使用

### 设置用户授权信息

```bash
gokins username admin
gokins token 11aa09e6009204c2fbd84826784999dab5
```

### 加载样例配置文件

```bash
gokins load http://cdn.huangxulin.cn/gokins/config.yaml
```

加载的配置信息会保存到 `$Home/.gokins/config.yaml` 文件中，请根据实际需求修改配置参数。

### 查看任务列表信息

```bash
gokins ls
```

### 执行部署操作

```bash
gokins run 1005
```

