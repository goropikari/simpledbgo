SimpleDB に関するコードは SimpleDB 3.4 を元にしている。

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


synchronized の中で wait を呼ぶと、スレッドそのオブジェクトの lock を開放する。そのため、他のスレッドが他の synchronized method を使うことができるようになる。

https://www.techscore.com/tech/Java/JavaSE/Thread/5-2/

Go で Java の wait, notify, notifyAll と同じことをしようとしたら、sync.Cond.Wait, Signal, Broadcast を使うのが良さそう。
ただ、Java と違って timeout がない。timeout がないせいで運の悪い goroutine がいつまでも残りそうな気がするが一旦無視する。


2022/4/6

goroutine の timeout は channel と select を使うのが Go ではよくあるので試しに実装してみたがだいぶ処理が複雑になってしまった。
https://github.com/goropikari/simpledbgo/blame/972526679d15cf5eb5d6d10a78b5192767714d38/backend/buffer/manager.go

```go
func (mgr *Manager) Pin(block *domain.Block) (*domain.Buffer, error) {
	mgr.mu.Lock()

	buf, err := mgr.tryToPin(block, naiveSearchUnpinnedBuffer)
	if err != nil {
		mgr.mu.Unlock()

		return nil, err
	}

	if buf == nil {
		mgr.mu.Unlock()
		select {
		case <-mgr.ch:
			mgr.mu.Lock()
			mgr.ch <- item

			buf, err = mgr.tryToPin(block, naiveSearchUnpinnedBuffer)
			if err != nil {
				mgr.mu.Unlock()

				return nil, err
			}
		case <-time.After(mgr.timeout):
			return nil, ErrTimeoutExceeded
		}
	}

	mgr.mu.Unlock()

	return buf, nil
}

// tryToPin tries to pin the block to a buffer.
func (mgr *Manager) tryToPin(block *domain.Block, chooseUnpinnedBuffer func([]*domain.Buffer) *domain.Buffer) (*domain.Buffer, error) {
	buf := mgr.findExistingBuffer(block)
	if buf == nil {
		buf = chooseUnpinnedBuffer(mgr.bufferPool)
		if buf == nil {
			return buf, nil
		}
		if err := buf.AssignToBlock(block); err != nil {
			return nil, err
		}
	}

	if !buf.IsPinned() {
		mgr.numAvailableBuffer--
		<-mgr.ch
	}

	buf.Pin()

	return buf, nil
}

// Unpin unpins buffer.
func (mgr *Manager) Unpin(buf *domain.Buffer) {
	mgr.mu.Lock()
	defer mgr.mu.Unlock()

	buf.Unpin()
	if !buf.IsPinned() {
		mgr.numAvailableBuffer++
		mgr.ch <- item
	}
}
```

```
case <-mgr.ch:
	mgr.mu.Lock()
	mgr.ch <- item
```

の部分で条件を満たした goroutine が Lock をとって処理を進めるという風に考えていたが、channel から要素を引っ張ってくる
goroutine と Lock を取る goroutine が同じになるとは限らないのでこの方法ではだめだと気づいた。

次のような場合を考える
- buffer pool size: 1
- goroutine 1 が block 1 を Pin している
- goroutine 2 が block 2 を Pin しようとするが空き buffer がないので select 文で待つ
- goroutine 1 が Unpin するタイミングとほぼ同時に goroutine 3 が Pin をする。

このような状況のとき、まず goroutine 2 は channel から要素を取れるので次に進める。
問題はここでの Lock を goroutine 2 と 3 でどちらが取るかということである。
goroutine 2 が Lock を取れれば期待した動きになるが、goroutine 3 が Lock を取ると問題がおきる。

goroutine 3 が Lock をとった時点で空き buffer があるので goroutine 3 は buffer に新たな block を割り当てることができる。
goroutine 3 が Unlock したあと goroutine 2 の select 内の処理が走るが、この中の処理では確実に buffer の割当が
できることを期待しているが実際は goroutine 3 による block 割当がすでになされているので goroutine 2 は割当をすることができない。
そのため select 内では Lock をとったあとにまた条件を満たしているかの判定をしないといけない。
channel を使った場合だと Lock を取るタイミングや、channel の要素の出し入れ等についてかなり考えないといけないので
最終的に素直に sync.Cond を使うのが一番最適だという結論に至った。

