CREATE DATABASE im_db; -- 创建IM数据库

USE im_db; -- 使用IM数据库

-- 用户表：存储用户基本信息
CREATE TABLE users (
                       id INT AUTO_INCREMENT PRIMARY KEY COMMENT '用户ID，自增主键',
                       phone_number VARCHAR(20) NOT NULL UNIQUE COMMENT '用户电话号码，唯一，不能为空',
                       email VARCHAR(255) NOT NULL UNIQUE COMMENT '用户邮箱，唯一，不能为空',
                       username VARCHAR(255) NOT NULL COMMENT '用户名，不能为空',
                       password_hash VARCHAR(255) NOT NULL COMMENT '密码哈希值，不能为空',
                       avatar_url VARCHAR(255) COMMENT '用户头像URL，允许为空',
                       bio TEXT COMMENT '用户个人简介，允许为空',
                       gender ENUM('male', 'female', 'other') NOT NULL COMMENT '用户性别，枚举类型，不能为空',
                       address VARCHAR(255) COMMENT '用户住址，允许为空',
                       city VARCHAR(100) COMMENT '用户所在城市，允许为空',
                       state VARCHAR(100) COMMENT '用户所在州/省，允许为空',
                       country VARCHAR(100) COMMENT '用户所在国家，允许为空',
                       postal_code VARCHAR(20) COMMENT '用户邮政编码，允许为空',
                       date_of_birth VARCHAR(256) COMMENT '用户出生日期，允许为空',
                       created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP COMMENT '记录创建时间，默认为当前时间',
                       updated_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP ON UPDATE CURRENT_TIMESTAMP COMMENT '记录更新时间，在更新时自动设置为当前时间'
) COMMENT='用户信息表';

-- 分组表：存储用户的好友分组信息
CREATE TABLE friend_groups (
                               id INT AUTO_INCREMENT PRIMARY KEY COMMENT '分组ID，自增主键',
                               user_id INT NOT NULL COMMENT '用户ID，不能为空',
                               group_name VARCHAR(255) NOT NULL COMMENT '分组名称，不能为空',
                               created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP COMMENT '记录创建时间，默认为当前时间',
                               updated_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP ON UPDATE CURRENT_TIMESTAMP COMMENT '记录更新时间，在更新时自动设置为当前时间',
                               FOREIGN KEY (user_id) REFERENCES users(id)
) COMMENT='好友分组表';

-- 好友关系表：存储用户之间的好友关系
CREATE TABLE friendships (
                             id INT AUTO_INCREMENT PRIMARY KEY COMMENT '好友关系ID，自增主键',
                             user_id INT NOT NULL COMMENT '用户ID，不能为空',
                             friend_id INT NOT NULL COMMENT '好友的用户ID，不能为空',
                             status ENUM('pending', 'accepted', 'blocked') DEFAULT 'pending' COMMENT '好友关系状态，默认为pending',
                             group_id INT COMMENT '分组ID，允许为空',
                             remark VARCHAR(255) COMMENT '好友备注，允许为空,默认的话采用用户名',
                             created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP COMMENT '记录创建时间，默认为当前时间',
                             updated_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP ON UPDATE CURRENT_TIMESTAMP COMMENT '记录更新时间，在更新时自动设置为当前时间',
                             FOREIGN KEY (user_id) REFERENCES users(id) ,
                             FOREIGN KEY (friend_id) REFERENCES users(id) ,
                             FOREIGN KEY (group_id) REFERENCES friend_groups(id)
) COMMENT='好友关系表';

-- 群组表：存储群组信息
CREATE TABLE `groups` (
                          id INT AUTO_INCREMENT PRIMARY KEY COMMENT '群组ID，自增主键',
                          name VARCHAR(255) NOT NULL COMMENT '群组名称，不能为空',
                          owner_id INT NOT NULL COMMENT '群主的用户ID，不能为空',
                          group_avatar VARCHAR(255) DEFAULT 'https://encrypted-tbn0.gstatic.com/images?q=tbn:ANd9GcQ3w2fqb71MsCj97IKLAUXoI6BS4IfeCeEoq_XGS3X2CErGlYyP4xxX4eQ&s' COMMENT '群组头像URL，不允许为空',
                          created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP COMMENT '记录创建时间，默认为当前时间',
                          updated_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP ON UPDATE CURRENT_TIMESTAMP COMMENT '记录更新时间，在更新时自动设置为当前时间',
                          FOREIGN KEY (owner_id) REFERENCES users(id)
) COMMENT='群组信息表';

