# プロジェクト固有の指示 - TypeSpec ECサイトAPI

このファイルはClaude Codeのメモリ機能用の指示書です。プロジェクトの技術的な決定事項と開発方針を記載します。

## プロジェクト概要

- **目的**: TypeSpecを使用したAPI定義からTypeScriptとGoの実装を生成するサンプルプロジェクト
- **ドメイン**: ECサイト（Eコマース）
- **認証**: なし（サンプルのため）
- **デプロイ**: なし（ローカル開発のみ）

## 技術的決定事項

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

## 開発時の注意事項

- 生成されたコードは直接編集しない
- TypeSpec定義を変更したら必ず再生成する
- コミット時は生成コードも含める
- エラーメッセージは日本語ではなく英語で記述
- コミットは適切な粒度で行う
