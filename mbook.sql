CREATE DATABASE  IF NOT EXISTS `mbook` /*!40100 DEFAULT CHARACTER SET utf8 */;
USE `mbook`;
-- MySQL dump 10.13  Distrib 5.6.24, for osx10.8 (x86_64)
--
-- Host: localhost    Database: mbook
-- ------------------------------------------------------
-- Server version	5.7.10

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
-- Table structure for table `md_attachment`
--

DROP TABLE IF EXISTS `md_attachment`;
/*!40101 SET @saved_cs_client     = @@character_set_client */;
/*!40101 SET character_set_client = utf8 */;
CREATE TABLE `md_attachment` (
  `attachment_id` int(11) NOT NULL AUTO_INCREMENT,
  `book_id` int(11) NOT NULL DEFAULT '0',
  `document_id` int(11) NOT NULL DEFAULT '0',
  `name` varchar(255) NOT NULL DEFAULT '',
  `path` varchar(2000) NOT NULL DEFAULT '',
  `size` double NOT NULL DEFAULT '0',
  `ext` varchar(50) NOT NULL DEFAULT '',
  `http_path` varchar(2000) NOT NULL DEFAULT '',
  `create_time` datetime NOT NULL,
  `create_at` int(11) NOT NULL DEFAULT '0',
  PRIMARY KEY (`attachment_id`)
) ENGINE=InnoDB DEFAULT CHARSET=utf8;
/*!40101 SET character_set_client = @saved_cs_client */;

--
-- Dumping data for table `md_attachment`
--

LOCK TABLES `md_attachment` WRITE;
/*!40000 ALTER TABLE `md_attachment` DISABLE KEYS */;
/*!40000 ALTER TABLE `md_attachment` ENABLE KEYS */;
UNLOCK TABLES;

--
-- Table structure for table `md_book_category`
--

DROP TABLE IF EXISTS `md_book_category`;
/*!40101 SET @saved_cs_client     = @@character_set_client */;
/*!40101 SET character_set_client = utf8 */;
CREATE TABLE `md_book_category` (
  `id` int(11) NOT NULL AUTO_INCREMENT,
  `book_id` int(11) NOT NULL DEFAULT '0',
  `category_id` int(11) NOT NULL DEFAULT '0',
  PRIMARY KEY (`id`),
  UNIQUE KEY `book_id` (`book_id`,`category_id`)
) ENGINE=InnoDB AUTO_INCREMENT=3 DEFAULT CHARSET=utf8;
/*!40101 SET character_set_client = @saved_cs_client */;

--
-- Dumping data for table `md_book_category`
--

LOCK TABLES `md_book_category` WRITE;
/*!40000 ALTER TABLE `md_book_category` DISABLE KEYS */;
INSERT INTO `md_book_category` VALUES (1,1,1),(2,1,4);
/*!40000 ALTER TABLE `md_book_category` ENABLE KEYS */;
UNLOCK TABLES;

--
-- Table structure for table `md_books`
--

DROP TABLE IF EXISTS `md_books`;
/*!40101 SET @saved_cs_client     = @@character_set_client */;
/*!40101 SET character_set_client = utf8 */;
CREATE TABLE `md_books` (
  `book_id` int(11) NOT NULL AUTO_INCREMENT,
  `book_name` varchar(500) NOT NULL DEFAULT '',
  `identify` varchar(100) NOT NULL DEFAULT '',
  `order_index` int(11) NOT NULL DEFAULT '0',
  `description` varchar(1000) NOT NULL DEFAULT '',
  `cover` varchar(1000) NOT NULL DEFAULT '',
  `editor` varchar(50) NOT NULL DEFAULT '',
  `status` int(11) NOT NULL DEFAULT '0',
  `privately_owned` int(11) NOT NULL DEFAULT '0',
  `private_token` varchar(500) DEFAULT NULL,
  `member_id` int(11) NOT NULL DEFAULT '0',
  `create_time` datetime NOT NULL,
  `modify_time` datetime NOT NULL,
  `release_time` datetime NOT NULL,
  `doc_count` int(11) NOT NULL DEFAULT '0',
  `comment_count` int(11) NOT NULL DEFAULT '0',
  `vcnt` int(11) NOT NULL DEFAULT '0',
  `star` int(11) NOT NULL DEFAULT '0',
  `score` int(11) NOT NULL DEFAULT '40',
  `cnt_score` int(11) NOT NULL DEFAULT '0',
  `cnt_comment` int(11) NOT NULL DEFAULT '0',
  `author` varchar(50) NOT NULL DEFAULT '',
  `author_url` varchar(1000) NOT NULL DEFAULT '',
  PRIMARY KEY (`book_id`),
  UNIQUE KEY `identify` (`identify`)
) ENGINE=InnoDB AUTO_INCREMENT=2 DEFAULT CHARSET=utf8;
/*!40101 SET character_set_client = @saved_cs_client */;