```go
// Pin pins buffer.
func (mgr *Manager) Pin(block *domain.Block) (*domain.Buffer, error) {
	mgr.mu.Lock()
	defer mgr.mu.Unlock()

	buf, err := mgr.tryToPin(block, naiveSearchUnpinnedBuffer)
	if err != nil {
		return nil, err
	}

	now := time.Now()

	for buf == nil {
		mgr.cond.Wait()
		if time.Since(now) > mgr.timeout {
			return nil, ErrTimeoutExceeded
		}
		buf, err = mgr.tryToPin(block, naiveSearchUnpinnedBuffer)
		if err != nil {
			return nil, err
		}
	}

	return buf, nil
}
```

https://github.com/goropikari/simpledbgo/blob/b1ba3f41fbc782214829c25b57e23e376f2cf052/backend/buffer/manager.go


上の方法で良いかと思ったが、`Wait()` が timeout しないせいで deadlock を起こす可能性に気づいた。
concurrency manager がうまく捌いてくれそうな気もするがやはり timeout が欲しくなったので goroutine 使った実装を改めて書いた

```
type result struct {
	buf *domain.Buffer
	err error
}

// Pin pins buffer.
func (mgr *Manager) Pin(block *domain.Block) (*domain.Buffer, error) {
	done := make(chan *result)

	go mgr.pin(done, block)
	select {
	case result := <-done:
		return result.buf, result.err
	case <-time.After(mgr.timeout):
		return nil, ErrTimeoutExceeded
	}
}

func (mgr *Manager) pin(done chan *result, block *domain.Block) {
	mgr.mu.Lock()
	defer mgr.mu.Unlock()
	defer close(done)

	buf, err := mgr.tryToPin(block, naiveSearchUnpinnedBuffer)
	if err != nil {
		done <- &result{err: err}
		return
	}

	for buf == nil {
		mgr.cond.Wait()
		buf, err = mgr.tryToPin(block, naiveSearchUnpinnedBuffer)
		if err != nil {
			done <- &result{err: err}
			return
		}
	}

	done <- &result{buf: buf}
}
```

# Chapter 5: Transaction Management

元の Java の実装だと RecoveryMgr class が Transaction class に依存し、また逆に
Transaction class が RecoveryMgr class に依存しているのがなんだか気持ち悪い.

```java
public class Transaction {
    private static int nextTxNum = 0;
    private static final int END_OF_FILE = -1;
    private RecoveryMgr    recoveryMgr;
    private ConcurrencyMgr concurMgr;
    private BufferMgr bm;
    private FileMgr fm;
    private int txnum;
    private BufferList mybuffers;

    public void rollback() {
        recoveryMgr.rollback();
        System.out.println("transaction " + txnum + " rolled back");
        concurMgr.release();
        mybuffers.unpinAll();
    }
    ...
}


public class RecoveryMgr {
    private LogMgr lm;
    private BufferMgr bm;
    private Transaction tx;
    private int txnum;

    private void doRollback() {
        Iterator<byte[]> iter = lm.iterator();
        while (iter.hasNext()) {
            byte[] bytes = iter.next();
            LogRecord rec = LogRecord.createLogRecord(bytes);
            if (rec.txNumber() == txnum) {
                if (rec.op() == START)
                    return;
                rec.undo(tx);
            }
        }
    }
    ...
}
```

下のように visitor pattern を使うといくぶんか違和感はなくなった気がする。

```go
package recovery

struct Manager {
    logMgr log.Manager
    bufMgr buffer.Manager
    txnum  int
}

func (mgr *Manager) rollback(tx TransactionVisitor) {
    mgr.doRollback(tx)
}

func (mgr *Manager) doRollback(tx TransactionVisitor) {
    for _, bytes := range mgr.logMgr.Iterator() {
        rec := factory(rec) // rec is LogRecord
        ...
            rec.undoAccept(tx)
    }
}

type LogRecord interface {
    undoAccept(TransactionVisitor)
}

type TransactionVisitor interface {
    VisitSetIntRecord(SetIntRecord)
    VisitSetStringRecord(SetStringRecord)
    VisitStartRecord(StartRecord)
    VisitRollbackRecord(RollbackRecord)
}

type SetIntRecord struct {}

func (rec *SetIntRecord) undoAccept(tx TransactionVisitor) {
    tx.visitSetIntRecord(rec)
}
...
```

```go
package transaction

type Transaction struct {}

func (tx *Transaction) VisitSetIntRecord(rec recovery.SetIntRecord) {
    tx.pin(rec.block)
    tx.setInt(rec.block, rec.offset, rec.val, false)
    tx.unpin(rec.block)
}
```

