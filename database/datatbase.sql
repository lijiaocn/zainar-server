-- MySQL dump 10.14  Distrib 5.5.35-MariaDB, for Linux (x86_64)
--
-- Host: localhost    Database: shangwei
-- ------------------------------------------------------
-- Server version	5.5.35-MariaDB-log

/*!40101 SET @OLD_CHARACTER_SET_CLIENT=@@CHARACTER_SET_CLIENT */;
/*!40101 SET @OLD_CHARACTER_SET_RESULTS=@@CHARACTER_SET_RESULTS */;
/*!40101 SET @OLD_COLLATION_CONNECTION=@@COLLATION_CONNECTION */;
/*!40101 SET NAMES utf8 */;
/*!40103 SET @OLD_TIME_ZONE=@@TIME_ZONE */;
/*!40103 SET TIME_ZONE='+00:00' */;
/*!40014 SET @OLD_UNIQUE_CHECKS=@@UNIQUE_CHECKS, UNIQUE_CHECKS=0 */;
/*!40014 SET @OLD_FOREIGN_KEY_CHECKS=@@FOREIGN_KEY_CHECKS, FOREIGN_KEY_CHECKS=0 */;
/*!40101 SET @OLD_SQL_MODE=@@SQL_MODE, SQL_MODE='NO_AUTO_VALUE_ON_ZERO' */;
/*!40111 SET @OLD_SQL_NOTES=@@SQL_NOTES, SQL_NOTES=0 */;

--
-- Current Database: `shangwei`
--

CREATE DATABASE /*!32312 IF NOT EXISTS*/ `shangwei` /*!40100 DEFAULT CHARACTER SET utf8 */;

USE `shangwei`;

--
-- Table structure for table `HisLoc`
--

DROP TABLE IF EXISTS `HisLoc`;
/*!40101 SET @saved_cs_client     = @@character_set_client */;
/*!40101 SET character_set_client = utf8 */;
CREATE TABLE `HisLoc` (
  `IT` smallint(5) DEFAULT NULL COMMENT 'ItemType',
  `IID` bigint(20) DEFAULT NULL COMMENT 'ItemID',
  `Lo` float DEFAULT '0' COMMENT '队员当前经度，东经为正数 ，西经为负数。例如东经: 16.33333，西经: -16.3333。精确度以设备能获取的精度为准。',
  `La` float DEFAULT '0' COMMENT '队员当前纬度，北纬为正数，南纬为负数。例如北纬: 23.0322，南纬: -43.22222。精确度以设备能获取的精度为准。',
  `He` float DEFAULT '0' COMMENT '队员当前海拔，海平面以上为正数，海平面以下为负数，以米为单位。例如海拔100.5米: 100.5。精确度以设备能获取的精度为准。',
  `Up` datetime NOT NULL COMMENT '时间',
  KEY `IT_IID` (`IT`,`IID`)
) ENGINE=InnoDB DEFAULT CHARSET=utf8 COMMENT='Item的历史位置（History Location）';
/*!40101 SET character_set_client = @saved_cs_client */;

--
-- Dumping data for table `HisLoc`
--

LOCK TABLES `HisLoc` WRITE;
/*!40000 ALTER TABLE `HisLoc` DISABLE KEYS */;
/*!40000 ALTER TABLE `HisLoc` ENABLE KEYS */;
UNLOCK TABLES;

--
-- Table structure for table `Human`
--

