# go-leaf
一个使用kratos框架的Snowflake生成器

# 简介
本项目是Snowflake的go版本实现，其中提供三种模式：
* 模式1：正常模式，出现时钟回滚会报错，此模式下
## 雪花算法模式
  Snowflake，Twitter开源的一种分布式ID生成算法。基于64位数实现，下图为Snowflake算法的ID构成图。
  ![](https://p0.meituan.net/travelcube/96034f8fa0f2cb14c21844a4fa12f50441574.png)
  * 第1位置为0。
  * 第2-42位是相对时间戳，通过当前时间戳减去一个固定的历史时间戳生成。
  * 第43-52位是机器号workerID，每个Server的机器ID不同。
  * 第53-64位是自增ID。
  这样通过时间+机器号+自增ID的组合来实现了完全分布式的ID下发。

# 任务列表
- [ ] 雪花算法实现
