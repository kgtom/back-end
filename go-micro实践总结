## 本节大纲
* [一、微服务理解](#1)
* [二、micro](#2)
     * [1.工具](#21)
     * [2.使用](#22)
* [三、go-micro](#3)
     * [1.命令行](#31)
     * [2.代码](#32)
* [四、api gateway](#4)
* [五、consul](#4)
* [六、docker-compose](#4)
* [七、k8s](#4)
* [八、elk](#4)
* [九、Zipkin+Prometheus Grafana/elk](#4)

## <span id="1"> 一、微服务理解</span>
微服务的关键理念:业务的拆分，这是从 unix 的设计哲学中得到的启示：
```doing one thing and doing it well```
* 框架：客户端-->api gateway -->多个服务(各个服务之间通过rpc、nsq进行通信解耦)
## <span id="2"> 二.micro</span>
* micro:微服务的工具包
## <span id="3"> 三.go-micro</span>
* go-micro:Go中用于开发微服务的RPC框架。提供了服务发现，客户端负载平衡，编码，同步和异步通信的库。
* micro API：是一个API网关或代理，用于提供HTTP并将请求路由到适当的微服务。
它可以作为单个入口点，用于反向代理或将HTTP请求转换为RPC。
* micro.NewService(...Option) 简化了微服务的注册流程， micro.Run() 简化了服务启动
* 限流器:golang.org/x/time/rate，防御下游，包含后端
* 熔断器：github.com/afex/hystrix-go，抵御上游，保护自身服务。
## <span id="4"> 四、api gateway</span>
  封装请求、减少通信次数；统一鉴权、流控
## <span id="5"> 五、consul</span>
微服务框架，默认consul 服务注册与发现
## <span id="6"> 六.docker-compose</span>
 容器编排，统一管理docker 镜像
## <span id="7"> 七、k8s</span>
将服务部署在pod，使用 service 做pod的负载均衡，使用Kube-DNS将 service 名解析成具体的ClusterIP，对 service实现负载。
## <span id="8"> 八、elk</span>
 elk:log-->es-->kibana
## <span id="9">九、Zipkin+Prometheus Grafana/elk]</span>
可以将 zipkin trace 数据推到 prometheus 监控系统，通过 grafana 可视化；
或者将 zipkin 数据存储放到 es，结合 kibana 生成图表。

