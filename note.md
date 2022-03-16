# Chapter 3: Disk and File Management
## Direct IO
普通のファイル I/O では disk と memory 間の情報の読み書きでは、間に page cache というものが入る。
page cache は memory 上にあり、ファイル I/O を高速化するために導入された。

disk の情報を memory に読み出すことを考える。
アプリケーションが OS に対して disk の読み出し要求をすると、
page cache が空だった場合、disk -> page cache -> process が管理する memory と情報が流れる。
page cache にすでに読み込みたい disk の内容があった場合は、OS は disk 読み出しをせず、page cache にある情報をそのまま返す。
memory の内容を disk に書き込む際は page cache の内容だけが更新され、すぐには disk への書き込みは行われない。
page cache の領域が満杯になったり、明示的に disk への書き込みが要求されたり、OS の決めたよしななタイミングで page cache の内容が disk に書き込まれデータの永続化がされる。
このように file IO を減らすことによってパフォーマンスの向上を図っている。

頻繁に使われるファイルは page cache に乗せ、あまり使われないものは page cache に乗せないという戦略を取るのが、OS 側の機能を使うとそのへんをコントロールすることができない。
例えば、小さいファイルだが頻繁にアクセスするものが page cache に乗っている状態で、アクセス頻度の少ない大きなファイルを読み込むと page cache の内容に取って代わってしまって全体的な効率が悪くなってしまう。

DBMS の場合、クエリを解析することで、どのファイルがどれくらいの頻度で使われるかということを OS よりも知っているため、
page cache の更新はファイルの書き出しは DBMS 自身が管理したくなる。
このようなときに Direct IO というのを使うと page cache を経由せずに
disk と process の memory 間で直接データのやり取りができるようになる。
DBMS では buffer pool というものを用いて、memory <-> disk の IO をコントロールしている。

CMU の講義で mmap を使うなと言われているが、それは mmap が上で書いたように OS 側の制御下にあるからだと思う。(mmap が page cache 使っているのかよくわからんが)
https://db.cs.cmu.edu/mmap-cidr2022/

Go の場合、C 言語で使えるような DIRECT_IO はどうやらサポートされていないらしいので [ncw/directio](https://github.com/ncw/directio) を使うといいっぽい。

ダイレクトI/O
- ダイレクトI/O
  - https://xtech.nikkei.com/it/article/Keyword/20070207/261244/
- その23 同期I/Oとdirect I/O
  - https://youtu.be/sn6EKG0_lOU
- Linux ファイルシステム 徹底入門
  - https://www.kimullaa.com/entry/2019/12/01/130347#direct-IO
- Goでdirect I/O
  - https://satoru-takeuchi.hatenablog.com/entry/2020/03/26/011423

## Synchonize

SimpleDB の file manager は以下のように read, write, append に synchronized がついている。

```java
public class FileMgr;
    public synchronized void read(BlockId blk, Page p);
    public synchronized void write(BlockId blk, Page p);
    public synchronized BlockId append(String filename);
```

Java に詳しくないのでこれはてっきり複数スレッドで read を呼べるのが一つだけであり、read を呼んでいるときに write も同時に呼んでいると思っていたがリンクの IPA の記述によるとそうではなく、インスタンスをロックしているらしい。

> あるスレッドがsynchronized指定されたメソッドを実行している間，実はそのメソッドの持ち主のオブジェクト全体がロックされている。同じオブジェクトの他のsynchronizedメソッドの呼び出しについてもブロックされてしまうのだ。本来は並列で動作しても差し支えない複数のメソッドがあっても， synchronizedが指定されていると同時にはただ一つのスレッドしか動作できないことになる。
https://www.ipa.go.jp/security/awareness/vendor/programmingv1/a03_06.html

下の2つの処理が等価らしいのでこれを見ると、`synchronized(this)` でインスタンスへのアクセスを1つのスレッドに限定しているから、IPA の記述も納得できる。
```
public synchronized T func {
    do something
}

public T func {
    synchronized(this){
        do something
    }
}
```
- https://www.techscore.com/tech/Java/JavaSE/Thread/3/
- https://www.techscore.com/tech/Java/JavaSE/Thread/3-2/

Go で同じようなことをしようと思ったら、struct に mutex をもたせる感じになるかな？


それはそれとして、read, write, append が一つのスレッドでしかできないという
制限はパフォーマンス的に大丈夫なのだろうか？
複数ユーザーが使う DBMS ならば同時に複数のファイルを disk から読み込みたいとか普通にありそうだけれども。
DBMS が管理している page cache は一つだから同時に読み書きすると page cache の
奪い合いが起こるとかそんなところだろうか？


# Chapter 4: Memory Management
4.3.1 p.84 には最初の `printLogRecords` で20行だけ表示されると書いてあるが、
iterator を呼ばれた段階で `flush` を最初にやっているので実際は page に書かれた 21 ~ 35 の
record も表示される。

```java
public Iterator<byte[]> iterator() {
  flush();
  return new LogIterator(fm, currentblk);
}
```
