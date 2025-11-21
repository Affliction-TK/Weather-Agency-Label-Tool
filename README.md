# Weather Agency Label Tool

面向气象行业的图像标注与审核工具，支持批量图片上传、OCR 自动识别时间地点、人工校验标注、经纬度反查和监测站点推荐。前后端完全开源，便于无锡气象局或其他气象单位快速部署与二次开发。

## 功能特性
- **批量上传 + 自动 OCR**：前端拖拽上传，后端 `ProcessImageOCR` 利用通义千问多模态模型提取拍摄时间与地点，并记录标准化结果。
- **标注工作台**：`AnnotationForm.svelte` 预填 OCR 结果，可一键调用 `/api/geocode` 获取经纬度，并根据经纬度推荐最近站点。
- **状态分组列表**：`ImageList.svelte` 按「未标注 / 已标注」分组，含缩略图、搜索过滤与展开折叠记忆。
- **站点/地理信息服务**：后台内置站点表，`/api/stations/nearest` 使用哈弗辛公式查找最近站点；`/api/geocode` 代理百度地图地理编码。
- **操作审计**：数据库 `annotations` 表保留创建/更新时间，便于追踪标注历史。

## 技术栈
- **前端**：Svelte 5 + Vite 7，纯前端构建，部署产物位于 `frontend/dist`。
- **后端**：Go 1.24，`gorilla/mux` 路由 + `database/sql` + `mysql` 驱动。
- **数据库**：MySQL 8（兼容 5.7+），`backend/schema.sql` 提供完整建表与示例站点。
- **AI/OCR**：通义千问视觉大模型（可配置），可根据需要替换。

## 目录结构
```
├── backend/              # Go 服务端
│   ├── main.go           # REST API、路由、存储逻辑
│   ├── ocr.go            # Qwen VLM OCR 管道
│   ├── schema.sql        # 数据库建表脚本
│   └── bin/server        # make build 后输出
├── frontend/             # Svelte 前端
│   ├── src/App.svelte    # 应用入口
│   └── src/lib/*.svelte  # 业务组件
├── uploads/              # 图片上传落地目录（运行时创建）
├── Makefile              # 常用命令封装
└── README.md
```

## 环境要求
- Go ≥ 1.24
- Node.js ≥ 20 & npm ≥ 10
- MySQL 实例（本地或云服务）
- 可选：通义千问 API Key（用于 OCR），百度地图开放平台 AK（用于地理编码）

## 快速开始
### 1. 克隆与依赖
```bash
git clone https://github.com/Affliction-TK/Weather-Agency-Label-Tool.git
cd Weather-Agency-Label-Tool
npm install --prefix frontend
```

### 2. 配置环境变量
在 `backend` 目录创建 `.env`（也可通过系统变量）：
```env
# 数据库
DB_HOST=127.0.0.1
DB_PORT=3306
DB_USER=weather_user
DB_PASSWORD=weather_password
DB_NAME=weather_label_db
DB_PARAMS=parseTime=true&charset=utf8mb4&loc=Asia%2FShanghai

# 运行配置
PORT=8080
UPLOAD_DIR=./uploads
STATIC_DIR=../frontend/dist

# OCR（可选）
QWEN_VLM_API_KEY=sk-xxx
QWEN_VLM_MODEL=qwen3-vl-plus
QWEN_VLM_ENABLE_THINKING=false
QWEN_VLM_THINKING_BUDGET=0

# 地理编码（可选）
BAIDU_MAP_AK=your-ak
LOCATION_PREFIX=江苏省无锡市
```

### 3. 初始化数据库
```bash
make init-db    # 或手动运行 backend/schema.sql
```

### 4. 本地开发
开两个终端：
```bash
make dev-backend    # 热编译 Go 服务，默认 8080 端口
make dev-frontend   # Vite 前端，默认 5173 端口
```
前端开发时需配置代理或将 `API_BASE` 指向 `http://localhost:8080/api`（源码已自动处理）。

