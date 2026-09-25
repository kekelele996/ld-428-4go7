# ArtVault · 艺术作品集与画廊管理平台（全栈Web应用）

**项目类型标签：全栈Web应用**

ArtVault 面向独立艺术家和小型画廊，提供作品管理、线上展览策划、观众互动（点赞/评论/收藏）和艺术家个人主页，帮助艺术家展示作品、画廊策划展览并与观众互动。

## 快速启动（Docker Compose 一键部署）

```bash
cp .env.example .env
docker compose config --quiet
docker compose up -d --build
```

访问地址：

- 前端：http://localhost:18808
- 后端健康检查：http://localhost:19308/healthz
- API 前缀：/api/v1

停止并清理：

```bash
docker compose down -v --remove-orphans
```

## 主要功能

- **作品管理**：上传作品（标题/描述/媒介/尺寸/图片/标签/价格），草稿→发布→售出→归档状态流转，发布内容经审核后公开。
- **线上展览**：策展人创建展览（Solo/Group/Thematic/Permanent），收录作品，发布后对外展示并统计参观人数。
- **观众互动**：点赞、评论、收藏、分享，作品点赞/收藏/浏览计数实时更新。
- **艺术家主页**：艺术家名、简介、擅长媒介、代表作品、社交媒体链接、粉丝关注数。
- **创作工作台（Studio）**：我的作品管理（草稿/已发布/已售 Tab）、上传新作品、创建展览、数据统计。
- **RBAC 权限**：Admin/Curator/Artist/Viewer，Viewer 只能浏览和互动。
- **内容审核**：作品与展览提交审核（Approved/Rejected/Flagged），未审核内容不公开返回。

## 本地开发方式

后端：

```bash
cd backend
go mod tidy
go run ./cmd/server
```

构建：`go build ./...`（在 backend/ 下执行）

前端：

```bash
cd frontend
npm install
npm run dev
```

本地运行时需自行提供 MongoDB 7，并通过 `MONGO_URI` 环境变量注入连接串。

## 技术栈

| 层级 | 技术 |
| --- | --- |
| 前端 | React 18、TypeScript、Tailwind CSS、Headless UI、ECharts、Zustand、Vite |
| 后端 | Go 1.22 + Gin + MongoDB（官方 mongo-driver，repository 模式） |
| 数据库 | MongoDB 7 |
| 认证 | JWT（golang-jwt/v5）+ RBAC，密码 bcrypt |
| 部署 | Docker Compose、Nginx |

## 项目目录结构

```text
.
├── backend/
│   ├── cmd/server/main.go        # 程序入口
│   ├── internal/
│   │   ├── config/               # 配置（MongoDB/JWT/端口）
│   │   ├── database/             # MongoDB 连接
│   │   ├── model/                # 实体（bson 标签，按实体分文件）
│   │   ├── repository/           # mongo-driver 数据访问层
│   │   ├── service/              # 业务逻辑层
│   │   ├── handler/              # HTTP 处理函数
│   │   ├── router/               # 路由注册（按实体分文件）
│   │   ├── middleware/           # auth/rbac/内容审核/审计/错误处理/限流
│   │   ├── dto/                  # 请求响应结构体
│   │   ├── constants/            # 枚举/错误码/日志模板
│   │   └── util/                 # logger/jwt/app_error/ids
│   └── database/seeds/           # 种子数据
├── frontend/
│   ├── src/                      # React 页面/组件/状态/类型
│   ├── Dockerfile
│   └── nginx.conf
├── docker-compose.yml
├── .env.example
└── README.md
```

## 环境变量说明

| 变量 | 说明 | 默认值 |
| --- | --- | --- |
| `COMPOSE_PROJECT_NAME` | Compose 项目名/容器名前缀 | `artvault` |
| `MONGO_INITDB_ROOT_USERNAME` | MongoDB root 用户 | `artvault_root` |
| `MONGO_INITDB_ROOT_PASSWORD` | MongoDB root 密码 | `artvault_pwd` |
| `MONGO_DB_NAME` | MongoDB 数据库名 | `artvault` |
| `JWT_SECRET` | JWT 签名密钥 | `change_me...` |
| `FRONTEND_PORT` | 前端宿主端口 | `18808` |
| `BACKEND_PORT` | 后端宿主端口 | `19308` |
| `DB_PORT` | MongoDB 宿主端口 | `33328` |

