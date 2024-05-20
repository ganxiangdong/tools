# 工具
1. 一些命令行工具
2. mystringer：和stringer类似，但是它会生成IsXX方法，用于判断是否是某个枚举值

### 安装
`go install codeup.aliyun.com/5f9118049cffa29cfdd3be1c/tools/cmd/tools@master`
由于是私有仓库，需要做一些
1. 配置账号到.netrc中，参考[这里](https://help.aliyun.com/document_detail/293638.html?spm=a2cl9.codeup_devops2020_goldlog_projectFiles.0.0.65312897cLAdzK)
2. go env -w GOPRIVATE=codeup.aliyun.com

### 使用
`tools --help`

### 加入到环境变量
如果直接运行 tools 提示命令不存在，表明gopath的bin没有在可执行的 path 中
1. mac
    * `vim ~/.bashrc`或者`vim ~/.zshrc`
    * 加入如下代码：
    `export PATH=$PATH:$GOPATH/bin`
2. windows

### 当前支持的工具
```
tools model -h：根据 ent 的模型生成 ddd 相关代码
tools wire -h：自识别 internal 目录下需要注入的方法，写入到 provider_set.go 文件中
```

### 工具的配置文件tools.yaml
配置到使用工具的项目根目录，工具会读取.tools.yaml文件，支持的配置文件如下：
```yaml
# tools 工具配置文件
wire:
  # 要排除的目录
  exclude:
    - "internal/conf"
```

### 项目结构
```
internal
    biz 用于测试
    data 用于测试
    model 生成model代码的代码
```

### Codegen 生成代码
```
# 安装 v1.3.10 以上版本
go install codeup.aliyun.com/5f9118049cffa29cfdd3be1c/tools/cmd/codegen@v1.3.10
# 生成代码
# path 生成代码的路径 feature commonPage生成page_bo table 生成代码的表名
codegen generate .\schema --path dev --feature commonPage --table platform
```