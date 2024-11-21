# fiber-sandbox
fiberが触りテェんだ俺は
このリポジトリはfiberを通してgolangのサーバサイドアプリケーションをいじるための遊び場です

## tips

- `prefork` を使う
  - OSレベルのport shardingがきいて複数の`fiber`プロセスが単一のソケットをリッスンできるようになる
  - 使う場合は`sh`で実行が必要（謎）（[根拠](https://github.com/gofiber/fiber/issues/1021#issuecomment-730537971)）
- `compress` を使う
  - アプリのリソース < ネットワークi/o の場合に有効になる（妄想でものを言っている）
- `ETag` を使う
  - 何それ
- `cache` を使う
  - どういう条件で何がどこにキャッシュされる？
- `DisableHeaderNormalizing: true`
  - 何が？