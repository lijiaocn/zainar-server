---
layout: default
title: DevelopDiary

---

# DevelopDiary
创建时间: 2014/08/07 06:03:13  修改时间: 2014/08/07 14:28:55 作者:lijiao

----
## Roadmap

v1.0:

	代号:荷花池塘, 支持中小规模的机构自建使用
	
	指标:
	1. 在线客户端5W以下
	2. 支持单机部署(不推荐)
	3. 支持主备部署: 数据库服务器主备，业务服务器主备(登陆系统、消息系统部署在同一台机器)
	4. 静态文件位于两台业务服务器上通过rsync同步。
	
	方案:
	1. mysql数据库后端
	2. 数据全部写入数据库
	3. Num族长
	4. DB女王

v2.0:

	代号: 鱼米之湖, 支持中等规模的用户
	
	指标:
	1. 在线客户端5W以上,上限待定
	2. 统一部署
	
	方案:
	1. 分布式数据库后端
	2. 数据库+静态文件CDN
	3. KV蚁群
	4. 实时蛛网
	5. 蜡质蛛网

v3.0:

	代号: 万众之海，支持大规模的用户

v5.0:

	代号: 璀璨星空，提供开放式云平台


## 2014-08-07 06:04:42

v1.0版设计思路:

	1. 服务端分为登陆系统和消息系统两部分。 登陆系统负责登陆，消息系统负责消息处理。
	2. 登陆采用https连接保证安全性，登陆验证成功后，为用户分配sid、sw、sp，回复给客户端后，立即断开连接。
	3. 客户端向消息系统发起一条连接用于消息传递。登陆系统依靠附加在每个消息中的sid、num验证消息的归属。
	4. 客户端发送的消息采用json格式。消息系统返回的数据的第一行确定消息的执行状态(OK/ER)以及紧邻的后续的内容的格式。
	5. 消息系统返回的数据中会指明返回的数据的类型, 客户端不必维护发送消息与返回数据的对应关系。
	6. 客户端发送的每条消息附带sid，sid是登陆系统分配给客户端的会话ID，消息系统通过sid从会话中获取发送者身份信息。
	7. 客户端发送的每条消息附带num，每个消息的num都是通过sw和sp计算出来的，客户端与消息系统的计算策略一致。如果消息系统发现消息中的num与自己计算出来的不一致，要求客户端重新登陆。
	8. 客户端与消息系统的通信采用tcp协议，由通信网络保证数据的准确性和传输的可靠性。
	9. 客户端通过订阅的方式获取特定目标的动态(例如team的动态),客户端的连接句柄(conn)被添加到指定目标的听众列表中，目标状态变化时向所有的听众广播变更。
	10. 目标状态每次变更时，会在执行变更的机器上立即触发一次广播。而不在同一台机器上的听众，需要在有所在的机器负责的每5s执行一次轮询中获取最新的变更。
	11. 消息分为实时消息(MsgRT,Message RealTime)和留言板消息(MsgBd,Message Board)。
	    实时消息只会被广播一次，如果客户端没能在广播时进入到听众列表，就永远的错过了这条实时消息。
	    留言板消息如果没有被接收目标确认为"已读",接收目标在每次登陆时都会被告知有多少条留言没有被确认。
	12. 所有的数据都存储在关系型数据库中。
	    (一些静态文件存储，例如头像的存储还没实现,这些静态文件计划存放在文件系统中,并且具有惟一的Uri,以充分利用CDN)

v1.0存在的问题:

	1. 严重依赖数据库，导致数据库写入非常频繁。例如v1.0中会话表、实时消息表、留言表和状态表(例如TeamStat)全部是以关系表的方式存放在数据库中。
	   每收取到一个正确的消息，都至少要更新一个条目。而从书籍、网络上获知，当前的分布是关系数据库的TPS都是几十万量级。
	   假设每个客户端每5s发送一条消息(例如位置变更消息), 那么同时在线客户端数目百万级别都达不到。而且几乎所有的更新都是需要在写入之后，因此实际情况可能更差。
	   (例如，淘宝的Oceanbase中将写操作集中在一台服务器,如果写入后立即去读，新写入的数据很有可能还没有被同步到其它的节点上, 因此读取时必须要与写入服务器上的数据进行比对)