DROP TABLE IF EXISTS `Human`;
/*!40101 SET @saved_cs_client     = @@character_set_client */;
/*!40101 SET character_set_client = utf8 */;
CREATE TABLE `Human` (
  `ID` bigint(20) NOT NULL AUTO_INCREMENT COMMENT 'UserID',
  `Mail` char(30) NOT NULL COMMENT '注册邮箱',
  `Pwd` char(40) NOT NULL COMMENT 'Sha(Password+salt)',
  `Nick` char(20) NOT NULL COMMENT '昵称',
  `MailIdenti` char(36) DEFAULT NULL COMMENT '邮箱验证的字符 mysql function uuid()',
  `IdentiFail` datetime DEFAULT NULL COMMENT '邮箱验证字符失效时间',
  `RealMail` tinyint(4) NOT NULL DEFAULT '0' COMMENT '0 邮箱未验证 1邮箱验证通过后，需要将MailIdentify清空',
  `Allow` datetime NOT NULL DEFAULT '0000-00-00 00:00:00' COMMENT '用于登陆时检查，在这个时间之后才允许用户登陆',
  PRIMARY KEY (`ID`),
  UNIQUE KEY `RegMail` (`Mail`),
  UNIQUE KEY `NickName` (`Nick`)
) ENGINE=InnoDB AUTO_INCREMENT=4 DEFAULT CHARSET=utf8;
/*!40101 SET character_set_client = @saved_cs_client */;

--
-- Dumping data for table `Human`
--

LOCK TABLES `Human` WRITE;
/*!40000 ALTER TABLE `Human` DISABLE KEYS */;
INSERT INTO `Human` VALUES (1,'La','e405742cae3b8ad6bb96305dcb04398415284934','','fcc3c6c2-2c82-11e4-80ea-000c299da5ad','2014-08-27 14:09:42',0,'2014-08-25 14:09:42'),(2,'Leader@znr.io','6bd37efb8c6505448a03d4ee1f5d3c7f462ae824','Leader','063b401b-2c83-11e4-80ea-000c299da5ad','2014-08-27 14:09:58',0,'2014-08-25 14:09:58'),(3,'m1@znr.io','4def1404431dedddc60df1e27dae72b8ce10807c','m1','7056a3f7-2c83-11e4-80ea-000c299da5ad','2014-08-27 14:12:56',0,'2014-08-25 14:12:56');
/*!40000 ALTER TABLE `Human` ENABLE KEYS */;
UNLOCK TABLES;

--
-- Table structure for table `HumanStat`
--

DROP TABLE IF EXISTS `HumanStat`;
/*!40101 SET @saved_cs_client     = @@character_set_client */;
/*!40101 SET character_set_client = utf8 */;
CREATE TABLE `HumanStat` (
  `ID` bigint(20) NOT NULL COMMENT '用户ID',
  `Des` text COMMENT '自我描述',
  `Teams` text COMMENT '用户参加的Team的ID，用逗号分隔。每个TeamID最长是20个数字。最多可参与2500个Team',
  `Friends` text COMMENT '用户好友的UserID_base64编码的备注，用逗号分隔。每个UserID最长是20个数字,备注最多24个字符(base64后是32)。最多可以有1000个好友。',
  `PhotoFull` text,
  `PhotoMini` text,
  `MsgBd` tinyint(3) DEFAULT NULL COMMENT '留言板状态，/关闭/允许任何人留言/只允许同一个Team的人留言/只允许好友留言/允许同一个Team和好友留言',
  PRIMARY KEY (`ID`)
) ENGINE=InnoDB DEFAULT CHARSET=utf8;
/*!40101 SET character_set_client = @saved_cs_client */;

--
-- Dumping data for table `HumanStat`
--

LOCK TABLES `HumanStat` WRITE;
/*!40000 ALTER TABLE `HumanStat` DISABLE KEYS */;
INSERT INTO `HumanStat` VALUES (1,NULL,',',',',NULL,NULL,NULL),(2,NULL,',5,',',',NULL,NULL,NULL),(3,NULL,',3,2,',',',NULL,NULL,NULL);
/*!40000 ALTER TABLE `HumanStat` ENABLE KEYS */;
UNLOCK TABLES;

--
-- Table structure for table `Item`
--

