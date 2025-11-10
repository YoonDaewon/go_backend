# Go Backend API Server

Clean Architecture 패턴을 따르는 RESTful API 서버입니다. PostgreSQL, MongoDB, Redis를 지원합니다.

## 프로젝트 구조

```
go_backend/
├── main.go                    # 애플리케이션 진입점
├── config/                    # 설정 관리
│   └── config.go
├── database/                  # 데이터베이스 연결 관리
│   ├── database.go           # 통합 연결 관리
│   ├── postgres.go           # PostgreSQL 연결
│   ├── mongodb.go            # MongoDB 연결
│   └── redis.go              # Redis 연결
├── router/                    # 라우팅 설정
│   └── router.go
├── controller/                # HTTP 요청/응답 처리
│   └── user_controller.go
├── usecase/                   # 비즈니스 로직
│   └── user_usecase.go
├── repository/                # 데이터 접근 계층
│   ├── user_repository.go    # 인터페이스 및 인메모리 구현
│   ├── postgres_user_repository.go
│   ├── mongo_user_repository.go
│   └── redis_cache.go        # Redis 캐싱
└── model/                     # 도메인 모델
    └── user.go
```

## 기능

- ✅ Clean Architecture 패턴 (Router → Controller → Usecase → Repository)
- ✅ PostgreSQL 지원 (GORM 사용)
- ✅ MongoDB 지원
- ✅ Redis 지원 (캐싱)
- ✅ 환경 변수 기반 설정
- ✅ Graceful shutdown
- ✅ 자동 마이그레이션 (PostgreSQL)

## 설치 및 실행

### 1. 의존성 설치

```bash
go mod download
```

### 2. 환경 변수 설정

`.env.example` 파일을 참고하여 `.env` 파일을 생성하세요:

```bash
cp .env.example .env
```

`.env` 파일 예시:

```env
# Server Configuration
SERVER_PORT=8080
ENV=development

# PostgreSQL Configuration
POSTGRES_HOST=localhost
POSTGRES_PORT=5432
POSTGRES_USER=postgres
POSTGRES_PASSWORD=postgres
POSTGRES_DB=go_backend
POSTGRES_SSLMODE=disable

# MongoDB Configuration
MONGODB_URI=mongodb://localhost:27017
MONGODB_DB=go_backend
MONGODB_USERNAME=
MONGODB_PASSWORD=

# Redis Configuration
REDIS_HOST=localhost
REDIS_PORT=6379
REDIS_PASSWORD=
REDIS_DB=0

# Database Type Selection (postgres, mongodb, or leave empty for auto)
DB_TYPE=postgres
```

### 3. 데이터베이스 실행 (Docker 예시)

```bash
# PostgreSQL
docker run --name postgres -e POSTGRES_PASSWORD=postgres -e POSTGRES_DB=go_backend -p 5432:5432 -d postgres

# MongoDB
docker run --name mongodb -p 27017:27017 -d mongo

# Redis
docker run --name redis -p 6379:6379 -d redis
```

### 4. 서버 실행

```bash
go run main.go
```

또는 빌드 후 실행:

```bash
go build -o bin/server
./bin/server
```

서버는 `http://localhost:8080`에서 실행됩니다.

## API 엔드포인트

### Health Check
- `GET /healthcheck` - 서버 상태 확인

### Users
- `POST /api/v1/users` - 사용자 생성
- `GET /api/v1/users` - 모든 사용자 조회
- `GET /api/v1/users/:id` - 특정 사용자 조회
- `PUT /api/v1/users/:id` - 사용자 수정
- `DELETE /api/v1/users/:id` - 사용자 삭제

## 데이터베이스 선택

`DB_TYPE` 환경 변수로 사용할 데이터베이스를 선택할 수 있습니다:

- `postgres` - PostgreSQL 사용
- `mongodb` - MongoDB 사용
- 설정하지 않으면 PostgreSQL이 연결되어 있으면 PostgreSQL, 없으면 인메모리 저장소 사용

## 사용 예시

### 사용자 생성

```bash
curl -X POST http://localhost:8080/api/v1/users \
  -H "Content-Type: application/json" \
  -d '{
    "name": "John Doe",
    "email": "john@example.com"
  }'
```

### 사용자 조회

```bash
curl http://localhost:8080/api/v1/users/1
```

### 모든 사용자 조회

```bash
curl http://localhost:8080/api/v1/users
```

### 사용자 수정

```bash
curl -X PUT http://localhost:8080/api/v1/users/1 \
  -H "Content-Type: application/json" \
  -d '{
    "name": "Jane Doe",
    "email": "jane@example.com"
  }'
```

### 사용자 삭제

```bash
curl -X DELETE http://localhost:8080/api/v1/users/1
```

## 개발 가이드

### 새로운 엔티티 추가하기

1. `model/` 디렉토리에 모델 정의
2. `repository/` 디렉토리에 Repository 인터페이스 및 구현체 추가
3. `usecase/` 디렉토리에 Usecase 추가
4. `controller/` 디렉토리에 Controller 추가
5. `router/router.go`에 라우트 추가

### Redis 캐싱 사용하기

```go
import "go_backend/repository"

cache := repository.NewRedisCache()
ctx := context.Background()

// 캐시에 저장
cache.SetUser(ctx, user, 1*time.Hour)

// 캐시에서 조회
cachedUser, err := cache.GetUser(ctx, userID)
```

## 기술 스택

- **Framework**: Gin
- **ORM**: GORM (PostgreSQL)
- **MongoDB Driver**: Official MongoDB Go Driver
- **Redis Client**: go-redis
- **Configuration**: godotenv

## 라이선스

MIT

