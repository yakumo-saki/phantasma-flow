# phantasma-flow

Work in progress.

# 現状

ジョブ実行：ローカル／リモートでジョブが実行可能。 指定時刻／n分ごと。ログ集積とりあえずOK。
ジョブ定義：設定ファイルを手動で記述することで定義可能。 WebUI実装予定
ジョブ結果確認：ログファイルを見るしかない。WebUIに実装予定

## What is phantasma flow?

* 小規模向け超簡易ジョブスケジューラ

## Goals

* DBを使わない→バックアップを容易に取得するため
* ログ集積（実行後ログ・実行中ログ）
* ジョブ実行（途中ステップからの再実行含む）
* SSH経由のエージェントレスなジョブ実行 (Windowsは当面対象外。WinRM?）

## Non goals

* High Availability (マルチマスタ）
* ジョブ実行のロバストネス（実行タイミングにデーモンが落ちていた場合はそのジョブは実行されない）
* プラグイン機構（有用だがむしろ混乱の元になるのでプラグインより本体に取り込むのを優先。取り込めないようなものなら諦める）
* Web UI （このプロジェクトではやらない。別プロジェクトとして作成）

## ラフな設計

* kubernetesライクなオブジェクトの集まり
* オブジェクト類はyamlにシリアライズ
* シャットダウン→すべての状態を失う
* 永続化されたものは起動時に全部読み込む（ログ、実行結果は除く）
* yamlファイルの在処だけは何らかの方法で指定してもらう必要がある
* phctl コマンドを作って通信できるようにする
* phctlの認証はなにかの証明書的なキーで行う（サーバー側ダイジェストに一致するなにか）

## ディレクトリ構造

PHFLOW_HOME
  definitions    `PHFLOW_DEF_DIR`
    config       設定ファイルyaml
    job          ジョブ定義yaml
    node         ノード定義yaml
  data           `PHFLOW_DATA`
    logs         ジョブ実行ログ 
    meta         ジョブ実行結果ログ
  tmp            `PHFLOW_TEMP_DIR` 実行中ログ書き込み os.temp使うべき？

## 関連リポジトリ

関連リポジトリはgolangの標準に則っていない＆モジュールの共有に問題があるので１つのリポジトリに統合予定。

* github.com/yakumosaki/phantasma-flow-cli  ... phantasma-flow CLI
* github.com/yakumosaki/phantasma-flow-web  ... phantasma-flow Web GUI using CLI