### 5. 打包部署
```bash
make build      # 构建前后端
cd backend && ./bin/server
```
部署时可将 `frontend/dist` 静态资源放置在任意 Web 服务器，或复用 Go 静态文件托管（默认 `STATIC_DIR`）。

## 常用 Make 命令
| 命令             | 作用 |
|------------------|------|
| `make build`     | 构建前后端产物 |
| `make run`       | 构建后启动 Go 服务 |
| `make dev-backend` | Go 热运行（`go run .`）|
| `make dev-frontend` | Vite Dev Server |
| `make test`      | 运行 Go 单元测试 |
| `make clean`     | 清理 `backend/bin` 与 `frontend/dist` |
| `make init-db`   | 使用默认账号初始化数据库 |

## 配置项速查
| 变量 | 说明 | 默认 |
|------|------|------|
| `PORT` | 后端监听端口 | `8080` |
| `UPLOAD_DIR` | 图片存储目录 | `./uploads` |
| `STATIC_DIR` | 前端静态资源目录 | `../frontend/dist` |
| `DB_*` / `DB_DSN` | MySQL 连接信息 | 参考 `.env` |
| `DB_MAX_OPEN_CONNS` | 最大连接数 | `10` |
| `DB_MAX_IDLE_CONNS` | 空闲连接 | `5` |
| `QWEN_*` | 通义千问配置 | 可选 |
| `BAIDU_MAP_AK` | 百度 Maps AK | 必填以启用 `/api/geocode` |
| `LOCATION_PREFIX` | 地址前缀（地理编码上下文） | 空字符串 |

## 数据库模型
见 `backend/schema.sql`：
- `stations`：监测站点，保存经纬度与别名。
- `images`：上传图片及 OCR 结果。`annotated` 标记人工标注状态，`is_standard` 表示 OCR 是否同时识别到时间+地点。
- `annotations`：标注结果，唯一关联 `image_id`，含天气类型、严重程度、观测时间、地点、经纬度及站点。

## API 说明（`/api` 前缀）
| 方法 | 路径 | 描述 |
|------|------|------|
| `GET` | `/stations` | 获取所有站点列表 |
| `GET` | `/stations/nearest?longitude=&latitude=` | 基于经纬度返回最近站点 |
| `GET` | `/images` | 获取图片列表（含 OCR 字段）|
| `GET` | `/images/{id}` | 返回图片详情 + 标注（若存在）|
| `DELETE` | `/images/{id}` | 删除未标注图片（已标注会被拒绝）|
| `POST` | `/upload` | 上传图片并触发 OCR |
| `POST` | `/annotations` | 新增或更新标注，自动写入 `annotations`，并将图片标记为已标注 |
| `DELETE` | `/annotations/{id}` | 删除标注，同时重置图片状态 |
| `POST` | `/geocode` | 通过百度 API 将地点转换为经纬度 |
| `GET` | `/images/{filename}` | 静态图片访问（非 `/api` 前缀）|

## 后端关键函数
| 函数 | 文件 | 说明 |
|------|------|------|
| `initDB()` | `backend/main.go` | 读取 DSN、建立 MySQL 连接并配置连接池。|
| `uploadImage()` | `backend/main.go` | 负责接收 `multipart/form-data`、保存至 `UPLOAD_DIR`、调用 `ProcessImageOCR` 并写入 `images` 表。|
| `ProcessImageOCR()` | `backend/ocr.go` | 构建 Qwen VLM 请求，解析 JSON 回包，标准化时间 (`normalizeTime`) 与地点 (`cleanLocationText`)，决定 `is_standard`。|
| `geocodeAddress()` | `backend/main.go` | 调用百度地理编码，自动添加 `LOCATION_PREFIX`，并对经纬度四舍五入到 5 位小数。|
| `findNearestStation()` | `backend/main.go` | 使用哈弗辛公式在 `stations` 表中选择最近站点，为前端自动推荐提供数据。|
| `createAnnotation()` | `backend/main.go` | 先查重，存在则更新，不存在则插入；随后将对应图片 `annotated` 字段置为 `TRUE`。|
| `deleteImage()` / `deleteAnnotation()` | `backend/main.go` | 确保业务约束（已标注图片不可删、删除标注需同步重置图片状态）。|

