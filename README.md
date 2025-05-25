# TypeSpec Hello World - ECサイトAPI

このプロジェクトは、TypeSpecを使用してECサイトのREST APIを定義し、TypeScriptとGoの実装を生成するサンプルプロジェクトです。

## 技術スタック

- **API定義**: [TypeSpec](https://typespec.io/)
- **API仕様**: OpenAPI 3.0/3.1
- **Go実装**: 
  - コード生成: [oapi-codegen](https://github.com/deepmap/oapi-codegen)
  - HTTPライブラリ: Go標準パッケージ
- **TypeScript実装**:
  - 型生成: [openapi-typescript](https://github.com/drwpow/openapi-typescript)
  - サーバーフレームワーク: [Hono](https://hono.dev/)

## プロジェクト構成

```
/typespec/          # TypeSpec定義
  /models/          # 共通モデル定義
  /services/        # サービス定義
  main.tsp          # エントリーポイント
  
/openapi/           # 生成されたOpenAPI仕様
  openapi.yaml      

/typescript/        # TypeScriptプロジェクト
  /generated/       # openapi-typescriptで生成された型定義
  /src/            
    /routes/        # Honoのルート実装
    /services/      # ビジネスロジック
    
/go/               # Goプロジェクト
  /generated/       # oapi-codegenで生成されたコード
  /cmd/api/         # メインアプリケーション
  /internal/        # 内部実装
```

## APIリソース

ECサイトのAPIとして以下のリソースを実装します：

- **Products** - 商品管理
- **Categories** - カテゴリ管理
- **Cart** - ショッピングカート
- **Orders** - 注文管理
- **Users** - ユーザー管理

## セットアップ

### 前提条件

- Node.js 18以上
- Go 1.21以上
- pnpm（TypeScriptプロジェクト用）

### 初期セットアップ

```bash
# TypeSpecの依存関係をインストール
cd typespec
npm install

# TypeScriptプロジェクトのセットアップ
cd ../typescript
pnpm install

# Goプロジェクトのセットアップ
cd ../go
go mod init github.com/blck-snwmn/hello-typespec
```

## ビルド手順

### 1. TypeSpecからOpenAPIを生成

```bash
cd typespec
npm run compile
```

### 2. OpenAPIからTypeScriptの型を生成

```bash
cd typescript
pnpm run generate
```

### 3. OpenAPIからGoのコードを生成

```bash
cd go
go generate ./...
```

## 開発サーバーの起動

### TypeScript

```bash
cd typescript
pnpm run dev
```

### Go

```bash
cd go
go run cmd/api/main.go
```

## コマンド一覧

### TypeSpec

- `npm run compile` - TypeSpecからOpenAPIを生成
- `npm run watch` - ファイル変更を監視して自動コンパイル
- `npm run format` - TypeSpecファイルのフォーマット

### TypeScript

- `pnpm run generate` - OpenAPIから型定義を生成
- `pnpm run dev` - 開発サーバーを起動
- `pnpm run build` - プロダクションビルド
- `pnpm run type-check` - 型チェック

### Go

- `go generate ./...` - OpenAPIからコードを生成
- `go run cmd/api/main.go` - サーバーを起動
- `go test ./...` - テストを実行

## 開発フロー

1. TypeSpecでAPI定義を修正
2. `npm run compile`でOpenAPIを生成
3. 各言語でコード生成コマンドを実行
4. 生成されたコードを基に実装を追加

## ライセンス

MIT