--
-- Dumping data for table `md_books`
--

LOCK TABLES `md_books` WRITE;
/*!40000 ALTER TABLE `md_books` DISABLE KEYS */;
INSERT INTO `md_books` VALUES (1,'演示','demo',0,'用于演示的书籍','/static/images/book.png','markdown',0,0,'',1,'2019-12-16 06:16:03','2019-12-16 06:16:03','2019-12-16 06:16:03',1,0,0,0,50,1,0,'','');
/*!40000 ALTER TABLE `md_books` ENABLE KEYS */;
UNLOCK TABLES;

--
-- Table structure for table `md_category`
--

DROP TABLE IF EXISTS `md_category`;
/*!40101 SET @saved_cs_client     = @@character_set_client */;
/*!40101 SET character_set_client = utf8 */;
CREATE TABLE `md_category` (
  `id` int(11) NOT NULL AUTO_INCREMENT,
  `pid` int(11) NOT NULL DEFAULT '0',
  `title` varchar(30) NOT NULL DEFAULT '',
  `intro` varchar(255) NOT NULL DEFAULT '',
  `icon` varchar(255) NOT NULL DEFAULT '',
  `cnt` int(11) NOT NULL DEFAULT '0',
  `sort` int(11) NOT NULL DEFAULT '0',
  `status` tinyint(1) NOT NULL DEFAULT '0',
  PRIMARY KEY (`id`),
  UNIQUE KEY `title` (`title`)
) ENGINE=InnoDB AUTO_INCREMENT=14 DEFAULT CHARSET=utf8;
/*!40101 SET character_set_client = @saved_cs_client */;

--
-- Dumping data for table `md_category`
--

LOCK TABLES `md_category` WRITE;
/*!40000 ALTER TABLE `md_category` DISABLE KEYS */;
INSERT INTO `md_category` VALUES (1,0,'演示','','',1,0,1),(2,0,'后端','','',0,0,1),(3,0,'前端','','',0,0,1),(4,1,'Demo','','',1,0,1),(5,2,'Go','','',0,0,1),(6,2,'JAVA','','',0,0,1),(7,2,'PHP','','',0,0,1),(8,2,'NET','','',0,0,1),(9,2,'Python','','',0,0,1),(10,3,'HTML','','',0,0,1),(11,3,'CSS','','',0,0,1),(12,3,'JavaScript','','',0,0,1),(13,3,'框架','','',0,0,1);
/*!40000 ALTER TABLE `md_category` ENABLE KEYS */;
UNLOCK TABLES;

--
-- Table structure for table `md_comments`
--

DROP TABLE IF EXISTS `md_comments`;
/*!40101 SET @saved_cs_client     = @@character_set_client */;
/*!40101 SET character_set_client = utf8 */;
CREATE TABLE `md_comments` (
  `id` int(11) NOT NULL AUTO_INCREMENT,
  `uid` int(11) NOT NULL DEFAULT '0',
  `book_id` int(11) NOT NULL DEFAULT '0',
  `content` varchar(255) NOT NULL DEFAULT '',
  `time_create` datetime NOT NULL,
  PRIMARY KEY (`id`),
  KEY `md_comments_uid` (`uid`),
  KEY `md_comments_book_id` (`book_id`)
) ENGINE=InnoDB DEFAULT CHARSET=utf8;
/*!40101 SET character_set_client = @saved_cs_client */;

--
-- Dumping data for table `md_comments`
--

LOCK TABLES `md_comments` WRITE;
/*!40000 ALTER TABLE `md_comments` DISABLE KEYS */;
/*!40000 ALTER TABLE `md_comments` ENABLE KEYS */;
UNLOCK TABLES;

--
-- Table structure for table `md_comments_0000`
--