DROP TABLE IF EXISTS `Item`;
/*!40101 SET @saved_cs_client     = @@character_set_client */;
/*!40101 SET character_set_client = utf8 */;
CREATE TABLE `Item` (
  `ID` smallint(5) NOT NULL AUTO_INCREMENT COMMENT 'ItemID',
  `Name` varchar(20) NOT NULL COMMENT '事物名称',
  `Des` text NOT NULL COMMENT '事物介绍的Uri',
  `PhotoMini` text NOT NULL COMMENT '事物的Mini头像的Uri',
  `PhotoFull` text NOT NULL COMMENT '事物的完整头像的Uri',
  `Video` text NOT NULL COMMENT '事物介绍视频的Uri',
  `Birth` datetime NOT NULL COMMENT '事物诞生日期',
  PRIMARY KEY (`ID`),
  UNIQUE KEY `ItemName` (`Name`)
) ENGINE=InnoDB AUTO_INCREMENT=7 DEFAULT CHARSET=utf8;
/*!40101 SET character_set_client = @saved_cs_client */;

--
-- Dumping data for table `Item`
--

LOCK TABLES `Item` WRITE;
/*!40000 ALTER TABLE `Item` DISABLE KEYS */;
INSERT INTO `Item` VALUES (1,'Human','人类交互终端','','','','2014-07-28 17:21:03'),(2,'PhyTrackTypeA','A型物理跟踪仪','','','','2014-07-28 17:18:05'),(3,'VirtualTrackerA1','A型虚拟迷你跟踪仪，Android平台','','','','2014-07-28 17:18:59'),(4,'VirtualTrackerA2','A型虚拟迷你跟踪仪，IOS平台','','','','2014-07-28 17:18:37'),(5,'Team','团队','','','','0000-00-00 00:00:00'),(6,'Public','公众','','','','0000-00-00 00:00:00');
/*!40000 ALTER TABLE `Item` ENABLE KEYS */;
UNLOCK TABLES;

--
-- Table structure for table `MsgBd`
--

DROP TABLE IF EXISTS `MsgBd`;
/*!40101 SET @saved_cs_client     = @@character_set_client */;
/*!40101 SET character_set_client = utf8 */;
CREATE TABLE `MsgBd` (
  `ID` bigint(20) NOT NULL AUTO_INCREMENT COMMENT 'MsgID，80亿个用户每秒产生一条消息，需要73年才会用完Unsigned BIGINT',
  `MT` tinyint(3) DEFAULT NULL COMMENT 'Type,Msg类型，0 系统消息  1短消息 2 长消息 3 文件 4 视频流 5语音流 6 加入Team请求',
  `Rd` tinyint(4) DEFAULT '0' COMMENT '0: 未读，1 已读',
  `ST` smallint(5) DEFAULT NULL COMMENT 'Sender Type发送者类型',
  `SID` bigint(20) DEFAULT NULL COMMENT 'Sender ID发送者ID',
  `RT` smallint(5) DEFAULT NULL COMMENT 'Receiver Type接收者类型，可以发送到：具体的终端（人/设备）、Team、Public(公众)',
  `RID` bigint(20) DEFAULT NULL COMMENT 'Receiver ID接受者ID',
  `Bo` tinytext COMMENT 'Body 消息体',
  `Birth` datetime DEFAULT NULL COMMENT '留言时间',
  `Exp` datetime DEFAULT NULL COMMENT '留言有效期，超过有效期的留言不再提供给接收者',
  PRIMARY KEY (`ID`),
  KEY `SType_SID` (`ST`,`SID`),
  KEY `RType_RID` (`RT`,`RID`)
) ENGINE=InnoDB AUTO_INCREMENT=31 DEFAULT CHARSET=utf8 ROW_FORMAT=COMPACT COMMENT='留言板消息(Message Board)，等待用户处理的消息。用户处理过后，Rd字段标记为1。';
/*!40101 SET character_set_client = @saved_cs_client */;

--
-- Dumping data for table `MsgBd`
--

