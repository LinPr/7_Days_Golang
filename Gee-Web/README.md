## gee 网络框架

#### 执行流程
1. 自定义一个多路复用器对象 engine,实现 (engine *Engine) ServeHTTP(w http.ResponseWriter, req *http.Request)方法，使其成为handler接口，这样用一个对象来实现接口的好处是，能够让这个对象带一些其他成员（闭包），或者让这个对象能有更对方法来操作自己
2. 将 req response 封装成一个 HTTP 连接事务的 Context，在serverHTTP 方法中统一调用handler, 使用 context 对象作为传入参数
3. 服务器初始化时，通过 addRouter 注册回调函数来指定对应 方法+路由 触发的行为，其中包含两部分内容：生成前缀树(先序插入书的节点) 和 回调函数集合
4. 服务器启动后，服务端通过路由查找前缀树匹配节点，由ServeHTTP方法来调用对应 method+path 注册好的回调函数