## 2022/4/10
### recovery manager は必要ないのでは?
本で recovery manager までを読んでいたときは visitor pattern でいいかと思っていたが、
chapter 5 を全て読んで concurrency manager も含めてどうするのが理想形をなのか再考すると
そもそも recovery manager なんてものは作らないので、Transaction の中に全て押し込めれば
よいのではないかという気がした。
grep する限り recovery manager は transaction でしか使われていないのでそれでも問題なさそう。
また recovery manager の methods にはどこにも synchronized ついてないからそういう意味でも
transaction の中に持っていって問題なさそう。


### `sync.RWMutex` の落とし穴

Go の `sync.RWMutex` は名前から想像できるように Read に関する lock は複数取れる、read が lock しているうちは write は lock を取れない。
逆に write が lock をとっているときは read も write をさらに lock はできない。
ここで次のような処理を考える
- R1: 時間 0 で read lock を取る。100 msec 後に unlock
- W: 時間 10 msec で write lock を取ろうとする。10 msec 後に unlock
- R2: 時間 20 msec で read lock を取ろうとする。10 msec 後に unlock

上のような処理を流したとき read lock はいくつでも取れるから R1, R2, W の順で lock を取って、
トータルの実行時間としては 110 msec になると思っていた。
だが実際に試してみると 120 msec になった。

```go
package main

import (
	"fmt"
	"sync"
	"time"
)

func main() {
	mu := &sync.RWMutex{}
	wg := &sync.WaitGroup{}
	wg.Add(3)

	now := time.Now()

	go func() {
		mu.RLock()
		time.Sleep(100 * time.Millisecond)
		mu.RUnlock()
		wg.Done()
	}()

	go func() {
		time.Sleep(10 * time.Millisecond)
		mu.Lock()
		time.Sleep(10 * time.Millisecond)
		mu.Unlock()
		wg.Done()
	}()

	go func() {
		time.Sleep(20 * time.Millisecond)
		mu.RLock()
		time.Sleep(10 * time.Millisecond)
		mu.RUnlock()
		wg.Done()
	}()

	wg.Wait()

	fmt.Println(time.Now().Sub(now))
}
```

自分の実装方法が悪いのだと思っていたが、read lock 取ったあとに write lock を
取ろうとする goroutine がいる場合、他の goroutine が lock を取ろうとしても
最初の read lock が release されるまで追加で read lock  はとれないと公式 doc に書いてあった。

> If a goroutine holds a RWMutex for reading and another goroutine might call Lock, no goroutine should expect to be able to acquire a read lock until the initial read lock is released. In particular, this prohibits recursive read locking. This is to ensure that the lock eventually becomes available; a blocked Lock call excludes new readers from acquiring the lock.

https://pkg.go.dev/sync#RWMutex

write lock が取れないのはわかるが read lock まで取れないとは思っていなかった。
`a blocked Lock call excludes new readers from acquiring the lock.` とあるのでどうやら
writer lock が優先されるようである。
reader lock を優先するか、writer lock を優先するかという問題は結構メジャーなものだったらしい。
wikipedia にそれぞれのメリット・デメリットが書いてあった。

https://en.wikipedia.org/wiki/Readers%E2%80%93writer_lock#Priority_policies

この辺の優先度がなかったら運の悪い writer がなかなか lock を取れないだろうなぁと思っていたが
writer 優先ならばその心配はなさそうである。
write の頻度が少なくて read の頻度が高いならば writer 優先にしても並列度への影響は
少ないということでなるほどなぁと思った。

## 2022/7/18
RWMutex を使った実装はどうにも間違った実装でバグっていた。
もともとテキストの実装だと exclusive lock は slock を取ったあとに xlock を呼ぶ(xlock を呼ぶ前に unlock は呼ばない)という方針になっていたところを無理やり LockTable だけ RWMutex を使う実装にしていたので全体的な整合性が崩れてしまうようになっていた。
そのせいで異なる transaction が同時に同じ領域を書き込める状態になってしまっていた。そのためテキスト通りに LockTable を作ることにしたがどうにもテキストの実装は実装でバグっている気がする。


