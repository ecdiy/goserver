/*
Navicat MySQL Data Transfer

Source Server         : localhost
Source Server Version : 50721
Source Host           : localhost:3306
Source Database       : goserver

Target Server Type    : MYSQL
Target Server Version : 50721
File Encoding         : 65001

Date: 2018-11-14 13:05:26
*/

SET FOREIGN_KEY_CHECKS=0;

-- ----------------------------
-- Table structure for GkUser
-- ----------------------------
DROP TABLE IF EXISTS `GkUser`;
CREATE TABLE `GkUser` (
  `Id` bigint(20) NOT NULL AUTO_INCREMENT,
  `Username` varchar(64) NOT NULL COMMENT '用户名',
  `Password` varchar(64) NOT NULL,
  `PasswordError` int(11) DEFAULT '0',
  PRIMARY KEY (`Id`),
  UNIQUE KEY `Username` (`Username`) USING BTREE
) ENGINE=MyISAM AUTO_INCREMENT=11049 DEFAULT CHARSET=utf8mb4 ROW_FORMAT=DYNAMIC;

-- ----------------------------
-- Records of GkUser
-- ----------------------------
INSERT INTO `GkUser` VALUES ('11047', 'test', 'test', '0');
INSERT INTO `GkUser` VALUES ('11048', 'rewrew', '89e12ff9a15a0027561bdf989e8e1388', '0');

