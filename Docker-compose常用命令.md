
### 以 user-service 为例：

* docker-compose ps                                   显示所有容器

* docker-compose up -d user-service                     构建建启动nignx容器

* docker exec -it  user-service  bash           登录到 user-service 容器中
~~~

# tom @ tom-pc in ~/goprojects/src/otw [23:21:57] C:1
$ docker exec -it  user-service  bash
root@bc2307cbf8c4:/# exit;
exit


~~~

* docker-compose down   user-service              删除所有 user-service 容器、镜像

* docker-compose restart user-service                    重新启动 user-service  容器

* docker-compose build user-service                      构建镜像 user-service         

* docker-compose build --no-cache user-service     不带缓存的构建 

* docker-compose logs  user-service                      查看 user-service  的日志 

* docker-compose logs -f user-service                    查看 user-service  的实时日志

* docker-compose start nginx                    启动 user-service 容器

* docker-compose rm nginx                       删除 user-service 容器（删除前需要先关闭容器）

* docker-compose stop nginx                    停止 user-service 容器