LOCK TABLES `MsgBd` WRITE;
/*!40000 ALTER TABLE `MsgBd` DISABLE KEYS */;
INSERT INTO `MsgBd` VALUES (1,31,0,1,3,5,2,'2:d','2014-08-26 04:17:20',NULL),(2,31,1,1,3,1,2,'2:1d','2014-08-26 04:17:20',NULL),(3,35,1,1,9,1,3,'2','2014-08-26 04:42:58',NULL),(4,31,1,1,3,1,2,'2:2d','2014-08-26 04:17:20',NULL),(5,31,1,1,3,1,2,'2:3d','2014-08-26 04:17:20',NULL),(6,31,1,1,3,1,2,'2:d1','2014-08-26 04:17:20',NULL),(7,31,1,1,3,1,2,'2:d2','2014-08-26 04:17:20',NULL),(8,31,1,1,3,1,2,'2:d3','2014-08-26 04:17:20',NULL),(9,31,1,1,3,1,2,'2:d4','2014-08-26 04:17:20',NULL),(10,31,1,1,3,1,2,'2:d','2014-08-26 04:17:20',NULL),(11,31,1,1,3,1,2,'2:d','2014-08-26 04:17:20',NULL),(12,31,1,1,3,1,2,'2:d','2014-08-26 04:17:20',NULL),(13,31,1,1,3,1,2,'2:md','2014-08-26 04:17:20',NULL),(14,31,1,1,3,1,2,'2:md','2014-08-26 04:17:20',NULL),(15,31,1,1,3,1,2,'2:md','2014-08-26 04:17:20',NULL),(16,31,1,1,3,1,2,'2:d','2014-08-26 04:17:20',NULL),(17,31,1,1,3,1,2,'2:d','2014-08-26 04:17:20',NULL),(18,31,1,1,3,1,2,'2:d','2014-08-26 04:17:20',NULL),(19,31,1,1,3,1,2,'2:d','2014-08-26 04:17:20',NULL),(20,31,1,1,3,1,2,'2:d','2014-08-26 04:17:20',NULL),(21,31,1,1,3,1,2,'2:d','2014-08-26 04:17:20',NULL),(22,31,1,1,3,1,2,'2:d','2014-08-26 04:17:20',NULL),(23,31,1,1,3,1,2,'2:xxxd','2014-08-26 04:17:20',NULL),(24,35,1,1,3,1,2,'2:d','2014-08-26 04:17:20',NULL),(25,34,1,1,9,1,3,'1-2:2:','2014-08-26 14:14:33',NULL),(26,34,1,1,2,1,3,'1-2:2:','2014-08-26 14:32:27',NULL),(27,34,1,1,2,1,3,'1-2:2:','2014-08-26 14:32:29',NULL),(28,35,0,1,2,5,3,'','2014-08-26 19:05:12',NULL),(29,35,0,1,2,5,4,'','2014-08-26 19:09:32',NULL),(30,35,0,1,2,5,2,'','2014-08-26 20:12:46',NULL);
/*!40000 ALTER TABLE `MsgBd` ENABLE KEYS */;
UNLOCK TABLES;

--
-- Table structure for table `MsgRT`
--

DROP TABLE IF EXISTS `MsgRT`;
/*!40101 SET @saved_cs_client     = @@character_set_client */;
/*!40101 SET character_set_client = utf8 */;
CREATE TABLE `MsgRT` (
  `ID` bigint(20) NOT NULL AUTO_INCREMENT COMMENT 'MsgID，80亿个用户每秒产生一条消息，需要73年才会用完Unsigned BIGINT',
  `MT` tinyint(3) DEFAULT NULL COMMENT 'Type,Msg类型，0 系统消息  1短消息 2 长消息 3 文件 4 视频流 5语音流 6 加入Team请求',
  `ST` smallint(5) DEFAULT NULL COMMENT 'Sender Type发送者类型',
  `SID` bigint(20) DEFAULT NULL COMMENT 'Sender ID发送者ID',
  `RT` smallint(5) DEFAULT NULL COMMENT 'Receiver Type接收者类型，可以发送到：具体的终端（人/设备）、Team、Public(公众)',
  `RID` bigint(20) DEFAULT NULL COMMENT 'Receiver ID接受者ID',
  `Bo` tinytext COMMENT 'Body 消息体',
  `Exp` datetime DEFAULT '9999-12-31 23:59:59' COMMENT 'Expire 过期时间',
  PRIMARY KEY (`ID`),
  KEY `SType_SID` (`ST`,`SID`),
  KEY `RType_RID` (`RT`,`RID`)
) ENGINE=InnoDB AUTO_INCREMENT=52 DEFAULT CHARSET=utf8 COMMENT='实时消息表';
/*!40101 SET character_set_client = @saved_cs_client */;

