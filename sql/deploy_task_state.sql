/*
 Navicat Premium Data Transfer

 Source Server         : cicd_mysql
 Source Server Type    : MySQL
 Source Server Version : 101102 (10.11.2-MariaDB)
 Source Host           : localhost:3306
 Source Schema         : wcicd

 Target Server Type    : MySQL
 Target Server Version : 101102 (10.11.2-MariaDB)
 File Encoding         : 65001

 Date: 28/04/2023 13:50:47
*/

SET NAMES utf8mb4;
SET FOREIGN_KEY_CHECKS = 0;

-- ----------------------------
-- Table structure for deploy_task_state
-- ----------------------------
DROP TABLE IF EXISTS `deploy_task_state`;
CREATE TABLE `deploy_task_state` (
  `id` bigint(20) unsigned NOT NULL AUTO_INCREMENT COMMENT '自增id',
  `task_id` varchar(64) NOT NULL COMMENT '构建任务id',
  `start_time` datetime DEFAULT '0000-00-00 00:00:00' COMMENT '构建任务开始时间',
  `node_id` varchar(255) DEFAULT NULL,
  `deploy_stats` int(11) NOT NULL COMMENT '构建实时状态',
  `deploy_details` varchar(255) NOT NULL COMMENT '构建状态详情',
  `push_image_name` varchar(255) DEFAULT NULL,
  `update_time` datetime NOT NULL DEFAULT '0000-00-00 00:00:00' ON UPDATE current_timestamp() COMMENT '记录更新时间',
  PRIMARY KEY (`id`),
  UNIQUE KEY `taskid` (`task_id`) COMMENT '任务id索引'
) ENGINE=InnoDB AUTO_INCREMENT=516 DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_general_ci;

SET FOREIGN_KEY_CHECKS = 1;
