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

 Date: 28/04/2023 13:50:38
*/

SET NAMES utf8mb4;
SET FOREIGN_KEY_CHECKS = 0;

-- ----------------------------
-- Table structure for deploy_task_details
-- ----------------------------
DROP TABLE IF EXISTS `deploy_task_details`;
CREATE TABLE `deploy_task_details` (
  `id` bigint(20) unsigned NOT NULL AUTO_INCREMENT COMMENT '自增id',
  `task_id` varchar(64) NOT NULL COMMENT '构建任务id',
  `push_date` varchar(128) NOT NULL DEFAULT '' COMMENT '任务提交时间',
  `project_name` varchar(128) NOT NULL COMMENT '项目名字',
  `branch_name` varchar(128) NOT NULL COMMENT '项目分支版本',
  `branch` varchar(255) NOT NULL,
  `git_ssh_url` varchar(128) NOT NULL COMMENT 'ssh格式url',
  `git_http_url` varchar(128) NOT NULL COMMENT 'http格式url',
  `push_user_name` varchar(32) NOT NULL COMMENT '提交的用户名',
  PRIMARY KEY (`id`),
  UNIQUE KEY `taskid` (`task_id`) COMMENT '任务id索引'
) ENGINE=InnoDB AUTO_INCREMENT=293 DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_general_ci;

SET FOREIGN_KEY_CHECKS = 1;