--
-- Dumping data for table `MsgRT`
--

LOCK TABLES `MsgRT` WRITE;
/*!40000 ALTER TABLE `MsgRT` DISABLE KEYS */;
INSERT INTO `MsgRT` VALUES (1,24,1,3,5,1,'1,3','2014-08-25 14:18:21'),(2,22,1,3,5,1,'','2014-08-25 14:19:01'),(3,23,1,3,5,1,'','2014-08-25 14:19:34'),(4,22,1,3,5,1,'','2014-08-25 14:24:46'),(5,23,1,3,5,1,'','2014-08-25 14:24:51'),(6,22,1,3,5,1,'','2014-08-25 14:27:38'),(7,22,1,2,5,1,'','2014-08-25 14:27:56'),(8,23,1,2,5,1,'','2014-08-25 14:28:47'),(9,22,1,2,5,2,'','2014-08-25 14:30:13'),(10,22,1,3,5,1,'','2014-08-25 14:39:43'),(11,23,1,3,5,1,'','2014-08-25 14:39:51'),(12,22,1,3,5,1,'','2014-08-25 14:40:46'),(13,23,1,3,5,1,'','2014-08-25 14:40:56'),(14,24,1,3,5,3,'1,3','2014-08-25 14:43:06'),(15,22,1,2,5,1,'','2014-08-25 18:32:15'),(16,23,1,2,5,1,'','2014-08-25 18:32:19'),(17,22,1,2,5,4,'','2014-08-25 20:38:35'),(18,23,1,2,5,4,'','2014-08-25 20:38:39'),(19,22,1,2,5,4,'','2014-08-25 21:46:18'),(20,23,1,2,5,4,'','2014-08-25 21:46:23'),(21,22,1,2,5,1,'','2014-08-25 22:19:31'),(22,23,1,2,5,1,'','2014-08-25 22:19:42'),(23,22,1,2,5,1,'','2014-08-25 22:19:46'),(24,23,1,2,5,1,'','2014-08-25 22:19:48'),(25,22,1,2,5,4,'','2014-08-25 22:46:20'),(26,23,1,2,5,4,'','2014-08-25 22:46:23'),(27,22,1,2,5,3,'','2014-08-25 22:46:26'),(28,23,1,2,5,3,'','2014-08-25 22:46:29'),(29,24,1,2,5,3,'1,2','2014-08-26 01:33:22'),(30,24,1,2,5,3,'1,2','2014-08-26 01:44:09'),(31,33,1,3,5,2,'','2014-08-26 03:10:57'),(32,33,1,3,5,2,'','2014-08-26 03:11:59'),(33,24,1,2,5,3,'1,2','2014-08-26 03:12:50'),(34,33,1,3,5,2,'','2014-08-26 04:13:40'),(35,24,1,2,5,3,'1,2','2014-08-26 04:47:58'),(36,22,1,2,5,1,'','2014-08-26 14:51:15'),(37,23,1,2,5,1,'','2014-08-26 14:51:20'),(38,25,1,2,5,3,'','2014-08-26 19:10:12'),(39,25,1,2,5,4,'','2014-08-26 19:14:32'),(40,22,1,2,5,1,'','2014-08-26 20:14:48'),(41,23,1,2,5,1,'','2014-08-26 20:14:53'),(42,25,1,2,5,2,'','2014-08-26 20:17:46'),(43,33,1,3,5,1,'','2014-08-26 20:19:14'),(44,22,1,2,5,1,'','2014-08-26 22:13:19'),(45,23,1,2,5,1,'','2014-08-26 22:13:22'),(46,22,1,2,5,1,'','2014-08-27 00:16:44'),(47,23,1,2,5,1,'','2014-08-27 00:16:48'),(48,22,1,2,5,5,'','2014-08-27 07:31:29'),(49,23,1,2,5,5,'','2014-08-27 07:31:31'),(50,22,1,2,5,5,'','2014-08-27 07:31:57'),(51,23,1,2,5,5,'','2014-08-27 07:32:02');
/*!40000 ALTER TABLE `MsgRT` ENABLE KEYS */;
UNLOCK TABLES;

