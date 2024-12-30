CREATE DATABASE im_db;

USE im_db;
-- 用户表
CREATE TABLE users (
                       id INT AUTO_INCREMENT PRIMARY KEY, -- 用户ID，自增主键
                       phone_number VARCHAR(20) NOT NULL UNIQUE, -- 用户电话号码，唯一，不能为空
                       email VARCHAR(255) NOT NULL UNIQUE, -- 用户邮箱，唯一，不能为空
                       username VARCHAR(255) NOT NULL, -- 用户名，不能为空
                       password_hash VARCHAR(255) NOT NULL, -- 密码哈希值，不能为空
                       avatar_url VARCHAR(255), -- 用户头像URL，允许为空
                       bio TEXT, -- 用户个人简介，允许为空
                       gender ENUM('male', 'female', 'other') NOT NULL, -- 用户性别，枚举类型，不能为空
                       address VARCHAR(255), -- 用户住址，允许为空
                       city VARCHAR(100), -- 用户所在城市，允许为空
                       state VARCHAR(100), -- 用户所在州/省，允许为空
                       country VARCHAR(100), -- 用户所在国家，允许为空
                       postal_code VARCHAR(20), -- 用户邮政编码，允许为空
                       date_of_birth VARCHAR(256), -- 用户出生日期，允许为空
                       created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP, -- 记录创建时间，默认为当前时间
                       updated_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP ON UPDATE CURRENT_TIMESTAMP -- 记录更新时间，在更新时自动设置为当前时间
);

-- 分组表
CREATE TABLE friend_groups (
                               id INT AUTO_INCREMENT PRIMARY KEY, -- 分组ID，自增主键
                               user_id INT NOT NULL, -- 用户ID，不能为空
                               group_name VARCHAR(255) NOT NULL, -- 分组名称，不能为空
                               created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP, -- 记录创建时间，默认为当前时间
                               updated_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP ON UPDATE CURRENT_TIMESTAMP, -- 记录更新时间，在更新时自动设置为当前时间
                               FOREIGN KEY (user_id) REFERENCES users(id) -- 外键，引用users表中的id
);

-- 好友关系表
CREATE TABLE friendships (
                             id INT AUTO_INCREMENT PRIMARY KEY, -- 好友关系ID，自增主键
                             user_id INT NOT NULL, -- 用户ID，不能为空
                             friend_id INT NOT NULL, -- 好友的用户ID，不能为空
                             status ENUM('pending', 'accepted', 'blocked') DEFAULT 'pending', -- 好友关系状态，默认为'pending'
                             group_id INT, -- 分组ID，允许为空
                             remark VARCHAR(255), -- 好友备注，允许为空,默认的话采用用户名
                             created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP, -- 记录创建时间，默认为当前时间
                             updated_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP ON UPDATE CURRENT_TIMESTAMP, -- 记录更新时间，在更新时自动设置为当前时间
                             FOREIGN KEY (user_id) REFERENCES users(id), -- 外键，引用users表中的id
                             FOREIGN KEY (friend_id) REFERENCES users(id), -- 外键，引用users表中的id
                             FOREIGN KEY (group_id) REFERENCES friend_groups(id) -- 外键，引用friend_groups表中的id
);

-- 群组表
CREATE TABLE `groups` (
                          id INT AUTO_INCREMENT PRIMARY KEY, -- 群组ID，自增主键
                          name VARCHAR(255) NOT NULL, -- 群组名称，不能为空
                          owner_id INT NOT NULL, -- 群主的用户ID，不能为空
                          created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP, -- 记录创建时间，默认为当前时间
                          updated_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP ON UPDATE CURRENT_TIMESTAMP, -- 记录更新时间，在更新时自动设置为当前时间
                          FOREIGN KEY (owner_id) REFERENCES users(id) -- 外键，引用users表中的id
);