```java
class LockTable {
    ...

    public synchronized void sLock(BlockId blk) {
        try {
            long timestamp = System.currentTimeMillis();
            while (hasXlock(blk) && !waitingTooLong(timestamp))
                wait(MAX_TIME);
            if (hasXlock(blk))
                throw new LockAbortException();
            int val = getLockVal(blk);  // will not be negative
            locks.put(blk, val+1);
        }
        catch(InterruptedException e) {
            throw new LockAbortException();
        }
    }

    synchronized void xLock(BlockId blk) {
        try {
            long timestamp = System.currentTimeMillis();
            while (hasOtherSLocks(blk) && !waitingTooLong(timestamp))
                wait(MAX_TIME);
            if (hasOtherSLocks(blk))
                throw new LockAbortException();
            locks.put(blk, -1);
        }
        catch(InterruptedException e) {
            throw new LockAbortException();
        }
    }

    synchronized void unlock(BlockId blk) {
        int val = getLockVal(blk);
        if (val > 1)
            locks.put(blk, val-1);
        else {
            locks.remove(blk);
            notifyAll();
        }
    }

    private boolean hasXlock(BlockId blk) {
        return getLockVal(blk) < 0;
    }

    private boolean hasOtherSLocks(BlockId blk) {
        return getLockVal(blk) > 1;
    }
}

public class ConcurrencyMgr {
    ...

    public void sLock(BlockId blk) {
        if (locks.get(blk) == null) {
            locktbl.sLock(blk);
            locks.put(blk, "S");
        }
    }

    public void xLock(BlockId blk) {
        if (!hasXLock(blk)) {
            sLock(blk);
            locktbl.xLock(blk);
            locks.put(blk, "X");
        }
    }
}
```


気になっているのは unlock の中で呼ばれる notifyAll のタイミングで、xlock 取っている状態(val = -1) のときはここに入れるだけで次に来る slock の待ちを正しく解消することができる。(exclusine lock のときでも最初は slock を取っているので現時点で xlock が取られているときの待っている lock は常に slock になる。)
問題は現時点で slock が取られている状態のときに exclusive lock を取ろうとした場合である。
この場合、(`slock`), (`slock`, `xlock`) の順に method が呼ばれるが、slock は重ねて取ることができるので最初の transaction が unlock を呼ぶ前までは `locks[blk]` の値は 2 となる。そしてこの状態で xlock は wait で待ち続ける。
ここで最初の transaction が unlock すると `locks.put(blk, val-1);` は呼ばれるが `notifyAll` は呼ばれないので xlock の wait は timeout になるまで待ち続けることになってしまう。
なので `val = 2` のときは decrement に加えて notifyAll も加えないといけないように思える。

```java
if (val > 1)
    locks.put(blk, val-1);
else {
    locks.remove(blk);
    notifyAll();
}
```


### 謎の newval

RecoveryMgr の `setInt`, `setString` で引数に newval というのがあるけど、実際の処理を見ると
全く使っていない。
`writeToLog` に渡している `oldval` を `newval` にするべきなのかと思ったがそういうわけでもない。
ここでの `writeToLog` は rollback, recovery 用に出している log であるが、元に戻すときに
必要な情報は新たに書き込んだ情報でなく、もともと何が書いてあったか？ということなので
ここはログに古い値を残すことが正しい。そのためここを `oldval` を `newval` に置き換えると
意味がわからないことになってしまう。

Undo だけでなく Redo もするようになったら `newval` は必要だろうが、少なからず
SimpleDB は redo はないのでここの `newval` は必要なさそう。
演習問題を全く見てないのであれだが、もしかしたら演習問題として redo の実装があるのかもしれない。
それを考慮してなら `newval` があることも納得できる。

```java
public int setInt(Buffer buff, int offset, int newval) {
	int oldval = buff.contents().getInt(offset);
	BlockId blk = buff.block();
	return SetIntRecord.writeToLog(lm, txnum, blk, offset, oldval);
}

public int setString(Buffer buff, int offset, String newval) {
	String oldval = buff.contents().getString(offset);
	BlockId blk = buff.block();
	return SetStringRecord.writeToLog(lm, txnum, blk, offset, oldval);
}
```

# Chapter 6: Record Management

今まで作ってきたものを組み合わせてレコードを挿入するところなので特段難しいところはない。
しかし、元実装の変数名が何を表しているのかよくわからなかったので、Go 実装ではそのへんを
リネームして実装した。

また元実装だと record package を作っていたが、他の package から呼ばれることが多そうだったので
Go 実装では domain package に置くようにした。

# Chapter 7: Metadata Management

