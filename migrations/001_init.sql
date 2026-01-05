-- 初始化数据库结构（严格按照《数据库设计文档》）
CREATE TABLE "users" (
    id                  BIGINT PRIMARY KEY,
    name                VARCHAR(50) NOT NULL,
    phone               VARCHAR(20) NOT NULL,
    password            VARCHAR(255) NOT NULL,
    status              SMALLINT NOT NULL DEFAULT 0,
    created_at          TIMESTAMP WITH TIME ZONE NOT NULL DEFAULT CURRENT_TIMESTAMP,
    updated_at          TIMESTAMP WITH TIME ZONE NOT NULL DEFAULT CURRENT_TIMESTAMP,
    deleted_at          TIMESTAMP WITH TIME ZONE
);

CREATE UNIQUE INDEX uk_users_shop_phone ON "users"(phone) WHERE deleted_at IS null;
CREATE INDEX idx_users_created_at ON "users"(created_at);

COMMENT ON TABLE "users" IS '用户表（包含用户账号信息）';
COMMENT ON COLUMN "users".id IS '用户ID（Snowflake生成，全局唯一）';
COMMENT ON COLUMN "users".name IS '用户名（对外展示使用）';
COMMENT ON COLUMN "users".phone IS '手机号（作为登录账号，非删除状态下全局唯一）';
COMMENT ON COLUMN "users".password IS '密码（6位以上位字符）';
COMMENT ON COLUMN "users".status IS '状态：0-正常, 1-停用';