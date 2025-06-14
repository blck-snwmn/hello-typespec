# プロジェクト固有の指示 - TypeSpec ECサイトAPI

このファイルはClaude Codeのメモリ機能用の指示書です。プロジェクトの技術的な決定事項と開発方針を記載します。

## プロジェクト概要

- **目的**: TypeSpecを使用したAPI定義からTypeScriptとGoの実装を生成するサンプルプロジェクト
- **ドメイン**: ECサイト（Eコマース）
- **認証**: JWT認証を実装予定
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
- ✅ TypeScriptのテスト実装（51テスト）
- ✅ Goのテスト実装（50以上のテストケース）
- ✅ メモリベースのデータストア
- ✅ ページネーション（Products, Orders, Users）
- ✅ 検索・フィルタリング機能
- ✅ CI/CDパイプライン（GitHub Actions）
  - PR時のチェック（ビルド、テスト、コード生成検証）
  - Push時のチェック（カバレッジ測定含む）
  - 依存関係のセキュリティチェック（週次）
- ✅ TypeScript: Viteのライブラリビルド設定
- ✅ テストカバレッジ測定環境（@vitest/coverage-v8）
- ✅ **エラーハンドリングの統一化**（2025/06/14）
  - TypeSpec: ErrorCode enumの定義
  - TypeScript: 集約型エラーハンドリング実装
  - TypeScript: グローバルエラーハンドラー追加
  - Go: エラーヘルパー関数の作成（部分的）

### 未実装・課題
- ❌ TypeSpec高度機能（バリデーション、@example等）
- ❌ 認証・認可
- ❌ 入力バリデーション

## 今後の改善計画

### 重要な作業方針
**各改善作業を開始する前に、必ずユーザーに実施の判断を仰ぐこと。**

### 高優先度の改善項目

#### 1. 入力バリデーション（未着手）
- TypeSpec: バリデーションデコレータ（`@minLength`, `@maxLength`, `@minValue`, `@format`等）
- TypeScript: Zodスキーマによる検証
- Go: カスタムバリデーション関数

#### 2. 認証・認可（未着手）
- TypeSpec: 認証定義の追加（`@useAuth`, `BearerAuth`）
- TypeScript: JWT認証ミドルウェアの実装
- Go: JWT認証ミドルウェアの実装
- 保護が必要なエンドポイントの特定と実装

### 中優先度の改善項目

#### 3. Goエラーハンドリングの完全統一
- products.goの重複errorResponse関数の削除
- 全ハンドラーでhelpers.goの関数を使用

#### 4. テストの充実
- エンドツーエンドテストの追加
- パフォーマンステストの実施
- セキュリティテストの追加

### 低優先度の改善項目

#### 5. 開発体験向上
- TypeSpec watchモードの設定
- Pre-commitフックの追加

#### 6. ドキュメント充実
- API利用ガイドの作成
- `@example`デコレータによるサンプルデータ定義

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

### 基本的な注意事項
- 生成されたコードは直接編集しない
- TypeSpec定義を変更したら必ず再生成する
- コミット時は生成コードも含める
- エラーメッセージは日本語ではなく英語で記述
- コミットは適切な粒度で行う
- **新機能追加前に必ずユーザーに確認を取る**

### Git操作のベストプラクティス

#### コミット前の確認事項
1. **Diagnosticsエラーの確認**
   - VSCodeのDiagnosticsパネルでエラーや警告がないことを確認
   - 特にTypeScriptの型エラー、Goのコンパイルエラーに注意

2. **テストの実行**
   - TypeScript: `pnpm test` (vitestの実行)
   - Go: `go test ./...`
   - すべてのテストがパスすることを確認

3. **ビルドの確認**
   - TypeScript: `pnpm build` または `pnpm -F typescript build`
   - Go: `go build ./...`
   - ビルドエラーがないことを確認

4. **差分の確認**
   - `git diff` で未ステージの変更を確認
   - `git diff --cached` でステージ済みの変更を確認
   - 意図しない変更が含まれていないことを確認

#### git addの前に行うこと
- **必ず `git status` で現在の状態を確認**
- **変更内容を `git diff` で詳細に確認**
- カバレッジファイルなど、コミット不要なファイルは`.gitignore`に追加
- 大きな変更は複数のコミットに分割することを検討

#### コミットメッセージの規約
- feat: 新機能の追加
- fix: バグ修正
- refactor: リファクタリング
- test: テストの追加・修正
- docs: ドキュメントの更新
- chore: ビルドプロセスや補助ツールの変更

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

## CI/CD設定

### GitHub Actionsワークフロー

#### 1. CI Check (ci.yml)
プルリクエストおよびプッシュ時に実行される検証：
- **TypeSpec検証**: コード生成が最新かチェック
- **TypeScriptチェック**: ビルド、テスト、型チェック、カバレッジ測定
- **Goチェック**: ビルド、テスト、フォーマット、静的解析、カバレッジ測定
- **カバレッジレポート**: アーティファクトとして保存

#### 2. Dependency Check (dependency-check.yml)
週次で実行されるセキュリティチェック：
- **npm/pnpm監査**: 高レベルの脆弱性をチェック
- **Go脆弱性チェック**: govulncheckによる検査
- **依存関係の更新確認**: 古いパッケージの確認

### CI/CDの注意事項
- **oapi-codegen**はCIで自動インストールされる
- **@vitest/coverage-v8**はテストカバレッジ測定に必要
- **typescript/dist/**は.gitignoreに追加してコミットしない
- **typescript/coverage/**は.gitignoreに追加してコミットしない

## エラーハンドリング実装の詳細

### TypeSpec定義 (`typespec/models/common.tsp`)
- ErrorCode enumでエラーコードを標準化
- クライアントエラー: BAD_REQUEST, NOT_FOUND, VALIDATION_ERROR等
- サーバーエラー: INTERNAL_ERROR, SERVICE_UNAVAILABLE

### TypeScript実装 (`typescript/src/types/errors.ts`)
- `sendError`関数: 統一的なエラーレスポンス送信
- `globalErrorHandler`: Honoのグローバルエラーハンドラー
- エラーコードとHTTPステータスのマッピング

### Go実装 (`go/internal/handlers/helpers.go`)
- `errorResponse`関数: 標準エラーレスポンス
- `errorResponseWithDetails`関数: 詳細付きエラーレスポンス
- ※注意: products.goに古いerrorResponse関数が残存（要リファクタリング）

## 次のセッションでの作業候補

1. **入力バリデーションの実装**（高優先度）
   - TypeSpecにバリデーションデコレータを追加
   - TypeScript/Goでバリデーションロジックを実装

2. **Goエラーハンドリングの完全統一**（中優先度）
   - products.goの重複関数削除
   - 全ハンドラーでhelpers.goの関数を使用

3. **認証・認可の実装**（高優先度）
   - JWT認証の基本実装
   - エンドポイントの保護