DROP TABLE IF EXISTS `md_comments_0000`;
/*!40101 SET @saved_cs_client     = @@character_set_client */;
/*!40101 SET character_set_client = utf8 */;
CREATE TABLE `md_comments_0000` (
  `id` int(11) NOT NULL AUTO_INCREMENT,
  `uid` int(11) NOT NULL DEFAULT '0',
  `book_id` int(11) NOT NULL DEFAULT '0',
  `content` varchar(255) NOT NULL DEFAULT '',
  `time_create` datetime NOT NULL,
  PRIMARY KEY (`id`),
  KEY `md_comments_uid` (`uid`),
  KEY `md_comments_book_id` (`book_id`)
) ENGINE=InnoDB DEFAULT CHARSET=utf8;
/*!40101 SET character_set_client = @saved_cs_client */;

--
-- Dumping data for table `md_comments_0000`
--

LOCK TABLES `md_comments_0000` WRITE;
/*!40000 ALTER TABLE `md_comments_0000` DISABLE KEYS */;
/*!40000 ALTER TABLE `md_comments_0000` ENABLE KEYS */;
UNLOCK TABLES;

--
-- Table structure for table `md_comments_0001`
--

DROP TABLE IF EXISTS `md_comments_0001`;
/*!40101 SET @saved_cs_client     = @@character_set_client */;
/*!40101 SET character_set_client = utf8 */;
CREATE TABLE `md_comments_0001` (
  `id` int(11) NOT NULL AUTO_INCREMENT,
  `uid` int(11) NOT NULL DEFAULT '0',
  `book_id` int(11) NOT NULL DEFAULT '0',
  `content` varchar(255) NOT NULL DEFAULT '',
  `time_create` datetime NOT NULL,
  PRIMARY KEY (`id`),
  KEY `md_comments_uid` (`uid`),
  KEY `md_comments_book_id` (`book_id`)
) ENGINE=InnoDB DEFAULT CHARSET=utf8;
/*!40101 SET character_set_client = @saved_cs_client */;

--
-- Dumping data for table `md_comments_0001`
--

LOCK TABLES `md_comments_0001` WRITE;
/*!40000 ALTER TABLE `md_comments_0001` DISABLE KEYS */;
/*!40000 ALTER TABLE `md_comments_0001` ENABLE KEYS */;
UNLOCK TABLES;

--
-- Table structure for table `md_document_store`
--

DROP TABLE IF EXISTS `md_document_store`;
/*!40101 SET @saved_cs_client     = @@character_set_client */;
/*!40101 SET character_set_client = utf8 */;
CREATE TABLE `md_document_store` (
  `document_id` int(11) NOT NULL AUTO_INCREMENT,
  `markdown` longtext NOT NULL,
  `content` longtext NOT NULL,
  PRIMARY KEY (`document_id`)
) ENGINE=InnoDB AUTO_INCREMENT=2 DEFAULT CHARSET=utf8;
/*!40101 SET character_set_client = @saved_cs_client */;

--
-- Dumping data for table `md_document_store`
--

LOCK TABLES `md_document_store` WRITE;
/*!40000 ALTER TABLE `md_document_store` DISABLE KEYS */;
INSERT INTO `md_document_store` VALUES (1,'','');
/*!40000 ALTER TABLE `md_document_store` ENABLE KEYS */;
UNLOCK TABLES;

--
-- Table structure for table `md_documents`
--

DROP TABLE IF EXISTS `md_documents`;
/*!40101 SET @saved_cs_client     = @@character_set_client */;
/*!40101 SET character_set_client = utf8 */;
CREATE TABLE `md_documents` (
  `document_id` int(11) NOT NULL AUTO_INCREMENT,
  `document_name` varchar(500) NOT NULL DEFAULT '',
  `identify` varchar(100) DEFAULT 'null',
  `book_id` int(11) NOT NULL DEFAULT '0',
  `parent_id` int(11) NOT NULL DEFAULT '0',
  `order_sort` int(11) NOT NULL DEFAULT '0',
  `release` longtext,
  `create_time` datetime NOT NULL,
  `member_id` int(11) NOT NULL DEFAULT '0',
  `modify_time` datetime NOT NULL,
  `modify_at` int(11) NOT NULL DEFAULT '0',
  `version` bigint(20) NOT NULL DEFAULT '0',
  `vcnt` int(11) NOT NULL DEFAULT '0',
  PRIMARY KEY (`document_id`),
  UNIQUE KEY `book_id` (`book_id`,`identify`),
  KEY `md_documents_identify` (`identify`),
  KEY `md_documents_book_id_parent_id_order_sort` (`book_id`,`parent_id`,`order_sort`)
) ENGINE=InnoDB AUTO_INCREMENT=2 DEFAULT CHARSET=utf8;
/*!40101 SET character_set_client = @saved_cs_client */;

