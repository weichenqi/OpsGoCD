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

 Date: 28/04/2023 13:50:29
*/

SET NAMES utf8mb4;
SET FOREIGN_KEY_CHECKS = 0;

-- ----------------------------
-- Table structure for deploy_node_state
-- ----------------------------
DROP TABLE IF EXISTS `deploy_node_state`;
CREATE TABLE `deploy_node_state` (
  `id` bigint(20) NOT NULL AUTO_INCREMENT,
  `node_id` varchar(64) NOT NULL,
  `host_name` varchar(64) NOT NULL,
  `load_num` int(11) NOT NULL,
  `node_update_ts` int(11) NOT NULL,
  `last_update` datetime NOT NULL DEFAULT '0000-00-00 00:00:00' ON UPDATE current_timestamp(),
  PRIMARY KEY (`id`,`node_id`) USING BTREE,
  UNIQUE KEY `node_id_index` (`node_id`) USING BTREE
) ENGINE=InnoDB AUTO_INCREMENT=29065 DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_general_ci;

SET FOREIGN_KEY_CHECKS = 1;