-- 消息表：存储用户发送的消息
CREATE TABLE messages (
                          id INT AUTO_INCREMENT PRIMARY KEY COMMENT '消息ID，自增主键',
                          sender_id INT NOT NULL COMMENT '发送者ID，不能为空',
                          receiver_user_id INT COMMENT '接收者的用户ID，用于私聊消息',
                          receiver_group_id INT COMMENT '接收者的群组ID，用于群聊消息',
                          content TEXT COMMENT '消息内容，允许为空',
                          message_type ENUM('text', 'image', 'video', 'file') DEFAULT 'text' COMMENT '消息类型，默认为text',
                          created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP COMMENT '记录创建时间，默认为当前时间',
                          updated_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP ON UPDATE CURRENT_TIMESTAMP COMMENT '记录更新时间，在更新时自动设置为当前时间',
                          FOREIGN KEY (sender_id) REFERENCES users(id) ,
                          FOREIGN KEY (receiver_user_id) REFERENCES users(id) ,
                          FOREIGN KEY (receiver_group_id) REFERENCES `groups`(id) ,
                          CONSTRAINT chk_receiver CHECK (
                              (receiver_user_id IS NOT NULL AND receiver_group_id IS NULL) OR
                              (receiver_user_id IS NULL AND receiver_group_id IS NOT NULL)
                              )
) COMMENT='消息表';

-- 消息已读状态表：存储消息的已读状态
CREATE TABLE message_read_status (
                                     id INT AUTO_INCREMENT PRIMARY KEY COMMENT '已读状态ID，自增主键',
                                     message_id INT NOT NULL COMMENT '消息ID，不能为空，引用messages表中的id',
                                     user_id INT NOT NULL COMMENT '用户ID，不能为空，引用users表中的id',
                                     read_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP ,
                                     FOREIGN KEY (message_id) REFERENCES messages(id) ,
                                     FOREIGN KEY (user_id) REFERENCES users(id)
) COMMENT='消息已读状态表';

-- 群组成员表：存储群组成员信息
CREATE TABLE group_members (
                               id INT AUTO_INCREMENT PRIMARY KEY COMMENT '群组成员关系ID，自增主键',
                               group_id INT NOT NULL COMMENT '群组ID，不能为空',
                               user_id INT NOT NULL COMMENT '用户ID，不能为空',
                               role ENUM('admin', 'member','owner') NOT NULL DEFAULT 'member' COMMENT '成员类型，默认为普通成员',
                               title VARCHAR(255) COMMENT '成员在群组中的称号，允许为空',
                               level INT NOT NULL DEFAULT 1 COMMENT '成员等级，默认为1',
                               nickname VARCHAR(255) COMMENT '用户在群组中的昵称，允许为空',
                               created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP COMMENT '记录创建时间，默认为当前时间',
                               updated_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP ON UPDATE CURRENT_TIMESTAMP COMMENT '记录更新时间，在更新时自动设置为当前时间',
                               FOREIGN KEY (group_id) REFERENCES `groups`(id) ,
                               FOREIGN KEY (user_id) REFERENCES users(id)
) COMMENT='群组成员表';

-- 示例通知表：存储系统通知信息
CREATE TABLE notifications (
                               id INT AUTO_INCREMENT PRIMARY KEY COMMENT '通知ID，自增主键',
                               sender_id INT NOT NULL COMMENT '发送者ID，不能为空',
                               receiver_id INT NOT NULL COMMENT '接收者ID，不能为空',
                               type ENUM('message', 'friend_request', 'group_request', 'group_invite', 'other') DEFAULT 'message' COMMENT '通知类型，默认为message',
                               content TEXT COMMENT '通知内容，允许为空',
                               is_read BOOLEAN DEFAULT FALSE COMMENT '通知是否已读，默认为未读',
                               status ENUM('pending', 'accepted', 'rejected') DEFAULT 'pending' COMMENT '通知状态，默认为pending',
                               group_id INT COMMENT '群组ID，允许为空',
                               created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP COMMENT '记录创建时间，默认为当前时间',
                               updated_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP ON UPDATE CURRENT_TIMESTAMP COMMENT '记录更新时间，在更新时自动设置为当前时间',
                               FOREIGN KEY (sender_id) REFERENCES users(id) ,
                               FOREIGN KEY (receiver_id) REFERENCES users(id) ,
                               FOREIGN KEY (group_id) REFERENCES `groups`(id)
) COMMENT='通知表';

-- 系统日志表：存储系统操作日志
CREATE TABLE system_logs (
                             id INT AUTO_INCREMENT PRIMARY KEY COMMENT '系统日志ID，自增主键',
                             log_type VARCHAR(255) COMMENT '日志类型',
                             message TEXT COMMENT '日志内容',
                             created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP COMMENT '记录创建时间，默认为当前时间'
) COMMENT='系统日志表';

# 新增的字段
ALTER TABLE `groups`
    ADD COLUMN `announcement` TEXT COMMENT '群公告，允许为空' AFTER `group_avatar`,
    ADD COLUMN `description` TEXT COMMENT '群描述，允许为空' AFTER `announcement`;