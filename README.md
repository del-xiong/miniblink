# blink
本库fork自 https://github.com/raintean/blink 并做了一些更新
使用html来编写golang的GUI程序(only windows), 基于[miniblink开源库](https://github.com/weolar/miniblink49)  

## Demo
[Demo项目地址](https://github.com/del-xiong/miniblink-example)

## 安装
```bash
go get github.com/del-xiong/miniblink
```

## 快速入门教程

- [基本使用](#base-init)
- [go与web页面通过js交互](#js-inject)
- [设置访问网址白名单/黑名单](#set-block-white-list)
- [任意篡改返回的数据](#data-hijack)

---
- [x] 一个可执行文件, miniblink的dll被嵌入库中
- [x] 生成的可执行文件灰常小,15M左右,upx后 12M左右
- [x] 支持无缝golang和浏览器页面js的交互 (Date类型都做了处理), 并支持异步调用golang中的方法(Promise), 支持异常获取.
- [x] 嵌入开发者工具(bdebug构建tags开启)
- [x] 支持虚拟文件系统, 基于golang的http.FileSystem, 意味着go-bindata出的资源可以直接嵌入程序, 无需开启额外的http服务
- [x] 添加了部分简单的接口(最大化,最小化,无任务栏图标等)
- [x] 设置窗口图标(参见icon.go文件)
- [ ] 支持文件拖拽
- [ ] 自定义dll,而不是使用内嵌的dll(防止更新不及时)
- [ ] golang调用js方法时的异步.
- [ ] dll的内存加载, 尝试过基于MemoryModule的方案, 没有成功, 目前是释放dll到临时目录, 再加载.
- [ ] 还有很多...

<h2 id="base-init">基本使用</h2>

```
package main

import (
    "github.com/del-xiong/miniblink"
    "log"
)

func main() {
    //设置调试模式
    miniblink.SetDebugMode(true)
    //初始化miniblink模块
    err := miniblink.InitBlink()
    if err != nil {
        log.Fatal(err)
    }
    // 启动1366x920普通浏览器
    view := miniblink.NewWebView(false, 1366, 920)
    // 启动1366x920透明浏览器(只有web界面会显示)
    //view := miniblink.NewWebView(true, 1366, 920)
    // 加载百度
    view.LoadURL("https://github.com/del-xiong/miniblink")
    // 设置窗体标题(会被web页面标题覆盖)
    view.SetWindowTitle("miniblink window")
    // 移动到屏幕中心位置
    view.MoveToCenter()
    // 显示窗口
    view.ShowWindow()
    // 开启调试模式(会调起chrome调试页面)
    view.ShowDevTools()
}

```

<h2 id="js-inject">go与web页面通过js交互</h2>

#### go中向web注入函数和值

```
# 调用函数 func (view *WebView) Inject(key string, value interface{})
# 例: 注入一个js函数 可通过js控制浏览器自身的位置 (加一个自定义限制 x不能超过1000否则抛异常)
# main.go
view.Inject("MoveWindow", func(x, y int32, relative bool) (int, error) {
    rectx, _ := view.GetWindowRect()
    if (relative && x+rectx > 1000) || (!relative && x > 1000) {
        return 0, errors.New("x位置不能超过1000")
    }
    view.Move(x, y, relative)
    time.Sleep(time.Second)
    return int(time.Now().Unix()), nil
})
# 向JS注入一个值变量
view.Inject("typeId", 456)
# main.html/main.js
# 获取typeId值
console.log(BlinkData.typeId); // 456
# 将浏览器窗口向右移动10px 向下移动5px
await BlinkFunc.MoveWindow(10,5,true)
# 将浏览器窗口移动到屏幕x=100 y=50的位置
await BlinkFunc.MoveWindow(10,5,false)
# 等待go执行完毕并获取返回值/错误
BlinkFunc.MoveWindow(5, 5, true).then(function(val) {
    console.error(val);
    // 其他回调处理
}).
catch (function(err) {
    // 异常捕获
    console.error(err);
});
```
#### golang调用/获取javascript中的方法或者值,异常可捕获(err变量返回)

```
# main.html/main.js
var getLocation = function() {
    return window.location;
}
# main.go
# js函数初始化之后 go可以这样调用并获取返回值
value, err := view.Invoke("getLocation")
if err != nil {
    log.Println(err)
} else {
    // 获取window.location对象并转为map然后获取hostname
    log.Println(value.GetInterface().(map[string]interface{})["hostname"])
}

```


<h2 id="set-block-white-list">设置访问白名单/黑名单</h2>

设置后只能访问指定白名单url/不能访问黑名单url 暂不支持模式匹配  
白名单优先级高于黑名单，如果设置了白名单且未url不在白名单中，无论是否匹配黑名单都会拒绝访问  
注意重复调用SetWhitelist/SetBlacklist不会重置之前的url，而是将新的url添加到黑白名单

```
// 只允许访问百度和google的robots
view.SetWhitelist([]string{"https://www.baidu.com","https://www.google.com/robots.txt"}...)
// 禁止访问头条
view.SetBlacklist([]string{"https://toutiao.com"}...)
// 清空黑名单
view.ClearBlacklist()
// 从黑名单中移除指定规则
view.RemoveFromBlacklist("https://toutiao.com")

```


<h2 id="data-hijack">任意篡改返回的数据</h2>

该功能可用于修改浏览器从网络接收的数据，也可以用于设置数据返回回调(例如监控到返回敏感关键词就报警)    
修改之前需要先定义需要监控的数据类型mime字符串，然后定义数据修改回调函数即可  
修改数据会对性能有影响，请尽量缩小数据类型范围(例如只需要修改json的数据就不要把js/html也加入监听)  
目前支持的数据类型-mime 字符串对应表
```
html: text/html
js: application/x-javascript
css: text/css
json: application/json
svg: image/svg+xml
```
使用方法
```
// 定义数据监听类型 可定义多个
view.AddUrlEndHandlerMimeTypes([]string{"text/html"}...)
// 将网页中的百度换成谷歌 可自己控制要处理的网址
view.SetUrlEndHandler(func(view *miniblink.WebView, mime string, url string, content []byte) []byte {
    log.Println(mime)
    if strings.Contains(url, ".baidu.com") {
        content = bytes.Replace(content,[]byte("百度"),[]byte("谷歌"),-1)
    }
    return content
})

```
<img align="right" width="100%" src="https://raw.githubusercontent.com/del-xiong/miniblink-example/master/static/baigoogledu.jpg">

## 更多 待续...

## 注意
- 网页调试工具默认不打包进可执行文件,请启用BuildTags **bdebug**, eg. `go build -tags bdebug`
- 使用本库需依赖cgo编译环境(mingw32)

## ...
再次感谢miniblink项目, 另外如果觉得本项目好用请点个星.  
欢迎PR, > o <
