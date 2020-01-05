Go Agar
=============

一个简单但功能完备的 Agar (类似 球球大作战）克隆实现，使用Go + HTML5 + es6 编写<br/><br/>
前端代码和后端处理逻辑部分源自`agar.io-clone`项目 https://github.com/huytd/agar.io-clone

---

## 如何开始
访问 http://localhost:38888 进行游戏

#### 基本描述
- 在屏幕上移动鼠标来指定前进方向.
- 吃掉食物、其他比你体积更小的玩家，并避免被其他玩家吃掉.
- 出于对新玩家的保护,如果玩家未开始"进食",你将不会被吃掉,
- 你的质量/体积由你吃下的食物和其他玩家总量决定,并因此影响你的速度.
- 质量将随时间缓慢下降.
- 目标：幸存下来,成为质量最高的玩家.
- 发射孢子(减少质量但增加速度): 按键 'Z'
- 主动分裂: 按键 'X'


---

## 安装
- 需要`go`1.13或以上环境

#### 依赖
- github.com/gin-gonic/gin
- github.com/gorilla/websocket
- github.com/satori/go.uuid
- github.com/spf13/viper

#### 下载依赖
在目录下执行
````
go mod download
````

#### 构建
进入以下目录
````
cmd/go-agar
````

编译(支持交叉编译)

````
go build
````


#### 启动服务
执行编译后的输出文件
````
go-agar
````

游戏将会在这个地址启动 `http://localhost:38888` .默认情况下,端口号为 `38888`,可以在配置文件或代码中更新这个值.

## 配置

#### 支持
配置文件为非必须的,支持名称为`game`的,包括yaml/properties/ini等格式的配置文件

#### 配置项
配置项参考
````
internel/game/config.go
````

## 额外说明
#### 前端说明
前端使用原生es6代码,并未基于任何编译,在部分旧浏览器中可能无法正常使用
#### 网页打包
使用go-bindata打包为go文件,用法如下
````
go-bindata -pkg asset -o internal/asset/asset.go web/...
````
## 许可
项目基于MIT许可