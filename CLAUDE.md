# プロジェクト固有の指示 - TypeSpec ECサイトAPI

このファイルはClaude Codeのメモリ機能用の指示書です。プロジェクトの技術的な決定事項と開発方針を記載します。

## プロジェクト概要

- **目的**: TypeSpecを使用したAPI定義からTypeScriptとGoの実装を生成するサンプルプロジェクト
- **ドメイン**: ECサイト（Eコマース）
- **認証**: なし（サンプルのため）
- **デプロイ**: なし（ローカル開発のみ）

## 技術的決定事項

### プロジェクト構造
- **pnpmワークスペース**を使用したモノレポ構造
- 3つの主要なpackage.json:
  - `/package.json` - ルート（ワークスペース管理）
  - `/typescript/package.json` - TypeScriptサブプロジェクト
  - `/typespec/package.json` - TypeSpecサブプロジェクト
- コマンド実行時は**ルートから `pnpm --filter` または `pnpm -F`** を使うか、**サブディレクトリで直接実行**

### API定義フロー
1. TypeSpecでAPI定義を記述
2. TypeSpecからOpenAPI 3.1を生成
3. OpenAPIから各言語のコードを生成

### 言語別の実装方針

#### TypeScript
- **型生成**: openapi-typescriptを使用（型定義のみ）
- **サーバー実装**: Honoを使用して手動実装
- **バリデーション**: 必要に応じてZodを使用
- **パッケージマネージャー**: pnpm

#### Go
- **コード生成**: oapi-codegenを使用
- **HTTPライブラリ**: Go標準パッケージ（net/http）
- **ルーティング**: chi（oapi-codegenのデフォルト）は使用せず、標準パッケージで実装

## コーディング規約

### TypeSpec
- モデルは`/models`ディレクトリに分割
- サービスは`/services`ディレクトリに分割
- 共通型は`common.tsp`に定義
- OpenAPIの`operationId`は`{リソース名}{動詞}`形式（例: `productsCreate`）

### TypeScript
- strictモードを有効化
- ESModuleを使用
- エラーハンドリングは明示的に行う
- ビジネスロジックは`/services`に分離
- **注意**: TypeScriptのpackage.jsonにはlint/typecheckスクリプトがない（必要に応じてvite buildで代用）

### Go
- 標準的なGoプロジェクトレイアウトに従う
- エラーハンドリングは標準的なGoの方式を使用
- contextは適切に伝播させる

## API設計方針

### RESTful設計
- リソース指向のURL設計
- 適切なHTTPメソッドの使用
- ステータスコードは標準に従う

### レスポンス形式
- 成功時: リソースをそのまま返す
- エラー時: `{error: {code: string, message: string}}`形式

### ページネーション
- クエリパラメータ: `limit`と`offset`
- レスポンスに`total`を含める

## ECサイトのデータモデル

### Products（商品）
- id, name, description, price, stock, categoryId, imageUrls
- 検索、フィルタリング機能を提供

### Categories（カテゴリ）
- id, name, parentId（階層構造）
- ネストされたカテゴリ取得をサポート

### Cart（カート）
- userId, items（productId, quantity）
- セッションベースまたはユーザーベース

### Orders（注文）
- id, userId, items, totalAmount, status, createdAt
- ステータス: pending, processing, shipped, delivered, cancelled

### Users（ユーザー）
- id, email, name, address
- 認証は実装しない（サンプルのため）

## 今後の拡張計画

1. WebSocket対応（リアルタイム在庫更新）
2. GraphQL対応（TypeSpecのカスタムエミッター）
3. gRPC対応
4. 認証・認可の追加（JWT）
5. テストコードの自動生成

## 現在の実装状況（2025/06/14時点）

### 完了済み
- ✅ 全APIエンドポイントの基本実装（TypeScript/Go）
- ✅ TypeScriptのテスト実装（50テスト）
- ✅ メモリベースのデータストア
- ✅ ページネーション（Products, Orders）
- ✅ 検索・フィルタリング機能
- ✅ エラーハンドリング

