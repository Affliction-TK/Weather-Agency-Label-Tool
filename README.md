# Weather Agency Label Tool

气象局图片标注系统 - 用于标注监控摄像头图片的气象现象（大雾、结冰、积劳）

## 功能特点

- 📸 图片列表展示（包含缩略图、文件名、标注状态）
- 🏷️ 气象现象标注（大雾、结冰、积劳）
- 📊 严重等级选择（无、轻度、中度、重度）
- 📍 地点和经纬度信息录入
- 🗺️ 智能匹配最近的监测站点（150个站点）
- ⏰ 观测时间记录
- ⬆️ 图片上传功能
- ✅ 自动跳转到未标注图片
- 💾 MySQL数据持久化存储

## 快速开始

```bash
# 1. 克隆仓库
git clone https://github.com/Affliction-TK/Weather-Agency-Label-Tool.git
cd Weather-Agency-Label-Tool

# 2. 运行自动化设置
./setup.sh

# 3. 启动 MySQL（使用 Docker）
docker-compose up -d

# 4. 构建并运行
make build
./server
```

打开浏览器访问 http://localhost:8080

📖 **详细说明**: 查看 [快速开始指南](QUICKSTART.md)

## 文档

- **[快速开始指南](QUICKSTART.md)** - 5分钟快速入门
- **[API 文档](API.md)** - 完整的 REST API 参考
- **[部署指南](DEPLOYMENT.md)** - 生产环境部署说明
- **[系统架构](ARCHITECTURE.md)** - 架构设计和技术决策
- **[贡献指南](CONTRIBUTING.md)** - 如何参与项目开发

## 技术栈

- **后端**: Go 1.24 + Gorilla Mux
- **前端**: Svelte + Vite
- **数据库**: MySQL 8.0+
- **文件存储**: 本地文件系统

## 系统截图

### 主界面
- 左侧：图片列表（缩略图、文件名、标注状态）
- 右侧：标注表单或上传界面

### 标注功能
- 气象类型选择
- 严重等级选择
- 时间、地点、经纬度输入
- 自动推荐最近的监测站点

### 上传功能
- 拖拽上传或点击选择
- 支持多文件同时上传
- 上传进度显示

## 快速开始

### 前置要求

- Go 1.24 或更高版本
- Node.js 20+ 和 npm
- MySQL 8.0 或更高版本

### 1. 数据库设置

创建数据库并导入数据结构：

```bash
mysql -u root -p < schema.sql
```

或手动执行：

```sql
CREATE DATABASE weather_label_db CHARACTER SET utf8mb4 COLLATE utf8mb4_unicode_ci;
```

然后导入 `schema.sql` 文件。

### 2. 配置环境变量

复制示例配置并按需调整：

```bash
cp .env.example .env
```

编辑 `.env`，填写数据库和服务参数：

```
DB_HOST=127.0.0.1
DB_PORT=3306
DB_NAME=weather_label_db
DB_USER=weather_user
DB_PASSWORD=your_password
DB_PARAMS=parseTime=true&charset=utf8mb4&loc=Local
PORT=8080
UPLOAD_DIR=./uploads
STATIC_DIR=./frontend/dist
```

如需一次性指定完整 DSN，可设置 `DB_DSN` 覆盖上述字段。

### 3. 安装依赖

#### 后端依赖

```bash
go mod download
```

#### 前端依赖

```bash
cd frontend
npm install
```

### 4. 构建前端

```bash
cd frontend
npm run build
cd ..
```

### 5. 运行应用

```bash
go run main.go
```

服务器将在 `http://localhost:8080` 启动。

## 开发模式

### 后端开发

```bash
go run main.go
```

### 前端开发

在另一个终端：

```bash
cd frontend
npm run dev
```

前端开发服务器将在 `http://localhost:5173` 启动。

## 项目结构

```
.
├── main.go              # Go 后端服务器
├── go.mod               # Go 依赖管理
├── schema.sql           # MySQL 数据库架构
├── uploads/             # 上传图片存储目录
└── frontend/            # Svelte 前端应用
    ├── src/
    │   ├── App.svelte           # 主应用组件
    │   ├── lib/
    │   │   ├── ImageList.svelte      # 图片列表组件
    │   │   ├── AnnotationForm.svelte # 标注表单组件
    │   │   └── UploadTab.svelte      # 上传标签页组件
    │   └── main.js
    └── dist/            # 构建输出目录
```

## API 接口

### 获取监测站点列表
```
GET /api/stations
```

### 获取最近的监测站点
```
GET /api/stations/nearest?longitude=<lon>&latitude=<lat>
```

### 获取图片列表
```
GET /api/images
```

### 获取图片详情
```
GET /api/images/:id
```

### 保存标注
```
POST /api/annotations
Content-Type: application/json

{
  "image_id": 1,
  "category": "大雾",
  "severity": "中度",
  "observation_time": "2024-01-01T12:00:00Z",
  "location": "北京市朝阳区",
  "longitude": 116.407396,
  "latitude": 39.904211,
  "station_id": 1
}
```

### 上传图片
```
POST /api/upload
Content-Type: multipart/form-data

image: <file>
```

### 获取图片文件
```
GET /images/:filename
```

## 数据库架构

### stations 表
- `id`: 站点ID
- `name`: 站点名称
- `longitude`: 经度
- `latitude`: 纬度

### images 表
- `id`: 图片ID
- `filename`: 文件名
- `filepath`: 文件路径
- `uploaded_at`: 上传时间
- `annotated`: 是否已标注

### annotations 表
- `id`: 标注ID
- `image_id`: 图片ID
- `category`: 气象类型（大雾/结冰/积劳）
- `severity`: 严重等级（无/轻度/中度/重度）
- `observation_time`: 观测时间
- `location`: 地点
- `longitude`: 经度
- `latitude`: 纬度
- `station_id`: 监测站点ID

## 使用说明

1. **查看图片列表**：左侧边栏显示所有图片，未标注的图片会优先显示
2. **标注图片**：选择气象类型、严重等级，填写时间、地点、经纬度信息
3. **选择监测点**：输入经纬度后系统会自动推荐最近的监测站点
4. **保存标注**：点击"保存标注"按钮，系统会自动跳转到下一张未标注的图片
5. **上传新图片**：点击"上传图片"标签页，拖拽或选择图片进行上传

## License

MIT