## 前端核心模块
- `src/App.svelte`：顶层状态管理，负责加载站点/图片、切换标注与上传 tab、触发模态框。
- `src/lib/ImageList.svelte`：带缩略图、搜索与折叠记忆的图片列表组件，按标注状态分组。
- `src/lib/AnnotationForm.svelte`：标注表单，包含 OCR 预填、地理编码按钮、最近站点推荐、删除标注/图片逻辑。
- `src/lib/UploadTab.svelte`：文件拖拽上传、去重、批量上传进度提示。
- `src/lib/toastStore.js` + `Toast.svelte`：全局提示系统，支持 success/error/warning。

## 二次开发指南
### 后端扩展
1. **新增 API**：在 `backend/main.go` 中通过 `api.HandleFunc` 注册，路由统一挂载在 `/api`。
2. **数据结构**：如需扩展 `images` / `annotations` 字段，先修改 `schema.sql`，再更新 `Image`、`Annotation` struct 与对应 SQL。
3. **OCR/外部服务**：`ocr.go` 中的 `QwenVLMClient` 可替换为其他供应商，保持 `OCRResult` 输出即可。若新增字段，可在 `uploadImage` 中扩展持久化。
4. **配置**：新增环境变量时建议使用 `getEnv` / `getEnvInt` 封装，保持默认值清晰。
5. **错误处理**：API 返回 `http.Error`，并在日志中记录详细错误，方便排查。建议新增 handler 时遵循相同模式。

### 前端扩展
1. **状态管理**：当前使用 Svelte 原生响应式语法。如需全局状态，可在 `src/lib` 中新增 store。
2. **API Base**：`App.svelte`/`AnnotationForm.svelte`/`UploadTab.svelte`/`ImageList.svelte` 通过 `window.location.hostname` 自动选择 `http://localhost:8080` 或相对路径，部署到同域后无需修改。跨域部署可替换为 `.env` 注入配置。
3. **UI 主题**：`app.css` 定义了大量 CSS 变量（`--primary-color` 等），可集中修改主题色。
4. **表单字段**：在 `AnnotationForm.svelte` 中扩展 `formData` 并调整 `handleSubmit()` 序列化逻辑，同时更新后端 `createAnnotation` 的 SQL。
5. **国际化**：界面文本集中在各组件内，可配合 Svelte store/字典方案实现多语言。

### 工作流建议
- 修改数据库后运行 `make init-db` 或手动迁移，保持结构同步。
- 后端新增逻辑后运行 `make test`，目前包含 `main_test.go`/`ocr_test.go`。
- 前端改动建议运行 `npm run build` 验证产物可用；生产上线前执行 `make build`。

## 测试
```bash
make test          # 运行全部 Go 单元测试
npm run test       # 若未来添加前端测试，可在 frontend 目录运行
```
当前仓库主要包含 Go 端测试（`backend/main_test.go`, `backend/ocr_test.go`）。可按需补充前端 e2e/单元测试，保持 README 更新。

## 故障排查
- **数据库连接失败**：确认 `DB_*` 配置与 `schema.sql` 已初始化；必要时开启 `DB_DSN` 直连。
- **OCR 未生效**：检查 `QWEN_VLM_API_KEY` 是否配置；未配置时后端会记录 warning 并默认 `is_standard=false`。
- **地理编码失败**：确保 `BAIDU_MAP_AK` 有效、`LOCATION_PREFIX` 符合实际区域；接口错误会返回中文提示和状态码。
- **图片删除受阻**：只有未标注且无标注记录的图片可被删除，如需强制删除需同时移除 annotations 记录。

至此，README 已覆盖部署、使用与二次开发要点。如需更多帮助，可直接查看对应源码文件或提交 Issue。祝开发顺利！