--
-- Table structure for table `Role`
--

DROP TABLE IF EXISTS `Role`;
/*!40101 SET @saved_cs_client     = @@character_set_client */;
/*!40101 SET character_set_client = utf8 */;
CREATE TABLE `Role` (
  `ID` smallint(5) NOT NULL AUTO_INCREMENT COMMENT '角色ID',
  `Na` char(20) DEFAULT NULL COMMENT '角色名称',
  `Des` tinytext COMMENT '角色描述',
  `Add` tinyint(3) DEFAULT '0' COMMENT '增加Team成员的权限',
  `Del` tinyint(3) DEFAULT '0' COMMENT '删除Team成员的权限',
  `Block` tinyint(3) DEFAULT '0' COMMENT '封锁Team成员的权限，被封锁的成员接收不到Team消息，发送到Team的所有消息被标记为封锁状态，用于警示其它的Team成员。',
  `Gag` tinyint(3) DEFAULT '0' COMMENT '禁言Team成员的权限，被禁言的成员正常接收Team消息，并且不能向Team发送消息',
  `Tag` tinyint(3) DEFAULT '0' COMMENT '为Team成员设置标签的权限',
  `Broad` tinyint(3) DEFAULT '0' COMMENT '向Team发送广播消息的权限',
  `SetRole` tinyint(3) DEFAULT '0' COMMENT '设置Team成员角色的权限',
  `Change` tinyint(3) DEFAULT '0' COMMENT '修改Team状态的权限',
  `Create` tinyint(3) DEFAULT '0' COMMENT '创建Team的权限',
  `TNum` tinyint(3) DEFAULT '0' COMMENT '最多创建的Team数量',
  PRIMARY KEY (`ID`)
) ENGINE=InnoDB AUTO_INCREMENT=10 DEFAULT CHARSET=utf8;
/*!40101 SET character_set_client = @saved_cs_client */;

--
-- Dumping data for table `Role`
--

LOCK TABLES `Role` WRITE;
/*!40000 ALTER TABLE `Role` DISABLE KEYS */;
INSERT INTO `Role` VALUES (1,'TeamOwer','Team的拥有者',1,1,1,1,1,1,1,1,0,0),(2,'TeamLeader','队长',1,1,1,1,1,1,0,0,0,0),(3,'TeamManage','Team的管理员',1,1,1,1,1,1,0,0,0,0),(4,'TeamMember','Team成员',0,0,0,0,0,1,0,0,0,0),(5,'SysAdmin','系统管理员',1,1,1,1,1,1,1,1,1,0),(6,'SilverUser','白银会员',0,0,0,0,0,0,0,0,0,0),(7,'GoldenUser','黄金会员',0,0,0,0,0,0,0,0,0,0),(8,'DiamondUser','钻石会员',0,0,0,0,0,0,0,0,0,0),(9,'NormalUser','普通会员',0,0,0,0,0,0,0,0,0,0);
/*!40000 ALTER TABLE `Role` ENABLE KEYS */;
UNLOCK TABLES;

--
-- Table structure for table `Session`
--