-- 消息表
CREATE TABLE messages (
                          id INT AUTO_INCREMENT PRIMARY KEY, -- 消息ID，自增主键
                          sender_id INT NOT NULL, -- 发送者ID，不能为空
                          receiver_user_id INT, -- 接收者的用户ID，用于私聊消息
                          receiver_group_id INT, -- 接收者的群组ID，用于群聊消息
                          content TEXT, -- 消息内容，允许为空
                          message_type ENUM('text', 'image', 'video', 'file') DEFAULT 'text', -- 消息类型，默认为'text'
                          created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP, -- 记录创建时间，默认为当前时间
                          updated_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP ON UPDATE CURRENT_TIMESTAMP, -- 记录更新时间，在更新时自动设置为当前时间
                          FOREIGN KEY (sender_id) REFERENCES users(id), -- 外键，引用users表中的id
                          FOREIGN KEY (receiver_user_id) REFERENCES users(id), -- 外键，引用users表中的id
                          FOREIGN KEY (receiver_group_id) REFERENCES `groups`(id), -- 外键，引用groups表中的id
                          CONSTRAINT chk_receiver CHECK (
                              (receiver_user_id IS NOT NULL AND receiver_group_id IS NULL) OR
                              (receiver_user_id IS NULL AND receiver_group_id IS NOT NULL)
                              ) -- 确保 receiver_user_id 和 receiver_group_id 不会同时为空，也不会同时有值
);

-- 消息已读状态表
CREATE TABLE message_read_status (
                                     id INT AUTO_INCREMENT PRIMARY KEY, -- 已读状态ID，自增主键
                                     message_id INT NOT NULL, -- 消息ID，不能为空，引用messages表中的id
                                     user_id INT NOT NULL, -- 用户ID，不能为空，引用users表中的id
                                     read_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP, -- 记录已读时间，默认为当前时间
                                     FOREIGN KEY (message_id) REFERENCES messages(id), -- 外键，引用messages表中的id
                                     FOREIGN KEY (user_id) REFERENCES users(id) -- 外键，引用users表中的id
);


-- 群组成员表
CREATE TABLE group_members (
                               id INT AUTO_INCREMENT PRIMARY KEY,                                          -- 群组成员关系ID，自增主键
                               group_id INT NOT NULL,                                                      -- 群组ID，不能为空
                               user_id INT NOT NULL,                                                       -- 用户ID，不能为空
                               role ENUM('admin', 'member') NOT NULL DEFAULT 'member',                     -- 成员类型，默认为普通成员
                               title VARCHAR(255),                                                         -- 成员在群组中的称号，允许为空
                               level INT NOT NULL DEFAULT 1,                                               -- 成员等级，默认为1
                               nickname VARCHAR(255),                                                      -- 用户在群组中的昵称，允许为空
                               created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,                             -- 记录创建时间，默认为当前时间
                               updated_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP ON UPDATE CURRENT_TIMESTAMP, -- 记录更新时间，在更新时自动设置为当前时间
                               FOREIGN KEY (group_id) REFERENCES `groups`(id),                             -- 外键，引用groups表中的id
                               FOREIGN KEY (user_id) REFERENCES users(id)                                  -- 外键，引用users表中的id
);

-- 通知表
CREATE TABLE notifications (
                               id INT AUTO_INCREMENT PRIMARY KEY, -- 通知ID，自增主键
                               user_id INT NOT NULL, -- 用户ID，不能为空
                               type ENUM('message', 'friend_request', 'group_request', 'other') DEFAULT 'message', -- 通知类型，默认为'message'
                               content TEXT, -- 通知内容，允许为空
                               is_read BOOLEAN DEFAULT FALSE, -- 通知是否已读，默认为未读
                               created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP, -- 记录创建时间，默认为当前时间
                               updated_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP ON UPDATE CURRENT_TIMESTAMP, -- 记录更新时间，在更新时自动设置为当前时间
                               FOREIGN KEY (user_id) REFERENCES users(id) -- 外键，引用users表中的id
);

-- 系统日志表
CREATE TABLE system_logs (
                             id INT AUTO_INCREMENT PRIMARY KEY, -- 系统日志ID，自增主键
                             log_type VARCHAR(255), -- 日志类型
                             message TEXT, -- 日志内容
                             created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP -- 记录创建时间，默认为当前时间
);