## Docker 部署说明

- Compose 顶层声明 `name: artvault`，容器名带 `${COMPOSE_PROJECT_NAME:-artvault}` 前缀，可在任意目录名（含中文）下启动。
- 前端 Nginx 将 `/api/` 反向代理到 `http://backend:8080`，前端调用 `/api/v1/...`。
- MongoDB 数据使用命名卷 `db_data` 持久化；MongoDB healthcheck 就绪后后端启动，后端 healthy 后前端启动。
- 后端使用官方 `go.mongodb.org/mongo-driver` + repository 模式访问 MongoDB。
- 常见问题：端口冲突时修改 `.env` 端口后重新 `docker compose up -d`。

## API 清单（/api/v1）

| 方法 | 路径 | 说明 |
| --- | --- | --- |
| POST | /api/v1/auth/login | 登录获取 JWT |
| POST | /api/v1/auth/register | 注册 |
| GET | /api/v1/artworks | 作品列表（仅公开已审核） |
| GET | /api/v1/artworks/:id | 作品详情 |
| POST | /api/v1/artworks | 上传作品（Artist/Admin） |
| PATCH | /api/v1/artworks/:id/status | 发布/下架作品 |
| PATCH | /api/v1/artworks/:id/review | 内容审核（Admin/Curator） |
| GET | /api/v1/exhibitions | 展览列表 |
| GET | /api/v1/exhibitions/:id | 展览详情 |
| POST | /api/v1/exhibitions | 创建展览（Curator/Admin） |
| PATCH | /api/v1/exhibitions/:id/status | 发布展览 |
| POST | /api/v1/exhibitions/:id/artworks/:artworkId | 展览收录作品 |
| GET | /api/v1/artists | 艺术家列表 |
| GET | /api/v1/artists/:id | 艺术家主页 |
| PATCH | /api/v1/artists/me | 完善我的艺术家主页 |
| GET/POST | /api/v1/interactions | 互动列表/点赞评论收藏 |
| DELETE | /api/v1/interactions | 取消点赞/收藏 |
| GET | /api/v1/reviews | 审核日志（Admin/Curator） |
| GET | /api/v1/audit-logs | 操作日志（Admin/Curator） |

演示账号：`lin / Artist@123`（艺术家林知微）、`chen / Artist@123`（艺术家陈序）、`viewer / Viewer@123`（观众）、`admin / Admin@123`（管理员）。

## 枚举出现位置清单

**ArtworkStatus（Draft/Published/Sold/Archived）**
- 前端：`frontend/src/types/enums.ts`、`frontend/src/pages/Studio.tsx`、`frontend/src/components/common/StatusBadge.tsx`
- 后端：`backend/internal/constants/artwork.go`、`backend/internal/model/artwork.go`、`backend/internal/service/artwork_service.go`

**Medium（OilPainting/Watercolor/Acrylic/Ink/Sculpture/Photography/Digital/Mixed/Installation/Other）**
- 前端：`frontend/src/constants/mediumOptions.ts`、`frontend/src/types/enums.ts`、`frontend/src/pages/Studio.tsx`
- 后端：`backend/internal/constants/artwork.go`、`backend/internal/model/artwork.go`

**ExhibitionStatus（Planning/Active/Ended/Archived）**
- 前端：`frontend/src/types/enums.ts`、`frontend/src/components/common/StatusBadge.tsx`、`frontend/src/pages/ExhibitionDetail.tsx`
- 后端：`backend/internal/constants/exhibition.go`、`backend/internal/model/exhibition.go`、`backend/internal/service/exhibition_service.go`

**InteractionType（Like/Comment/Bookmark/Share）**
- 前端：`frontend/src/types/enums.ts`、`frontend/src/components/common/InteractionBar.tsx`、`frontend/src/hooks/useInteraction.ts`
- 后端：`backend/internal/constants/interaction.go`、`backend/internal/model/interaction.go`、`backend/internal/service/interaction_service.go`

## License

MIT
