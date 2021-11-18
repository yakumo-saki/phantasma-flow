# messagehub

MQMTみたいなもの。

## POINT

* map[topic]*[]listeners
* ロックは2段階。 topic毎のlockとmap自体に対するlock。
* 