## 2022/4/24
`synchronized` がついている `getStatInfo` から同様に `synchronized` がついている
`refreshStatistics`/`calcTableStats` を呼んでいるが、lock はインスタンス単位で取られるからこれだと
`refreshStatistics`/`calcTableStats` がいつまでもつかえないのではないかと思った。
しかし、どうやら nest した synchronized は許されるらしい。

https://code-examples.net/en/q/464063


```
public synchronized StatInfo getStatInfo(String tblname,
                              Layout layout, Transaction tx)
private synchronized void refreshStatistics(Transaction tx)

private synchronized StatInfo calcTableStats(String tblname,
                              Layout layout, Transaction tx)
```

ただ `refreshStatistics`/`calcTableStats` は単体で呼び出されることはなく、常に `getStatInfo` またはコンストラクタからしか呼ばれないので
わざわざ `synchronized` をつける必要はなさそうである。
Go で実装するときは `refreshStatistics`/`calcTableStats` 内で `mutex.Lock` を取る必要はなさそう。

## 2022/4/26

`IndexInfo` が Chap.12 で出てくる `HashIndex` に依存しているという思わぬ実装になっていた。
Java で実装していたら実装を拝借してくるとかができたが、こちらは Go で実装しているのでそうもいかない。
さいわい、`HashIndex` の実装はシンプルだったのでほぼそのまま実装した。
コードを読むに index name が prefix につくファイルが最大で `NUM_BUCKET` 個できるようである。
index ファイルの suffix は `searchKey` の Hash 値を `NUM_BUCKET` で割った余り。
 `searchKey` というのが index を張ったカラムの値を表している模様。
同じ `searchKey` だと同じ index ファイルに記録されるから検索早くなるよねということらしい。
ただ、Hash 値から計算しているからわかる通り、`searchKey` が違うからといって違う index ファイルに記録されるわけではない。

# Chapter 8: Query Processing, Chapter 9: Parsing, Chapter 10: Planning

Chap 8 の `Term#reductionFactor(Plan p)` の Plan は Chap 10 で出てくるものなので
Chap 8 の時点では実装しないで Chap 8 ~ 10 の3章を読んでから実装するようにした。
Chap 10 まで読んでわかったが、`Plan` は `interface` だし、
Chapter 10 で作る `BasicQueryPlanner`, `BetterQueryPlanner` だと `reductionFactor` は
実装していなくても問題なかったので、3章一気読みしてから実装でなく Chapter 8 から順に実装していっても特段の問題はなかったようである。


# Chapter 12: Indexing
## 2022/7/10

IndexInfo の情報は `idxcat` テーブルに保存されている。
leaf node の layout は `IndexInfo.java#createIdxLayout` で定義されているやつが入ってくる。
```
---------------------------------------------------------------
| block number (int32) | slot id (int32) | dataval (Constant) |
---------------------------------------------------------------
```

## 2022/7/12
BTPage の flag は
- directory page のときは level
- leaf page のときは overflow block の block number を指している。oveflow block がないときは -1 が flag として入っている。`BTreeLeaf#insert` では overflow block の判定を flag が 0 以上で判定しているけど

(ref: text p.337)

## 2022/7/13

`BTreeLeaf#delete` は node split はするけど node merge はしていない。
この場合 overflow block は永遠に残り続けるので次の overflow block には record が1個もないのに `BTreeLeaf#tryOverflow` は true を返してしまいそうな気がする。

directory node に関しては dir entry を消すことをそもそもしていない。


バグっているか調査実験
- `SimpleDB_3.4/simpledb/metadata/IndexInfo.java` の HashIndex を BTreeIndex に変更
- `SimpleDB_3.4/simpledb/server/SimpleDB.java` の `QueryPlanner`, `UpdatePlanner` を `HeuristicQueryPlanner`, `IndexUpdatePlanner` に変更

```sql
jdbc:simpledb:hoge
create table hoge (name varchar(100));
create index idx_hoge on hoge (name);
insert into hoge (name) values ('hoge');
insert into hoge (name) values ('hoge');
insert into hoge (name) values ('hoge');
insert into hoge (name) values ('hoge');
insert into hoge (name) values ('hoge');
insert into hoge (name) values ('hoge');
select name from hoge;
select name from hoge where name = 'hoge';
delet from hoge where name = 'hoge';
select name from hoge;
select name from hoge where name = 'hoge';
```