-- ----------------------------
-- Table structure for Project
-- ----------------------------
DROP TABLE IF EXISTS `Project`;
CREATE TABLE `Project` (
  `Id` bigint(20) NOT NULL AUTO_INCREMENT,
  `HomeUrl` varchar(255) DEFAULT NULL,
  `Name` varchar(160) DEFAULT NULL,
  `CreateAt` datetime DEFAULT NULL,
  `CatId` bigint(20) DEFAULT NULL,
  `Star` varchar(16) DEFAULT NULL,
  `Site` varchar(64) DEFAULT NULL,
  `Summary` mediumtext,
  `ItemCount` int(11) NOT NULL DEFAULT '0' COMMENT '文档总数',
  `UserId` bigint(20) DEFAULT '0',
  PRIMARY KEY (`Id`),
  UNIQUE KEY `HomeUrl` (`HomeUrl`)
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4;

-- ----------------------------
-- Records of Project
-- ----------------------------

-- ----------------------------
-- Table structure for Token
-- ----------------------------
DROP TABLE IF EXISTS `Token`;
CREATE TABLE `Token` (
  `UserId` bigint(20) NOT NULL,
  `Ua` varchar(8) NOT NULL,
  `Token` varchar(64) NOT NULL,
  `CreateAt` datetime NOT NULL,
  UNIQUE KEY `Token` (`Token`) USING BTREE,
  UNIQUE KEY `UserId` (`UserId`,`Ua`) USING BTREE
) ENGINE=MyISAM DEFAULT CHARSET=latin1;

-- ----------------------------
-- Records of Token
-- ----------------------------
INSERT INTO `Token` VALUES ('11047', 'web', '57b4d577798e05a16726209ec1917526', '2018-11-14 13:04:57');
INSERT INTO `Token` VALUES ('11048', 'web', '1a760e15136a02daea0fbb619b6deaad', '2018-11-14 10:50:44');

-- ----------------------------
-- Procedure structure for GetMailAjax
-- ----------------------------
DROP PROCEDURE IF EXISTS `GetMailAjax`;
DELIMITER ;;
CREATE DEFINER=`root`@`localhost` PROCEDURE `GetMailAjax`(ginUserId bigint)
    COMMENT 'user:map'
BEGIN

	select Email,Username from GkUser where Id=ginUserId;

END
;;
DELIMITER ;

-- ----------------------------
-- Procedure structure for HelloWorldAjax
-- ----------------------------
DROP PROCEDURE IF EXISTS `HelloWorldAjax`;
DELIMITER ;;
CREATE DEFINER=`root`@`localhost` PROCEDURE `HelloWorldAjax`()
    COMMENT 'test:string'
BEGIN

	select 'Hello' Result;

END
;;
DELIMITER ;

-- ----------------------------
-- Procedure structure for ParamAjax
-- ----------------------------
DROP PROCEDURE IF EXISTS `ParamAjax`;
DELIMITER ;;
CREATE DEFINER=`root`@`localhost` PROCEDURE `ParamAjax`(inParam text)
    COMMENT 'param:map,list:list'
BEGIN

	select inParam Param;

	select 1 A,2 B
	union select 2 A ,3 B;

END
;;
DELIMITER ;

-- ----------------------------
-- Procedure structure for TestSys
-- ----------------------------
DROP PROCEDURE IF EXISTS `TestSys`;
DELIMITER ;;
CREATE DEFINER=`root`@`localhost` PROCEDURE `TestSys`(insp text,ginUserId bigint)
    COMMENT 'result:map'
BEGIN
	#Routine body goes here...

	select insp SpName;

END
;;
DELIMITER ;

-- ----------------------------
-- Procedure structure for UserListAjax
-- ----------------------------
DROP PROCEDURE IF EXISTS `UserListAjax`;
DELIMITER ;;
CREATE DEFINER=`root`@`localhost` PROCEDURE `UserListAjax`()
    COMMENT 'list:list,total:int'
BEGIN

	SELECT * from GkUser;
	select count(*) from GkUser;

END
;;
DELIMITER ;

-- ----------------------------
-- Procedure structure for UserLoginAjax
-- ----------------------------
DROP PROCEDURE IF EXISTS `UserLoginAjax`;
DELIMITER ;;
CREATE DEFINER=`root`@`localhost` PROCEDURE `UserLoginAjax`(inUsername varchar(64),inPassword varchar(64),ua varchar(8))
    COMMENT 'status:o'
proc:
BEGIN
	DECLARE pErr,pId bigint DEFAULT 0;
	DECLARE pPass,pToken VARCHAR(64);

	select PasswordError,Id,`Password` into pErr,pId,pPass from GkUser where Username=inUsername;

	if pId=0 THEN
		select 1 code,'用户名不存在' msg;
		LEAVE proc ;
	end if;

	if pErr>3 THEN
		select 2 code,'最大错误次数' msg;
		LEAVE proc ;
	end if;
	
	if inPassword=pPass THEN
		update `GkUser` set PasswordError=0 where Id=pId;
		set pToken=md5(CONCAT(pId,inUsername,now()));
		
		DELETE from Token where UserId=pId and Ua=ua;
		INSERT INTO `Token` (`UserId`, `Ua`,`Token`,`CreateAt`) VALUES (pId,ua,pToken,now());
		 
		select 0 code,pToken Token,inUsername Username, Id 
		from GkUser where Id=pId;

	ELSE
		update GkUser Set PasswordError=PasswordError+1 where Id=pId and PasswordError<4;
		select 3 code,'密码错误' msg;
	end if;

END
;;
DELIMITER ;

-- ----------------------------
-- Procedure structure for UserLoginCaptcha
-- ----------------------------
DROP PROCEDURE IF EXISTS `UserLoginCaptcha`;
DELIMITER ;;
CREATE DEFINER=`root`@`localhost` PROCEDURE `UserLoginCaptcha`(inUsername varchar(64),inPassword varchar(64),ua varchar(8))
    COMMENT 'status:o'
proc:
BEGIN
	DECLARE pErr,pId bigint DEFAULT 0;
	DECLARE pPass,pToken VARCHAR(64);

	select PasswordError,Id,`Password` into pErr,pId,pPass from GkUser where Username=inUsername;

	if pId=0 THEN
		select 1 code,'用户名不存在' msg;
		LEAVE proc ;
	end if;
	
-- 	if md5(CONCAT(inUsername,',',inPassword))=pPass THEN

	if  inPassword=pPass THEN
		update `GkUser` set PasswordError=0 where Id=pId;
		set pToken=md5(CONCAT(pId,inUsername,now()));
		DELETE from Token where UserId=pId and Ua=ua;
		INSERT INTO `Token` (`UserId`, `Ua`,`Token`,`CreateAt`) VALUES (pId,ua,pToken,now());
		select 0 code,'' msg,pId UserId ,pToken Token;
	ELSE
		update GkUser Set PasswordError=PasswordError+1 where Id=pId and PasswordError<4;
		select 3 code,'密码错误' msg;
	end if;

END
;;
DELIMITER ;

-- ----------------------------
-- Procedure structure for UserRegister
-- ----------------------------
DROP PROCEDURE IF EXISTS `UserRegister`;
DELIMITER ;;
CREATE DEFINER=`root`@`localhost` PROCEDURE `UserRegister`(inPassword varchar(32),inUsername varchar(64),ua varchar(8))
    COMMENT 'result:map'
proc:
BEGIN
	DECLARE pExt int DEFAULT 0;
	DECLARE pId bigint;
	DECLARE pToken VARCHAR(64);
	
	if LENGTH(inUsername)<5 or LENGTH(inPassword)<6   then
		select 1 Code,'参数错误' msg;
		LEAVE proc;
	end if;

	select count(*) into pExt from `GkUser` where Username=inUsername;
	if pExt>0 THEN
		select 1000 Code,'用户名已存在' msg;
		LEAVE proc;
	end if;

 
	insert into GkUser(Username,Password,PasswordError) 			VALUES (inUsername, inPassword,0);

	set pId=@@IDENTITY;
	set pToken=md5(CONCAT(pId,inUsername,now()));
	update GkUser set `Password`=md5(CONCAT(pId,',',inPassword)) where Username=inUsername;
 
	INSERT INTO `Token` (`UserId`, `Ua`,`Token`,`CreateAt`) VALUES (pId,ua,pToken,now());


	select 0 Code,pToken Token,pId UserId;

END
;;
DELIMITER ;

-- ----------------------------
-- Procedure structure for UserRule
-- ----------------------------
DROP PROCEDURE IF EXISTS `UserRule`;
DELIMITER ;;
CREATE DEFINER=`root`@`localhost` PROCEDURE `UserRule`(insp text,ginUserId bigint)
BEGIN
	#Routine body goes here...
	
	select 1;
END
;;
DELIMITER ;
SET FOREIGN_KEY_CHECKS=1;
