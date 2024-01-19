DROP TABLE IF EXISTS `borrow_fee`;

CREATE TABLE `borrow_fee` (
  `member_fk` int(11) NOT NULL COMMENT '會員pk',
  `pk` int(10) unsigned NOT NULL AUTO_INCREMENT COMMENT '主鍵PK',
  `type` int(11) NOT NULL DEFAULT 1 COMMENT '業務類型 1.新合約 2.續期',
  `borrow_fee` decimal(12,2) NOT NULL COMMENT '管理費',
  `create_time` datetime DEFAULT NULL COMMENT '發生時間',
  PRIMARY KEY (`pk`)
);

DROP TABLE IF EXISTS `member`;

CREATE TABLE `member` (
  `pk` int(10) unsigned NOT NULL AUTO_INCREMENT COMMENT '用戶pk',
  `username` varchar(16) NOT NULL COMMENT '登入帳號',
  `create_time` datetime NOT NULL COMMENT '註冊時間',
  PRIMARY KEY (`pk`)
);
