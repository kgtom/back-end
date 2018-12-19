
## 本节知识

* [一、单个删除](#1)
     * [1.删除镜像](#21)
     * [2.删除容器](#22)
* [二、批量删除](#3)
     * [1.删除镜像](#31)
     * [2.删除容器](#32)
* [三、todo](#4)

## <span id="1">一、单个删除</span>

### <span id="11">删除容器</span>
* 删除容器的命令：docker rm 容器ID或者名字，查看容器ID、名字
~~~

# tom @ tom-pc in ~ [23:43:39]
$ docker container ls -a
CONTAINER ID        IMAGE               COMMAND             CREATED             STATUS                     PORTS                      NAMES
6e3fc01a9cdc        web-service         "/web-service"      8 hours ago         Up 8 hours                 0.0.0.0:8084->8084/tcp     nervous_grothendieck
222c104147c9        u-service           "/u-service"        9 hours ago         Up 9 hours                 0.0.0.0:8083->8083/tcp     elegant_montalcini

~~~

* 如果要删除的 container 处于运行状态，那么先把容器停止了，然后再删除

~~~
docker   stop   ID
~~~


### 删除 image 镜像
