# go-leaf
美团分布式ID-leaf的go版本

# 简介
本项目可以认为是美团leaf的go版本实现。  
美团leaf 实现号段模式和雪花算法模式，并且给与了HTTP请求的服务示例
[美团leaf](https://github.com/Meituan-Dianping/Leaf)
[美团leaf介绍](https://tech.meituan.com/2019/03/07/open-source-project-leaf.html)  
Leaf在设计之初就秉承着几点要求:
  * 全局唯一，绝对不会出现重复的ID，且ID整体趋势递增。
  * 高可用，服务完全基于分布式架构，即使MySQL宕机，也能容忍一段时间的数据库不可用。
  * 高并发低延时
  * 接入简单，直接通过公司RPC服务或者HTTP调用即可接入。
# 分布式ID生成简介
* 号段模式
  号段模式基于MySQL，采用预取的方式。也即是每次取一个区间范围的数字，用完之后再取。  
  每个Server启动时，都会去DB拿固定长度的ID List。这样就做到了完全基于分布式的架构，同时因为ID是由内存分发，所以也可以做到很高效。接下来是数据持久化问题，Leaf每次去DB拿固定长度的ID List，然后把最大的ID持久化下来，也就是并非每个ID都做持久化，仅仅持久化一批ID中最大的那一个。这个方式有点像游戏里的定期存档功能，只不过存档的是未来某个时间下发给用户的ID，这样极大地减轻了DB持久化的压力。  
  ![模式如下](https://p1.meituan.net/travelcube/210ca1564c70b228ed46f3b33c9bb9b161120.png)
* 雪花算法模式
  Snowflake，Twitter开源的一种分布式ID生成算法。基于64位数实现，下图为Snowflake算法的ID构成图。  
  ![](https://p0.meituan.net/travelcube/96034f8fa0f2cb14c21844a4fa12f50441574.png)
  * 第1位置为0。
  * 第2-42位是相对时间戳，通过当前时间戳减去一个固定的历史时间戳生成。
  * 第43-52位是机器号workerID，每个Server的机器ID不同。
  * 第53-64位是自增ID。
  这样通过时间+机器号+自增ID的组合来实现了完全分布式的ID下发。

# 任务列表
- [ ] 雪花算法实现
- [ ] 基于etcd替换 leaf的zk
- [ ] 号段模式实现
- [ ] 对外提供HTTP服务
- [ ] 支持zk调用
- [ ] 支持RPC调用
- [ ] 提供RPC客户端
