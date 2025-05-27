# Go KVS ベンチマークツール

Go言語における主要なKey-Value Store（KVS）ライブラリの性能を比較するためのベンチマークツールです。

## 対象ライブラリ

以下の5つのKVSライブラリの性能を比較します：

- **LevelDB** - Googleが開発したLSM-treeベースのKVS
- **BBolt** - BoltDBのフォーク、B+treeベースのKVS
- **Badger** - DGraphが開発したLSM-treeベースのKVS
- **Pebble** - CockroachDBが開発したLSM-treeベースのKVS（RocksDBインスパイア）
- **SQLite** - SQLiteをKVSとして使用（インデックス付きテーブル）

## プロジェクト構成

```
├── cmd/
│   ├── generate-data/main.go  # テストデータ生成プログラム
│   └── benchmark/main.go      # ベンチマーク実行プログラム
├── internal/
│   ├── kvs/
│   │   ├── interface.go       # KVS共通インターフェース
│   │   ├── leveldb.go         # LevelDB実装
│   │   ├── bbolt.go           # BBolt実装
│   │   ├── badger.go          # Badger実装
│   │   ├── pebble.go          # Pebble実装
│   │   └── sqlite.go          # SQLite実装
│   └── benchmark/
│       └── runner.go          # ベンチマーク実行ロジック
├── data/
│   └── keys.tsv              # テストデータ（10万件）
├── go.mod
└── README.md
```

## セットアップ

### 依存関係のインストール

```bash
go mod tidy
```

### テストデータの生成

```bash
go run cmd/generate-data/main.go
```

このコマンドにより、`data/keys.tsv`に10万件のテストデータが生成されます。各行は以下の形式です：
- SHA256ハッシュ値（キー）
- 乱数値（値）

## ベンチマーク実行

```bash
go run cmd/benchmark/main.go
```

## ベンチマーク内容

各KVSライブラリに対して以下の4つのベンチマークを実行します：

### 1. Append（追加）
- 10万件のkey-valueペアを新規追加
- value形式: `{"single": 乱数値}`
- 所要時間（ミリ秒）を測定

### 2. Update（更新）
- 同じキーに対してvalueを更新
- value形式: `{"double": 乱数値 * 2}`
- 所要時間（ミリ秒）を測定

### 3. Get（取得）
- 全キーの値を取得し、正しい値が格納されているかを検証
- 期待値: `{"single": 乱数値, "double": 乱数値 * 2}`
- 所要時間（ミリ秒）を測定

### 4. File Size（ファイルサイズ）
- データディレクトリの占有ディスク容量を計測
- 単位: バイト

## パフォーマンス最適化設定

ベンチマークの公平性を保つため、全てのライブラリでfsync（ディスク同期）を無効化しています：

- **LevelDB**: `NoSync: true`
- **BBolt**: `NoSync: true`
- **Badger**: `SyncWrites: false`
- **Pebble**: `pebble.NoSync`
- **SQLite**: `synchronous=OFF&journal_mode=MEMORY`

これにより、ディスクI/Oの影響を最小化し、より純粋なメモリ/CPU性能を測定できます。

## 結果の出力形式

ベンチマーク結果はタブ区切り（TSV）形式で出力されます：

```
Library	Append(ms)	Update(ms)	Get(ms)	FileSize(bytes)
leveldb	1234	567	890	12345678
bbolt	2345	678	901	23456789
badger	3456	789	012	34567890
pebble	4567	890	123	45678901
sqlite	5678	901	234	56789012
```

## 新しいKVSライブラリの追加

新しいKVSライブラリを追加する場合：

1. `internal/kvs/`に新しいファイルを作成
2. `kvs.KVS`インターフェースを実装
3. `cmd/benchmark/main.go`の`kvsLibraries`スライスに追加

### 実装例

```go
package kvs

type NewKVS struct {
    // 必要なフィールド
}

func NewNewKVS() *NewKVS {
    return &NewKVS{}
}

func (n *NewKVS) Name() string {
    return "newkvs"
}

func (n *NewKVS) Open(path string) error {
    // 初期化処理
}

func (n *NewKVS) Close() error {
    // クリーンアップ処理
}

func (n *NewKVS) Set(key string, value *Value) error {
    // キー・バリューの保存
}

func (n *NewKVS) Get(key string) (*Value, error) {
    // キーによる値の取得
}
```

## 注意事項

- ベンチマーク実行時は他のプロセスの影響を最小化することを推奨
- 各ライブラリのデータディレクトリは自動的に削除・再作成されます
- 大量のディスクI/Oが発生するため、SSDでの実行を推奨
- メモリ使用量も考慮に入れてテスト環境を選択してください

## ライセンス

このプロジェクトはMITライセンスの下で公開されています。