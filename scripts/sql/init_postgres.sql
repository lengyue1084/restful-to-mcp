-- ============================================================
-- RESTful-to-MCP PostgreSQL 初始化脚本
-- 说明：项目已支持 GORM AutoMigrate 自动建表，
--       此脚本供手动初始化或参考表结构使用。
-- ============================================================

-- ----------------------------
-- 序列
-- ----------------------------
CREATE SEQUENCE IF NOT EXISTS "public"."mcp_connect_token_id_seq" INCREMENT 1 MINVALUE 1 START 1 CACHE 1;
CREATE SEQUENCE IF NOT EXISTS "public"."mcp_file_id_seq" INCREMENT 1 MINVALUE 1 START 1 CACHE 1;
CREATE SEQUENCE IF NOT EXISTS "public"."mcp_server_id_seq" INCREMENT 1 MINVALUE 1 START 1 CACHE 1;
CREATE SEQUENCE IF NOT EXISTS "public"."mcp_tools_id_seq" INCREMENT 1 MINVALUE 1 START 1 CACHE 1;

-- ----------------------------
-- mcp_connect_token - MCP 连接令牌表
-- ----------------------------
CREATE TABLE IF NOT EXISTS "public"."mcp_connect_token" (
  "id"              int8         NOT NULL DEFAULT nextval('mcp_connect_token_id_seq'::regclass),
  "created_at"      timestamptz(6),
  "updated_at"      timestamptz(6),
  "deleted_at"      timestamptz(6),
  "mcp_server_uuid" varchar(36)  NOT NULL DEFAULT '',
  "mcp_server_id"   int8         NOT NULL,
  "mcp_server_name" varchar(500) NOT NULL DEFAULT '',
  "connect_token"   varchar(36)  NOT NULL DEFAULT '',
  CONSTRAINT "mcp_connect_token_pkey" PRIMARY KEY ("id")
);

COMMENT ON COLUMN "public"."mcp_connect_token"."mcp_server_uuid" IS '关联的MCP服务器UUID';
COMMENT ON COLUMN "public"."mcp_connect_token"."mcp_server_id"   IS '关联的MCP服务器ID';
COMMENT ON COLUMN "public"."mcp_connect_token"."mcp_server_name" IS '服务器名称';
COMMENT ON COLUMN "public"."mcp_connect_token"."connect_token"   IS '连接令牌';

-- ----------------------------
-- mcp_file - 文件记录表
-- ----------------------------
CREATE TABLE IF NOT EXISTS "public"."mcp_file" (
  "id"          int8         NOT NULL DEFAULT nextval('mcp_file_id_seq'::regclass),
  "created_at"  timestamptz(6),
  "updated_at"  timestamptz(6),
  "deleted_at"  timestamptz(6),
  "md5"         varchar(36)  NOT NULL DEFAULT '',
  "name"        varchar(500) NOT NULL DEFAULT '',
  "source_name" varchar(500) NOT NULL DEFAULT '',
  "description" text         DEFAULT '',
  "suffix"      varchar(255) NOT NULL DEFAULT 'yaml',
  CONSTRAINT "mcp_file_pkey" PRIMARY KEY ("id")
);

-- ----------------------------
-- mcp_server - MCP 服务器表
-- ----------------------------
CREATE TABLE IF NOT EXISTS "public"."mcp_server" (
  "id"              int8         NOT NULL DEFAULT nextval('mcp_server_id_seq'::regclass),
  "created_at"      timestamptz(6),
  "updated_at"      timestamptz(6),
  "deleted_at"      timestamptz(6),
  "uuid"            varchar(36)  NOT NULL DEFAULT '',
  "name"            varchar(500) NOT NULL DEFAULT '',
  "description"     text         DEFAULT '',
  "urls"            varchar(255) NOT NULL DEFAULT '',
  "all_tools"       text         DEFAULT '',
  "version"         varchar(20)  DEFAULT 'v1.0.0',
  "mcp_server_type" int2         DEFAULT 1,
  "have_tools"      int2         DEFAULT 1,
  "is_auth"         int2         DEFAULT 1,
  "service_token"   text         DEFAULT '',
  "platform_token"  text         DEFAULT '',
  "security"        text         DEFAULT '',
  "status"          int2         DEFAULT 1,
  "serial_number"   varchar(36)  DEFAULT '',
  "source"          int2         DEFAULT 1,
  "header"          text         DEFAULT '',
  CONSTRAINT "mcp_server_pkey" PRIMARY KEY ("id")
);

CREATE INDEX  IF NOT EXISTS "idx_deleted_at"   ON "public"."mcp_server" USING btree ("deleted_at" ASC NULLS LAST);
CREATE UNIQUE INDEX IF NOT EXISTS "idx_uuid_unique" ON "public"."mcp_server" USING btree ("uuid" ASC NULLS LAST);