--
-- Dumping data for table `md_documents`
--

LOCK TABLES `md_documents` WRITE;
/*!40000 ALTER TABLE `md_documents` DISABLE KEYS */;
INSERT INTO `md_documents` VALUES (1,'空白文档','blank',1,0,0,'','2019-12-16 14:16:03',1,'2019-12-16 14:16:03',0,0,0);
/*!40000 ALTER TABLE `md_documents` ENABLE KEYS */;
UNLOCK TABLES;

--
-- Table structure for table `md_fans`
--

DROP TABLE IF EXISTS `md_fans`;
/*!40101 SET @saved_cs_client     = @@character_set_client */;
/*!40101 SET character_set_client = utf8 */;
CREATE TABLE `md_fans` (
  `id` int(11) NOT NULL AUTO_INCREMENT,
  `member_id` int(11) NOT NULL DEFAULT '0',
  `fans_id` int(11) NOT NULL DEFAULT '0',
  PRIMARY KEY (`id`),
  UNIQUE KEY `member_id` (`member_id`,`fans_id`),
  KEY `md_fans_fans_id` (`fans_id`)
) ENGINE=InnoDB DEFAULT CHARSET=utf8;
/*!40101 SET character_set_client = @saved_cs_client */;

--
-- Dumping data for table `md_fans`
--

LOCK TABLES `md_fans` WRITE;
/*!40000 ALTER TABLE `md_fans` DISABLE KEYS */;
/*!40000 ALTER TABLE `md_fans` ENABLE KEYS */;
UNLOCK TABLES;

--
-- Table structure for table `md_members`
--

DROP TABLE IF EXISTS `md_members`;
/*!40101 SET @saved_cs_client     = @@character_set_client */;
/*!40101 SET character_set_client = utf8 */;
CREATE TABLE `md_members` (
  `member_id` int(11) NOT NULL AUTO_INCREMENT,
  `account` varchar(30) NOT NULL DEFAULT '',
  `nickname` varchar(30) NOT NULL DEFAULT '',
  `password` varchar(255) NOT NULL DEFAULT '',
  `description` varchar(640) NOT NULL DEFAULT '',
  `email` varchar(100) NOT NULL DEFAULT '',
  `phone` varchar(20) DEFAULT 'null',
  `avatar` varchar(255) NOT NULL DEFAULT '',
  `role` int(11) NOT NULL DEFAULT '1',
  `status` int(11) NOT NULL DEFAULT '0',
  `create_time` datetime NOT NULL,
  `create_at` int(11) NOT NULL DEFAULT '0',
  `last_login_time` datetime DEFAULT NULL,
  PRIMARY KEY (`member_id`),
  UNIQUE KEY `account` (`account`),
  UNIQUE KEY `nickname` (`nickname`),
  UNIQUE KEY `email` (`email`)
) ENGINE=InnoDB AUTO_INCREMENT=3 DEFAULT CHARSET=utf8;
/*!40101 SET character_set_client = @saved_cs_client */;

--
-- Dumping data for table `md_members`
--

LOCK TABLES `md_members` WRITE;
/*!40000 ALTER TABLE `md_members` DISABLE KEYS */;
INSERT INTO `md_members` VALUES (1,'admin','admin','6fVynJQW4iV-KmCfHPrFucWFxwBKfGB-OY6Gu-9_QsHEFoEqCmgj-M-RwvM6WoIirokO|15|ced0f3c3ba8a223007bd5da110af9c0a3d3985e3c451e80c59789d91|7fec678fcc990d025b378232314a5339e96b26cb55b4ac2b13010f4a8d23c6af','','admin@ziyoubiancheng.com','','/static/images/avatar.png',0,0,'2019-12-16 06:13:31',0,'2019-12-16 14:13:31'),(2,'user1','user1','4mSZoWt1u91t3q6tcSZwFdIMT1wFR9o8Qzo53NRIhmd2FYqschKLYQknxcAADlHdfWLJ|15|98b702a40e8da1402a477983ab3b8fbbf5215b5dc4f5df526af28aa5|7ace8c5c5a49594446197ead810e34c4959e9f72ebcdd64218ecbff23500c5cd','','user1@ziyoubiancheng.com','','/static/images/avatar.png',2,0,'2019-12-19 17:04:26',0,'2019-12-20 01:04:26');
/*!40000 ALTER TABLE `md_members` ENABLE KEYS */;
UNLOCK TABLES;

