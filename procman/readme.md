# procman

* channelのメッセージは誰か一人が受信すると終わりなのでメッセージを分配する
* 

## 機能

* シャットダウン時のメッセージ伝達
* シャットダウン完了まで待つ(shutdown)

## スタートアップシーケンス

* procman.Start() -> Added modules
* modules -> "STARTUP_DONE" -> channel -> procman
* NOTE: start all services -> wait for started up -> start all modules -> wait

## シャットダウンシーケンス

* procman.Shutdown() -> Added modules
* modules -> "SHUTDOWN_DONE" -> channel -> procman
* All modules shutdown or timeout, then return REASON("DONE" / "TIMEOUT")