DROP TABLE IF EXISTS `Session`;
/*!40101 SET @saved_cs_client     = @@character_set_client */;
/*!40101 SET character_set_client = utf8 */;
CREATE TABLE `Session` (
  `ID` bigint(20) NOT NULL AUTO_INCREMENT COMMENT '会话ID，必须是唯一',
  `IID` bigint(20) NOT NULL COMMENT 'Client Item的ID，根据ItemType的值参照不同的ID',
  `IT` smallint(5) NOT NULL COMMENT 'Client Item Type的类型',
  `TID` bigint(20) NOT NULL DEFAULT '0' COMMENT 'TeamID, Current Team当前参与的Team, 0表示没有进入任何Team',
  `Ok` tinyint(3) NOT NULL DEFAULT '0' COMMENT '0 登陆验证失败  1 登陆验证成功 ',
  `Err` tinyint(3) NOT NULL DEFAULT '0' COMMENT '密码错误次数，达到五次后，将Item锁定30分钟',
  `CW` int(10) NOT NULL COMMENT 'Client Wait, Client正在等待的报文的编号， 随机生成',
  `CP` tinyint(3) NOT NULL COMMENT 'Client Policy, Client正在等待的报文编号的变化策略，这里存放的是递增值，可以取负数',
  `SW` int(10) NOT NULL COMMENT 'Server Wait, Server正在等待的报文的编号',
  `SP` tinyint(3) unsigned NOT NULL COMMENT 'Server Policy, Server等待的报文编号的变化策略，这里存放的是递增值，可以取负数',
  `Birth` datetime DEFAULT NULL COMMENT '会话创建时间',
  `Up` datetime NOT NULL COMMENT '客户端上一次活动时间',
  PRIMARY KEY (`ID`),
  UNIQUE KEY `ItemID_ItemType` (`IID`,`IT`)
) ENGINE=InnoDB AUTO_INCREMENT=44 DEFAULT CHARSET=utf8 COMMENT='Session放入类似Redis的共享内存中。';
/*!40101 SET character_set_client = @saved_cs_client */;

--
-- Dumping data for table `Session`
--

LOCK TABLES `Session` WRITE;
/*!40000 ALTER TABLE `Session` DISABLE KEYS */;
INSERT INTO `Session` VALUES (43,2,1,0,1,0,1014923696,40,634536052,84,'2014-08-27 07:26:17','2014-08-27 07:46:10');
/*!40000 ALTER TABLE `Session` ENABLE KEYS */;
UNLOCK TABLES;

--
-- Table structure for table `Team`
--

DROP TABLE IF EXISTS `Team`;
/*!40101 SET @saved_cs_client     = @@character_set_client */;
/*!40101 SET character_set_client = utf8 */;
CREATE TABLE `Team` (
  `ID` bigint(20) NOT NULL AUTO_INCREMENT COMMENT 'Team ID',
  `Na` varchar(50) NOT NULL DEFAULT '0' COMMENT 'Team Name',
  `CT` smallint(5) NOT NULL DEFAULT '0' COMMENT 'Team Owner''s Type',
  `CID` bigint(20) NOT NULL DEFAULT '0' COMMENT 'Team Owner''s ID',
  `De` tinytext NOT NULL COMMENT 'Describe',
  `Tg` char(10) NOT NULL DEFAULT '0' COMMENT 'Team标签',
  `Pub` tinyint(3) NOT NULL DEFAULT '0' COMMENT '0:信息完全公开  1:信息仅对Team成员公开 2:隐藏Team(不能被搜索到，信息只对Team成员公开)',
  `Alw` tinyint(3) NOT NULL DEFAULT '0' COMMENT '0: 随意加入 1:加入需审批  2:用“加入码”加入 3:禁止加入',
  `Pro` tinytext NOT NULL COMMENT '加入码提示',
  `Cod` tinytext NOT NULL COMMENT '加入码，如果Team的Join是3（用加入码加入），用户回答加入码后自动加入。',
  `Rol` smallint(5) NOT NULL COMMENT '新成员的默认角色',
  PRIMARY KEY (`ID`),
  KEY `TeamName` (`Na`)
) ENGINE=InnoDB AUTO_INCREMENT=6 DEFAULT CHARSET=utf8;
/*!40101 SET character_set_client = @saved_cs_client */;

--
-- Dumping data for table `Team`
--

LOCK TABLES `Team` WRITE;
/*!40000 ALTER TABLE `Team` DISABLE KEYS */;
INSERT INTO `Team` VALUES (5,'fgg',1,2,'dd','dd',1,0,'','',4);
/*!40000 ALTER TABLE `Team` ENABLE KEYS */;
UNLOCK TABLES;

