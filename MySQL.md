## 学习大纲




## sql 优化



## 索引使用原则
 
* 使用索引：最左前缀匹配、
* 不能使用索引：前置模糊查询不能使用、反向查找not in、
* 使用join时:可以先where过滤一下，再join，再使用索引
* 字段长度尽量短
* 单表建立 主键单例索引
* 多表链接：联合索引
* 根据业务场景，建立覆盖索引
* where、分组、排序建立索引