--
-- Table structure for table `md_relationship`
--

DROP TABLE IF EXISTS `md_relationship`;
/*!40101 SET @saved_cs_client     = @@character_set_client */;
/*!40101 SET character_set_client = utf8 */;
CREATE TABLE `md_relationship` (
  `relationship_id` int(11) NOT NULL AUTO_INCREMENT,
  `member_id` int(11) NOT NULL DEFAULT '0',
  `book_id` int(11) NOT NULL DEFAULT '0',
  `role_id` int(11) NOT NULL DEFAULT '0',
  PRIMARY KEY (`relationship_id`),
  UNIQUE KEY `member_id` (`member_id`,`book_id`)
) ENGINE=InnoDB AUTO_INCREMENT=2 DEFAULT CHARSET=utf8;
/*!40101 SET character_set_client = @saved_cs_client */;

--
-- Dumping data for table `md_relationship`
--

LOCK TABLES `md_relationship` WRITE;
/*!40000 ALTER TABLE `md_relationship` DISABLE KEYS */;
INSERT INTO `md_relationship` VALUES (1,1,1,0);
/*!40000 ALTER TABLE `md_relationship` ENABLE KEYS */;
UNLOCK TABLES;

--
-- Table structure for table `md_score`
--

DROP TABLE IF EXISTS `md_score`;
/*!40101 SET @saved_cs_client     = @@character_set_client */;
/*!40101 SET character_set_client = utf8 */;
CREATE TABLE `md_score` (
  `id` int(11) NOT NULL AUTO_INCREMENT,
  `book_id` int(11) NOT NULL DEFAULT '0',
  `uid` int(11) NOT NULL DEFAULT '0',
  `score` int(11) NOT NULL DEFAULT '0',
  `time_create` datetime NOT NULL,
  PRIMARY KEY (`id`),
  UNIQUE KEY `uid` (`uid`,`book_id`)
) ENGINE=InnoDB DEFAULT CHARSET=utf8;
/*!40101 SET character_set_client = @saved_cs_client */;

--
-- Dumping data for table `md_score`
--

LOCK TABLES `md_score` WRITE;
/*!40000 ALTER TABLE `md_score` DISABLE KEYS */;
/*!40000 ALTER TABLE `md_score` ENABLE KEYS */;
UNLOCK TABLES;

--
-- Table structure for table `md_star`
--

DROP TABLE IF EXISTS `md_star`;
/*!40101 SET @saved_cs_client     = @@character_set_client */;
/*!40101 SET character_set_client = utf8 */;
CREATE TABLE `md_star` (
  `id` int(11) NOT NULL AUTO_INCREMENT,
  `member_id` int(11) NOT NULL DEFAULT '0',
  `book_id` int(11) NOT NULL DEFAULT '0',
  PRIMARY KEY (`id`),
  UNIQUE KEY `member_id` (`member_id`,`book_id`),
  KEY `md_star_member_id` (`member_id`)
) ENGINE=InnoDB DEFAULT CHARSET=utf8;
/*!40101 SET character_set_client = @saved_cs_client */;

--
-- Dumping data for table `md_star`
--

LOCK TABLES `md_star` WRITE;
/*!40000 ALTER TABLE `md_star` DISABLE KEYS */;
/*!40000 ALTER TABLE `md_star` ENABLE KEYS */;
UNLOCK TABLES;
/*!40103 SET TIME_ZONE=@OLD_TIME_ZONE */;

/*!40101 SET SQL_MODE=@OLD_SQL_MODE */;
/*!40014 SET FOREIGN_KEY_CHECKS=@OLD_FOREIGN_KEY_CHECKS */;
/*!40014 SET UNIQUE_CHECKS=@OLD_UNIQUE_CHECKS */;
/*!40101 SET CHARACTER_SET_CLIENT=@OLD_CHARACTER_SET_CLIENT */;
/*!40101 SET CHARACTER_SET_RESULTS=@OLD_CHARACTER_SET_RESULTS */;
/*!40101 SET COLLATION_CONNECTION=@OLD_COLLATION_CONNECTION */;
/*!40111 SET SQL_NOTES=@OLD_SQL_NOTES */;

-- Dump completed on 2019-12-20  9:11:18
