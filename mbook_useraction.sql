CREATE DATABASE  IF NOT EXISTS `mbook_useraction` /*!40100 DEFAULT CHARACTER SET utf8 */;
USE `mbook_useraction`;
-- MySQL dump 10.13  Distrib 5.6.24, for osx10.8 (x86_64)
--
-- Host: localhost    Database: mbook_useraction
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
) ENGINE=InnoDB AUTO_INCREMENT=2 DEFAULT CHARSET=utf8;
/*!40101 SET character_set_client = @saved_cs_client */;

--
-- Dumping data for table `md_comments_0001`
--

LOCK TABLES `md_comments_0001` WRITE;
/*!40000 ALTER TABLE `md_comments_0001` DISABLE KEYS */;
INSERT INTO `md_comments_0001` VALUES (1,1,1,'hello','2019-12-20 09:02:52');
/*!40000 ALTER TABLE `md_comments_0001` ENABLE KEYS */;
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
) ENGINE=InnoDB AUTO_INCREMENT=2 DEFAULT CHARSET=utf8;
/*!40101 SET character_set_client = @saved_cs_client */;

--
-- Dumping data for table `md_fans`
--

LOCK TABLES `md_fans` WRITE;
/*!40000 ALTER TABLE `md_fans` DISABLE KEYS */;
INSERT INTO `md_fans` VALUES (1,1,2);
/*!40000 ALTER TABLE `md_fans` ENABLE KEYS */;
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
) ENGINE=InnoDB AUTO_INCREMENT=3 DEFAULT CHARSET=utf8;
/*!40101 SET character_set_client = @saved_cs_client */;

--
-- Dumping data for table `md_score`
--

LOCK TABLES `md_score` WRITE;
/*!40000 ALTER TABLE `md_score` DISABLE KEYS */;
INSERT INTO `md_score` VALUES (2,1,1,50,'2019-12-20 01:03:34');
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

-- Dump completed on 2019-12-20  9:13:39
