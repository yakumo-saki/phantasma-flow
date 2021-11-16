# procman

* channelのメッセージは誰か一人が受信すると終わりなのでメッセージを分配する
* 

## 機能

* シャットダウン時のメッセージ伝達
* シャットダウン完了まで待つ(shutdown)

## シャットダウンシーケンス

* procman.Shutdown() -> "SHUTDOWN" -> Subscribed modules
* Subscribed module -> "SHUTDOWN_DONE" -> procman
* All modules shutdown or timeout, then return REASON("DONE" / "TIMEOUT")