COMMENT ON COLUMN "public"."mcp_server"."uuid"            IS '服务器唯一标识';
COMMENT ON COLUMN "public"."mcp_server"."name"            IS '服务器名称';
COMMENT ON COLUMN "public"."mcp_server"."description"     IS '服务器详细描述';
COMMENT ON COLUMN "public"."mcp_server"."urls"            IS '服务器访问地址';
COMMENT ON COLUMN "public"."mcp_server"."all_tools"       IS '允许使用的工具列表，JSON格式存储';
COMMENT ON COLUMN "public"."mcp_server"."version"         IS 'OpenAPI版本号';
COMMENT ON COLUMN "public"."mcp_server"."mcp_server_type" IS '服务器类型 1:openapi 2:grpc';
COMMENT ON COLUMN "public"."mcp_server"."have_tools"      IS '是否支持工具 0:未知 1:不支持 2:支持';
COMMENT ON COLUMN "public"."mcp_server"."is_auth"         IS '认证类型 0:未知 1:不开启 2:平台授权 3:service授权 4:全部授权';
COMMENT ON COLUMN "public"."mcp_server"."service_token"   IS '服务认证Token';
COMMENT ON COLUMN "public"."mcp_server"."platform_token"  IS '平台认证Token';
COMMENT ON COLUMN "public"."mcp_server"."security"        IS '认证配置信息，JSON格式';
COMMENT ON COLUMN "public"."mcp_server"."status"          IS '状态 0:未知 1:未设置token 2:已设置token 3:正常工作';
COMMENT ON COLUMN "public"."mcp_server"."serial_number"   IS '服务序列号';
COMMENT ON COLUMN "public"."mcp_server"."source"          IS '数据来源 1:OpenAPI文档 2:表单创建';
COMMENT ON COLUMN "public"."mcp_server"."header"          IS '请求头信息，JSON格式存储';

-- ----------------------------
-- mcp_tools - MCP 工具表
-- ----------------------------
CREATE TABLE IF NOT EXISTS "public"."mcp_tools" (
  "id"               int8         NOT NULL DEFAULT nextval('mcp_tools_id_seq'::regclass),
  "created_at"       timestamptz(6),
  "updated_at"       timestamptz(6),
  "deleted_at"       timestamptz(6),
  "uuid"             varchar(36)  NOT NULL DEFAULT '',
  "mcp_server_id"    int8         NOT NULL,
  "mcp_server_uuid"  varchar(36)  NOT NULL DEFAULT '',
  "name"             varchar(500) NOT NULL DEFAULT '',
  "description"      text         DEFAULT '',
  "mcp_server_type"  int2         DEFAULT 1,
  "method"           varchar(10)  DEFAULT '',
  "endpoint"         varchar(255) DEFAULT '',
  "headers"          text         DEFAULT '',
  "args"             text         DEFAULT '',
  "request_body"     text         DEFAULT '',
  "response_body"    text         DEFAULT '',
  "tool_schema"      text         DEFAULT '',
  "annotations"      text         DEFAULT '',
  "security"         text         DEFAULT '',
  "is_auth"          int2         DEFAULT 1,
  "auth_mode"        varchar(20)  DEFAULT '',
  "is_platform_auth" int2         DEFAULT 1,
  "is_show"          int2         DEFAULT 1,
  "serial_number"    varchar(36)  DEFAULT '',
  "is_repeat"        int2         DEFAULT 1,
  "is_test"          int2         DEFAULT 1,
  CONSTRAINT "mcp_tools_pkey" PRIMARY KEY ("id")
);

COMMENT ON COLUMN "public"."mcp_tools"."uuid"             IS '工具唯一标识';
COMMENT ON COLUMN "public"."mcp_tools"."mcp_server_id"    IS '关联的MCP服务器ID';
COMMENT ON COLUMN "public"."mcp_tools"."mcp_server_uuid"  IS '关联的MCP服务器UUID';
COMMENT ON COLUMN "public"."mcp_tools"."name"             IS '工具名称';
COMMENT ON COLUMN "public"."mcp_tools"."description"      IS '工具描述信息';
COMMENT ON COLUMN "public"."mcp_tools"."mcp_server_type"  IS '服务器类型 1:openapi 2:grpc';
COMMENT ON COLUMN "public"."mcp_tools"."method"           IS 'HTTP请求方法(GET/POST/PUT/DELETE等)';
COMMENT ON COLUMN "public"."mcp_tools"."endpoint"         IS 'API端点地址';
COMMENT ON COLUMN "public"."mcp_tools"."headers"          IS '请求头信息，JSON格式存储';
COMMENT ON COLUMN "public"."mcp_tools"."args"             IS '参数配置，JSON格式存储';
COMMENT ON COLUMN "public"."mcp_tools"."request_body"     IS '请求体模板';
COMMENT ON COLUMN "public"."mcp_tools"."response_body"    IS '响应体模板';
COMMENT ON COLUMN "public"."mcp_tools"."tool_schema"      IS '工具参数Schema定义';
COMMENT ON COLUMN "public"."mcp_tools"."annotations"      IS '工具注解信息，JSON格式存储';
COMMENT ON COLUMN "public"."mcp_tools"."security"         IS '认证配置信息';
COMMENT ON COLUMN "public"."mcp_tools"."is_auth"          IS '是否需要认证 0:未知 1:不需要 2:需要';
COMMENT ON COLUMN "public"."mcp_tools"."auth_mode"        IS '认证模式';
COMMENT ON COLUMN "public"."mcp_tools"."is_platform_auth" IS '是否平台认证';
COMMENT ON COLUMN "public"."mcp_tools"."is_show"          IS '是否显示';
COMMENT ON COLUMN "public"."mcp_tools"."serial_number"    IS '服务序列号';
COMMENT ON COLUMN "public"."mcp_tools"."is_repeat"        IS '是否重复';
COMMENT ON COLUMN "public"."mcp_tools"."is_test"          IS '是否已测试';

-- ----------------------------
-- 序列归属
-- ----------------------------
ALTER SEQUENCE "public"."mcp_connect_token_id_seq" OWNED BY "public"."mcp_connect_token"."id";
ALTER SEQUENCE "public"."mcp_file_id_seq"          OWNED BY "public"."mcp_file"."id";
ALTER SEQUENCE "public"."mcp_server_id_seq"        OWNED BY "public"."mcp_server"."id";
ALTER SEQUENCE "public"."mcp_tools_id_seq"         OWNED BY "public"."mcp_tools"."id";