--
-- Table structure for table `TeamStat`
--

DROP TABLE IF EXISTS `TeamStat`;
/*!40101 SET @saved_cs_client     = @@character_set_client */;
/*!40101 SET character_set_client = utf8 */;
CREATE TABLE `TeamStat` (
  `TID` bigint(20) NOT NULL COMMENT 'TeamID',
  `MT` smallint(5) NOT NULL COMMENT '队员类型，参照Item',
  `MID` bigint(20) DEFAULT NULL COMMENT '队员ID',
  `LT` smallint(5) NOT NULL COMMENT '队员直属Leader的类型（队员的直接上级 senior Leader的Type）',
  `LID` bigint(20) NOT NULL COMMENT '队员直属Leader的ID（队员的直接上级 senior Leader的ID）',
  `Rol` smallint(5) DEFAULT NULL COMMENT '队员在Team中的角色，不同角色拥有不同的管理权限',
  `Stat` tinyint(4) NOT NULL DEFAULT '0' COMMENT '队员的状态 0: 未进入 1 进入',
  `Attr` tinyint(4) NOT NULL DEFAULT '0' COMMENT '队员属性: 0:无 1:被封锁 2 被禁言',
  `De` tinytext NOT NULL COMMENT '备注信息',
  `Tg` tinytext NOT NULL COMMENT '队员的标签',
  `Lo` float DEFAULT '0' COMMENT '队员当前经度，东经为正数 ，西经为负数。例如东经: 16.33333，西经: -16.3333。精确度以设备能获取的精度为准。',
  `La` float DEFAULT '0' COMMENT '队员当前纬度，北纬为正数，南纬为负数。例如北纬: 23.0322，南纬: -43.22222。精确度以设备能获取的精度为准。',
  `He` float DEFAULT '0' COMMENT '队员当前海拔，海平面以上为正数，海平面以下为负数，以米为单位。例如海拔100.5米: 100.5。精确度以设备能获取的精度为准。',
  `Birth` datetime DEFAULT NULL COMMENT '队员加入时间',
  `Up` datetime NOT NULL COMMENT '队员最近一次位置更新时间',
  KEY `TID` (`TID`),
  KEY `MType_MID` (`MT`,`MID`),
  KEY `LType_LID` (`LT`,`LID`)
) ENGINE=InnoDB DEFAULT CHARSET=utf8;
/*!40101 SET character_set_client = @saved_cs_client */;

--
-- Dumping data for table `TeamStat`
--

LOCK TABLES `TeamStat` WRITE;
/*!40000 ALTER TABLE `TeamStat` DISABLE KEYS */;
INSERT INTO `TeamStat` VALUES (1,1,2,0,0,1,0,0,'','',0,0,0,'2014-08-25 14:10:56','2014-08-25 14:10:56'),(5,1,2,0,0,1,0,0,'','',0,0,0,'2014-08-27 01:26:29','2014-08-27 01:26:29');
/*!40000 ALTER TABLE `TeamStat` ENABLE KEYS */;
UNLOCK TABLES;
/*!40103 SET TIME_ZONE=@OLD_TIME_ZONE */;

/*!40101 SET SQL_MODE=@OLD_SQL_MODE */;
/*!40014 SET FOREIGN_KEY_CHECKS=@OLD_FOREIGN_KEY_CHECKS */;
/*!40014 SET UNIQUE_CHECKS=@OLD_UNIQUE_CHECKS */;
/*!40101 SET CHARACTER_SET_CLIENT=@OLD_CHARACTER_SET_CLIENT */;
/*!40101 SET CHARACTER_SET_RESULTS=@OLD_CHARACTER_SET_RESULTS */;
/*!40101 SET COLLATION_CONNECTION=@OLD_COLLATION_CONNECTION */;
/*!40111 SET SQL_NOTES=@OLD_SQL_NOTES */;

-- Dump completed on 2015-05-30 19:32:40
