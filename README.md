# phantasma-flow
Work in progress... (maybe WIP forever)

## What is phantasma flow?

* 小規模向け超簡易ジョブスケジューラ

## Goals

* DBを使わない→バックアップを容易に取得するため
* ログ集積（実行後ログ・実行中ログ）
* ジョブ実行（途中ステップからの再実行含む）
* SSH経由のエージェントレスなジョブ実行 (Windowsは当面対象外。WinRM?）

## Non goals

* High Availability (マルチマスタ）
* ジョブ実行のロバストネス（実行タイミングに落ちていた場合はそのジョブは実行されない）
* プラグイン機構（有用だがむしろ混乱の元になるのでプラグインより本体に取り込むのを優先。取り込めないようなものなら諦める）

