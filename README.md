# go-common

Go 项目开发公共模块， 包含了一组可重用的组件库, 不包含任何业务模块。

- 可以快速的初始化一个包含数据库连接的 WEB 服务器
- 提供 Excel 数据快速导出工具
- 提供数据校验
- 提供 aes 加解密
- 提供基础日志工具 
- 提供 Systemd 脚本安装工具
- 其他一些常用工具函数


## 使用

    go get -u github.com/ca17/go-common
    
    import "github.com/ca17/go-common"