v1.0存在问题的解决思路以及引发的新问题:

	1. v1.0问题1的分析:
	   导致频繁写入的原因有两个:

	   原因1:

	   因为消息系统中对消息的num进行实时变更，而num是存放在会话条目中。因此每收到一个正确的消息，就需要在会话表中变更一次num。
	   这种每次变更方式相比固定的cookie更为安全(而且可以扩展变更策略增加被猜解的难度)，而存放在会话表中的好处是消息服务器可以被随时切换。
	   例如客户端最开始被负载均衡到消息服务器A上，之后客户端的IP变动被负载均衡到了消息服务器B上，不会影响对消息的验证和处理。而且对消息服务器的宕机也有一定的抵抗性。
	   例如消息服务器A宕机时，会话表中num已经被成功更新的客户端只需要重新连接一台消息服务器，而不需要重新登陆获取新的变更策略。

	   方案命名: Num族长，寓意是Num的事情族长说了算。
	   方案结论: 5W以下并发采用

	   解决思路:

	   a. 使用分布式kv管理会话, 分布式kv相比分布式关系数据库要求要少，性能应该要高不少。
	      但是在aws和阿里云上没有找到分布式kv服务, 自己组建显然不现实, pass。

	   b. 将num从会话表条目中剥离，由消息服务器执行变更操作，变更后的num不写回会话表。
	      在单台消息服务器上可以直接使用单机kv管理num，性能提高，并且可以线性扩展。
	      但是使用会话表管理num的好处也随之消失。切换消息服务器或消息服务器宕机时，对应的客户端必须重新登陆。
	      用户体验不好，但是可以容忍。例如，可以在客户端中增加自动登陆的选项，如果用户选择自动登陆，客户端保存用户的账号口令，需要重新登陆时自动登陆。
	      保存账号口令是非常不安全的，注重安全的用户可以不选择这个选项，这样消息服务器被切换时，提示用户网络环境变动，让用户重新登陆。

	      这是一个折衷方案，不够完美，但是可以考虑。

	      方案命名: KV蚁群，寓意是将KV分散到每一个小蚂蚁身上。
	      方案结论: 采用

	  原因2:

	  因为所有的实时消息和留言消息都保存在数据库表中，因此每产生一条，就需要执行一次写入。而且实时消息中的位置更新消息会是很频繁的，用户量如果一旦上去，写入压力非常重。

	  方案命名: DB女王，寓意数据库大包大揽
	  方案结论: 5W以下并发采用

	  解决思路:

	  a. 完全摒弃用数据库表存放实时消息的方式。在分布式数据库中建立一张路由表，每个路由项由客户端标记和所在的消息服务器组成，例如 <ClientX, MsgServerB>。
	     用户连接到消息服务器时，更新路由项。向指定目标发送实时消息时，查找路由表获得目标所在的消息服务器地址，将实时消息发送到对应的消息服务器。
	     每个消息服务器开启一个端口，接收到其它消息服务器之间传递过来的消息后，立即将其广播给自己管理的客户端。

	     这种方式极大的降低数据库的压力。只有在用户登陆到消息服务器时执行一次写入操作。(如果每秒登陆用户数都成为瓶颈了，哈哈，我就可以退休了！！...想多了...)
	     用这种方式需要考虑路由项失效的问题。例如有一个目标是TeamA，这个目标同时位于三台消息服务器上(MsgServerA,MsgServerB,MsgServerC)。其中MsgServerC宕机了,
	     向MsgServerC发送消息的消息服务器会探测到, 可以将这条路由项标记为可疑路由。(具体的路由删除策略还需要考虑)

	     这个方案可以很好的解决实时消息的问题，实现的复杂度完全可以接受，而且消息服务器不再需要轮询目标状态，采用。

	     方案命名: 实时蛛网，寓意是消息服务器之间通过网络连接起来，如同蛛网一般，路由表坐镇蛛网。
	     方案结论: 采用

	 b.  留言板消息被认为是比较重要的消息，必须可靠，并且需要让接收者确认已经接收。
	     产生一条新的留言的时候，首先将留言写入到留言板消息表(MsgBd)，然后查找路由表，向目标发送新的留言。
	     目标只有在登陆时需要读取留言板消息表，检查是否存在未确认的留言。

	     留言板消息必须进行写入操作，最早对留言的定义就是一个低频的操作，并且允许延迟，因为这个功能是留言而不是聊天。(这个系统不计划做聊天...)

	     这个方案只需要在"实时蛛网"上增加一个入库操作, 采用。

	     方案命名: 蜡质蛛网，寓意是在实时蛛网上打了一层蜡，蜡意味着缓慢、固化，消息被固化在那里，等你来取。

## 文献
1. http://xxx  "Name"

