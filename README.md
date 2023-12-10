# gobase  基于go的原生SDK封装的常用功能与框架（为了包的简洁性，完全不引用第三方包），持续保持更新
## 项目主要分为 frame 与 util 两个板块，基本上所有的封装包都有测试用例 。

## frame ： 主要是一些项目开发中的常用框架封装，具何如下：
### mybind : 功能是：将两个对象A[监听],B[被监听] 建立一个绑定关系,则B 的数据变更将会被A实时获取 。
### mybus  : 基于 订阅者与发布者模型封装 
### myctrl : 提供 延时执时，定时执时 与 并发数控制执行的接口 
### mydata : 提供并发安全的数据操作模型 。
### myfuture : 仿flutter中future函数封装 。
### mysignal : 基于基于QT 的信号 与 槽 的思路来封装 ，对象间通信号，可以定义多个信号,对象间的通信则通过信号与槽建立联接来实现通信。
### mysql   ： SQL语句构造器 。


## util  基础功能库封装 ，部分常用库介绍
### mycache: 封装并发安全的 数据缓存功能
### mych ： 基于生产者，消费者功能封装
### mylist : 封装并发安全的数组功能
### mymap : 封装并发安全的数组功能
### myvalidator :封装常用的数据模型校验功能
#### 支持格式如下：
#### enum :针对数组，字符串
#### valid: 针对 结构体，map, 切片 的递归遍历标识
#### 数值类：
#### min;max;range
#### 字符串类：
#### min_len/min_length ;max_len/max_length;range_len/range_length
#### 数组类：
#### arr_minlen/arr_minlength ;arr_maxlen/arr_maxlength;arr_rangelen/arr_rangelength
#### map类：
#### map_minlen/map_minlength ;map_maxlen/map_maxlength;map_rangelen/map_rangelength
