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
  - ビルドツール: [Vite](https://vitejs.dev/)
  - テストフレームワーク: [Vitest](https://vitest.dev/)
- **CI/CD**: GitHub Actions
- **パッケージ管理**: pnpm workspaces (モノレポ構成)

## プロジェクト構成

```
/typespec/          # TypeSpec定義
  /models/          # 共通モデル定義
  /services/        # サービス定義
  main.tsp          # エントリーポイント
  
/openapi/           # 生成されたOpenAPI仕様
  openapi.yaml      

/typescript/        # TypeScriptプロジェクト
  /src/            
    /routes/        # Honoのルート実装
    /services/      # ビジネスロジック  
    /types/         # 型定義とエラーハンドリング
      api.ts        # 自動生成された型定義
      errors.ts     # エラーハンドリング統一実装
  /tests/          # テストコード
  /coverage/       # カバレッジレポート (gitignore)
  /dist/           # ビルド成果物 (gitignore)
    
/go/               # Goプロジェクト
  /generated/       # oapi-codegenで生成されたコード
  /cmd/api/         # メインアプリケーション
  /internal/        # 内部実装
    /handlers/      # HTTPハンドラー実装
    /storage/       # データストレージ層
  /tests/          # テストコード

/.github/workflows/ # CI/CD設定
  ci.yml           # PR/Push時の検証
  dependency-check.yml # セキュリティチェック
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
- pnpm 8以上（パッケージマネージャー）

### 初期セットアップ

```bash
# リポジトリをクローン
git clone https://github.com/blck-snwmn/hello-typespec.git
cd hello-typespec

# pnpm workspaceで全依存関係をインストール（ルートディレクトリから実行）
pnpm install

# Goの依存関係をインストール
cd go
go mod download
```

## ビルド手順

### 一括生成（推奨）

```bash
# ルートディレクトリから実行
pnpm generate:all
```

### 個別生成

#### 1. TypeSpecからOpenAPIを生成

```bash
# ルートから
pnpm compile:spec
# または
cd typespec && npm run compile
```

#### 2. OpenAPIからTypeScriptの型を生成

```bash
# ルートから
pnpm generate:typescript
# または
cd typescript && pnpm run generate
```

#### 3. OpenAPIからGoのコードを生成

```bash
# ルートから
pnpm generate:go
# または
cd go && go generate ./...
```

## 開発サーバーの起動

### TypeScript

```bash
# ルートから
pnpm -F typescript dev
# または
cd typescript && pnpm run dev
```

デフォルトでは http://localhost:3000 で起動します。

### Go

```bash
cd go
go run cmd/api/main.go
```

デフォルトでは http://localhost:8080 で起動します。

## OpenAPIドキュメント（Elements）

`elements/index.html` から、Stoplight Elements を使って `openapi/openapi.yaml` をブラウザ表示できます。

### 簡易サーバーで閲覧

ルートディレクトリで静的サーバーを起動し、`elements/index.html` にアクセスします。

```bash
# 推奨（スクリプト）
pnpm run docs:serve

# 代替（直接実行）
npx http-server -p 8081 -c-1
# または
pnpm dlx http-server -p 8081 -c-1
```

- ブラウザで http://127.0.0.1:8081/elements/ を開くと表示されます。
- `elements/index.html` は相対パスで `../openapi/openapi.yaml` を参照するため、必ずリポジトリのルートでサーバーを起動してください。
- `-c-1` はキャッシュ無効化（編集反映を早めるため、任意）。

OpenAPI仕様を更新した場合の再生成手順は、本READMEの「ビルド手順 > 個別生成」を参照してください。

## コマンド一覧

### ルートディレクトリから実行可能なコマンド

- `pnpm generate:all` - TypeSpec → OpenAPI → 各言語のコード生成を一括実行
- `pnpm compile:spec` - TypeSpecからOpenAPIを生成
- `pnpm generate:typescript` - TypeScriptの型を生成
- `pnpm generate:go` - Goのコードを生成
- `pnpm test` - 全プロジェクトのテストを実行
- `pnpm -F <package> <command>` - 特定のワークスペースでコマンドを実行

### TypeSpec

- `npm run compile` - TypeSpecからOpenAPIを生成
- `npm run watch` - ファイル変更を監視して自動コンパイル
- `npm run format` - TypeSpecファイルのフォーマット

### TypeScript

- `pnpm run generate` - OpenAPIから型定義を生成
- `pnpm run dev` - 開発サーバーを起動 (Port: 3000)
- `pnpm run build` - プロダクションビルド
- `pnpm test` - Vitestによるテスト実行
- `pnpm test:coverage` - カバレッジ付きテスト実行

### Go

- `go generate ./...` - OpenAPIからコードを生成
- `go run cmd/api/main.go` - サーバーを起動 (Port: 8080)
- `go test ./...` - テストを実行
- `go test -cover ./...` - カバレッジ付きテスト実行
- `go fmt ./...` - コードフォーマット
- `go vet ./...` - 静的解析

## 開発フロー

1. TypeSpecでAPI定義を修正
2. `pnpm generate:all`で全コードを再生成（ルートから実行）
3. 生成されたコードを基に実装を追加
4. テストを実行して動作確認
5. CI/CDが自動でビルド・テスト・検証を実行

## テスト

### 全テストの実行

```bash
# ルートから全プロジェクトのテストを実行
pnpm test
```

### TypeScriptテスト

```bash
# 単体テスト
pnpm -F typescript test

# カバレッジ付きテスト
pnpm -F typescript test:coverage
```

### Goテスト

```bash
cd go
# 単体テスト
go test ./...

# カバレッジ付きテスト  
go test -cover ./...
```

## CI/CD

GitHub Actionsによる自動化:

- **PR/Push時**: ビルド、テスト、コード生成の検証、カバレッジ測定
- **週次**: 依存関係のセキュリティチェック

## エラーハンドリング

統一されたエラーレスポンス形式:

```json
{
  "error": {
    "code": "NOT_FOUND",
    "message": "Product not found",
    "details": {} // オプション
  }
}
```

ErrorCodeは`typespec/models/common.tsp`で定義されています。

## 実装状況

- ✅ 全APIエンドポイントの実装完了
- ✅ TypeScript/Go両方のテスト実装（カバレッジ80%以上）
- ✅ エラーハンドリングの統一化
- ✅ CI/CDパイプライン構築
- ✅ ページネーション機能
- ✅ 検索・フィルタリング機能
- ❌ 認証・認可（未実装）
- ❌ 入力バリデーション（基本的なもののみ）

## トラブルシューティング

### コード生成でエラーが出る場合

```bash
# oapi-codegenをインストール
go install github.com/deepmap/oapi-codegen/v2/cmd/oapi-codegen@latest

# パスを通す
export PATH=$PATH:$(go env GOPATH)/bin
```

### TypeScriptのビルドエラー

```bash
# node_modulesを削除して再インストール
rm -rf node_modules pnpm-lock.yaml
pnpm install
```

## ライセンス

MIT
