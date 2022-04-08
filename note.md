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
https://github.com/goropikari/simpledb-go/blame/972526679d15cf5eb5d6d10a78b5192767714d38/backend/buffer/manager.go

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

https://github.com/goropikari/simpledb-go/blob/b1ba3f41fbc782214829c25b57e23e376f2cf052/backend/buffer/manager.go


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