```sql
~/simpledb_go ⑂main* ↑8 $ docker run --rm -it -v $(pwd)/SimpleDB_3.4:/app/SimpleDB_3.4 simpledb
Note: ./simpledb/jdbc/ResultSetAdapter.java uses or overrides a deprecated API.
Note: Recompile with -Xlint:deprecation for details.
creating new database
transaction 1 committed
database server ready
Connect>
jdbc:simpledb:hoge
creating new database
transaction 1 committed

SQL> create table hoge (name varchar(100));
transaction 2 committed
0 records processed

SQL> create index idx_hoge on hoge (name);
transaction 3 committed
0 records processed

SQL> insert into hoge (name) values ('hoge');
transaction 4 committed
1 records processed

SQL> insert into hoge (name) values ('hoge');
transaction 5 committed
1 records processed

SQL> insert into hoge (name) values ('hoge');
transaction 6 committed
1 records processed

SQL> insert into hoge (name) values ('hoge');
transaction 7 committed
1 records processed

SQL> insert into hoge (name) values ('hoge');
transaction 8 committed
1 records processed

SQL> insert into hoge (name) values ('hoge');
transaction 9 committed
1 records processed

SQL> select name from hoge;
                                                                                                 name
-----------------------------------------------------------------------------------------------------
                                                                                                 hoge
                                                                                                 hoge
                                                                                                 hoge
                                                                                                 hoge
                                                                                                 hoge
                                                                                                 hoge
transaction 10 committed

SQL> select name from hoge where name = 'hoge';
index on name used
                                                                                                 name
-----------------------------------------------------------------------------------------------------
                                                                                                 hoge
                                                                                                 hoge
                                                                                                 hoge
                                                                                                 hoge
                                                                                                 hoge
                                                                                                 hoge
transaction 11 committed

SQL> delete from hoge where name = 'hoge';
transaction 12 committed
6 records processed

SQL> select name from hoge;
                                                                                                 name
-----------------------------------------------------------------------------------------------------
transaction 13 committed

SQL> select name from hoge where name = 'hoge';
index on name used
                                                                                                 name
-----------------------------------------------------------------------------------------------------
                                                                                                 hoge
                                                                                                 hoge
transaction 14 committed
```

全レコードを消したあとでも `where` をつけて `SELECT` すると2件 record が返ってきた。
printf debug によると overflow block 由来だったのでやはり `tryOverflow` はバグっている。


# おまけ
## 2022/7/16
b-tree を実装するのが辛くなってきたので気分転換で `psql` から接続できるようにした。

```
cd exp
docker compose build
docker compose up

$ docker compose exec client bash
# psql -h db -Upostgres
postgres=# BEGIN;
BEGIN

$ docker compose exec client bash
# ngrep -x -q '.' 'host db'
```


```
postgres=# BEGIN;
BEGIN

postgres=*# select * from hoge;
 id
----
(0 rows)

postgres=*# ROLLBACK;
ROLLBACK
postgres=#
```

```
T 172.18.0.2:54100 -> 172.18.0.3:5432 [AP] #47
  51 00 00 00 0b 42 45 47    49 4e 3b 00                Q....BEGIN;.

T 172.18.0.3:5432 -> 172.18.0.2:54100 [AP] #49
  43 00 00 00 0a 42 45 47    49 4e 00 5a 00 00 00 05    C....BEGIN.Z....
  54                                                    T



T 172.18.0.2:54100 -> 172.18.0.3:5432 [AP] #51
  51 00 00 00 18 73 65 6c    65 63 74 20 2a 20 66 72    Q....select * fr
  6f 6d 20 68 6f 67 65 3b    00                         om hoge;.

T 172.18.0.3:5432 -> 172.18.0.2:54100 [AP] #53
  54 00 00 00 1b 00 01 69    64 00 00 00 40 00 00 01    T......id...@...
  00 00 00 17 00 04 ff ff    ff ff 00 00 43 00 00 00    ............C...
  0d 53 45 4c 45 43 54 20    30 00 5a 00 00 00 05 54    .SELECT 0.Z....T



T 172.18.0.2:54100 -> 172.18.0.3:5432 [AP] #55
  51 00 00 00 0e 52 4f 4c    4c 42 41 43 4b 3b 00       Q....ROLLBACK;.

T 172.18.0.3:5432 -> 172.18.0.2:54100 [AP] #57
  43 00 00 00 0d 52 4f 4c    4c 42 41 43 4b 00 5a 00    C....ROLLBACK.Z.
  00 00 05 49                                           ...I
```