### 未実装・課題
- ❌ Goのテスト（テストファイルが存在しない）
- ❌ CI/CDパイプライン
- ❌ TypeSpec高度機能（バリデーション、@example等）
- ❌ 認証・認可
- ✅ ~~TypeScript: 注文作成時のshippingAddress処理にバグ~~ **修正済み**
- ✅ ~~TypeScript: UsersのページネーションがAPI定義と不一致~~ **修正済み**

## 改善提案と作業方針

### 重要な作業方針
**各改善作業を開始する前に、必ずユーザーに実施の判断を仰ぐこと。**

### 高優先度の改善案

#### 1. バグ修正
- **TypeScript注文作成バグ**: shippingAddressの処理を修正
- **作業前の確認事項**: 仕様通り（毎回住所指定）でよいか、ユーザーのデフォルトアドレスを使用すべきか

#### 2. テスト実装
- **Goテスト**: 全APIエンドポイントのテスト実装（50以上のテストケース）
- **作業前の確認事項**: TypeSpecサンプルとしてGoのテストは必須か

#### 3. CI/CD整備
- **GitHub Actions**: コード生成検証、テスト実行、lint
- **作業前の確認事項**: サンプルプロジェクトにCI/CDは必要か

### 中優先度の改善案

#### 4. TypeSpec機能の活用
- **PATCH警告解消**: `@patch(#{implicitOptionality: true})`の追加
- **バリデーション**: `@minLength`, `@maxValue`, `@format`等の追加
- **認証定義**: `@useAuth`によるBearerトークン認証
- **作業前の確認事項**: TypeSpecの高度な機能をどこまで含めるか

#### 5. 実装の一貫性
- **Usersページネーション**: TypeScriptでの実装追加
- **作業前の確認事項**: API定義との一貫性を優先するか

#### 6. ドキュメント充実
- **API利用ガイド**: 各エンドポイントの詳細な使用例
- **作業前の確認事項**: どの程度詳細なドキュメントが必要か

### 低優先度の改善案

#### 7. 開発体験向上
- **Watchモード**: TypeSpecの自動再生成
- **Pre-commitフック**: 生成忘れ防止
- **作業前の確認事項**: 開発ツールをサンプルに含めるか

#### 8. TypeSpec例示機能
- **@example**: 各モデルにサンプルデータ定義
- **作業前の確認事項**: OpenAPIドキュメントの充実は必要か

## TypeSpecで使用可能な高度な機能（未使用）

### バリデーション
```typescript
@minLength(3) @maxLength(100) name: string;
@minValue(0.01) @maxValue(999999.99) price: float64;
@format("email") email: string;
```

### セキュリティ
```typescript
@useAuth(BearerAuth)
@secret apiKey: string;
```

### 可視性制御
```typescript
@visibility("create", "update") password: string;
@visibility("read") id: uuid;
```

### その他
- `@deprecated`: 非推奨マーキング
- `@discriminator`: ユニオン型の判別
- `@versioned`: APIバージョニング
- `@knownValues`: 拡張可能なenum

## 開発時の注意事項

- 生成されたコードは直接編集しない
- TypeSpec定義を変更したら必ず再生成する
- コミット時は生成コードも含める
- エラーメッセージは日本語ではなく英語で記述
- コミットは適切な粒度で行う
- **新機能追加前に必ずユーザーに確認を取る**

## TypeSpec定義変更時の影響範囲

### 変更時の再生成手順
1. `pnpm compile:spec` - TypeSpecからOpenAPIを生成
2. `pnpm generate:typescript` - TypeScript型定義を生成
3. `pnpm generate:go` - Goコードを生成
4. または`pnpm generate:all`で一括実行

### 影響を受けるファイル
- **OpenAPI定義**: `/openapi/openapi.yaml`
- **TypeScript**: `/typescript/src/types/api.ts`（自動生成）
- **Go**: `/go/generated/server.gen.go`（自動生成）
- **実装コード**: APIの型定義が変わった場合は手動修正が必要

### TypeSpecのベストプラクティス
- **成功レスポンスとエラーレスポンスは分離する**
  - ❌ `list(): PaginatedResponse<T> | ErrorResponse`
  - ✅ `list(): PaginatedResponse<T>` （エラーは適切なHTTPステータスで返す）
- **HTTPステータスコードは@statusDecoratorで明示的に指定**
- **エラーレスポンスは4xx/5xxステータスで返す**
