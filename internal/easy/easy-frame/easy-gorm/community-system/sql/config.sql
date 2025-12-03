-- 用户表
CREATE TABLE users (
                       id BIGINT UNSIGNED PRIMARY KEY AUTO_INCREMENT,
                       username VARCHAR(50) UNIQUE NOT NULL,
                       email VARCHAR(100) NOT NULL,
                       password_hash VARCHAR(255) NOT NULL,
                       created_at DATETIME(3),
                       updated_at DATETIME(3),
                       deleted_at DATETIME(3),
                       INDEX idx_deleted_at (deleted_at)
);

-- 频道表
CREATE TABLE channels (
                          id BIGINT UNSIGNED PRIMARY KEY AUTO_INCREMENT,
                          name VARCHAR(100) NOT NULL,
                          description TEXT,
                          user_id BIGINT UNSIGNED NOT NULL,
                          created_at DATETIME(3),
                          updated_at DATETIME(3),
                          FOREIGN KEY (user_id) REFERENCES users(id)
);

-- 主题表
CREATE TABLE topics (
                        id BIGINT UNSIGNED PRIMARY KEY AUTO_INCREMENT,
                        title VARCHAR(200) NOT NULL,
                        channel_id BIGINT UNSIGNED NOT NULL,
                        user_id BIGINT UNSIGNED NOT NULL,
                        created_at DATETIME(3),
                        updated_at DATETIME(3),
                        FOREIGN KEY (channel_id) REFERENCES channels(id),
                        FOREIGN KEY (user_id) REFERENCES users(id)
);

-- 文章表
CREATE TABLE articles (
                          id BIGINT UNSIGNED PRIMARY KEY AUTO_INCREMENT,
                          content TEXT NOT NULL,
                          topic_id BIGINT UNSIGNED NOT NULL,
                          user_id BIGINT UNSIGNED NOT NULL,
                          created_at DATETIME(3),
                          updated_at DATETIME(3),
                          FOREIGN KEY (topic_id) REFERENCES topics(id),
                          FOREIGN KEY (user_id) REFERENCES users(id)
);

-- 评论表
CREATE TABLE comments (
                          id BIGINT UNSIGNED PRIMARY KEY AUTO_INCREMENT,
                          content TEXT NOT NULL,
                          article_id BIGINT UNSIGNED NOT NULL,
                          user_id BIGINT UNSIGNED NOT NULL,
                          parent_id BIGINT UNSIGNED,
                          created_at DATETIME(3),
                          updated_at DATETIME(3),
                          FOREIGN KEY (article_id) REFERENCES articles(id),
                          FOREIGN KEY (user_id) REFERENCES users(id),
                          FOREIGN KEY (parent_id) REFERENCES comments(id)
);
