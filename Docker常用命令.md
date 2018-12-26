
## 本节知识

* [一、单个删除](#1)
     * [1.删除镜像](#21)
     * [2.删除容器](#22)
* [二、批量删除](#3)
     * [1.删除镜像](#31)
     * [2.删除容器](#32)
* [三、进入postgres的容器 ](#4)

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

## <span id="3">进入postgres容器</span>
~~~
# tom @ tom-pc in ~ [21:34:05] C:127
$ docker ps
CONTAINER ID        IMAGE                    COMMAND                  CREATED             STATUS                PORTS                                                                                                                       
bc2307cbf8c4        postgres                 "docker-entrypoint.s…"   About an hour ago   Up About an hour      0.0.0.0:5432->5432/tcp                                                                                                                           otw_database_1
19bfd67dba7f        mongo                    "docker-entrypoint.s…"   5 days ago          Up About an hour      0.0.0.0:27017->27017/tcp                                                                                                                         otw_datastore_1
727ad15029ea        progrium/consul:latest   "/bin/start -ui-dir …"   6 days ago          Up 6 days (healthy)   0.0.0.0:53->53/udp, 0.0.0.0:8300-8302->8300-8302/tcp, 0.0.0.0:8400->8400/tcp, 53/tcp, 0.0.0.0:8301-8302->8301-8302/udp, 0.0.0.0:8500->8500/tcp   consul-shop
222c104147c9        u-service                "/u-service"             7 days ago          Up 7 days             0.0.0.0:8083->8083/tcp                                                                                                                           elegant_montalcini

# tom @ tom-pc in ~ [21:34:08]
$ docker exec -it otw_database_1 psql -U postgres -d postgres
psql (11.1 (Debian 11.1-1.pgdg90+1))
Type "help" for help.

postgres=# select * from users;
                  id                  | name | company | email | password | xxx_unrecognized | xxx_sizecache
--------------------------------------+------+---------+-------+----------+------------------+---------------
 7dee641b-0420-4185-aaf3-a33acb48fc5b |      |         |       |          | \x               |             0
 f44cdfa6-0c3b-49a6-856d-cb7414b9d910 |      |         |       |          | \x               |             0
 87573a08-e805-4461-a20d-e14c2bcfeb49 |      |         |       |          | \x               |             0
 953c9d1f-d833-41ee-99c0-cf8a78a21745 |      |         |       |          | \x               |             0
 c3a065ce-87df-42bc-9e00-0789924503fe |      |         |       |          | \x               |             0
 96779bb2-1c6e-4da0-a9fe-7715139adbf7 |      |         |       |          | \x               |             0
 92f1acc5-13c6-411e-9f61-f3e5837754c5 |      |         |       |          | \x               |             0
 a63e0c61-5fd8-4adb-954e-e95b0d4e77c9 |      |         |       |          | \x               |             0
 3c012178-1435-4e49-908b-6c2a001ae3d1 |      |         |       |          | \x               |             0
(9 rows)

postgres=#

~~~
