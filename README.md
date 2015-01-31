小田鳥子(おだとりこ)
==========

# これはなに？
音声自動応答システム「小田鳥子(おだとりこ)」です。

# できること
ファミレスの注文のような、以下のような会話を音声でおこなえます。

```
あなた：オーダーお願いします
小田　：いらっしゃいませ。ご注文をお伺いします
あなた：ハンバーグ
小田　：ハンバーグをお一つ
あなた：コーヒー
小田　：コーヒーをお一つ
あなた：やっぱり二つ
小田　：コーヒーをお二つ
あなた：やっぱりキャンセル
小田　：コーヒーをキャンセルしました
あなた：ウーロン茶
小田　：ウーロン茶をお一つ
あなた：以上でお願いします
小田　：ご注文を確認いたします。ハンバーグをお一つ、ウーロン茶をお一つ。以上でよろしいですか？
あなた：はい
小田　：ご注文ありがとうございました
```

家に居ながらファミレス気分を味わえます。


# 必要環境
## ハードウェア

* Linux が動作する PC
 * 想定環境は Raspberry Pi (Raspbian) ですが、開発は Mac で VirtualBox に Debian wheezy で行いました。
* マイク
* スピーカー

## ソフトウェア
2015/1/31 現在

* [Go 言語 1.4.1](https://golang.org/doc/install/source)
* [Open JTalk 1.08](http://open-jtalk.sp.nitech.ac.jp/)
 * 音声モデルとして [MMDAgent](http://www.mmdagent.jp/) の Mei を想定。違う音声モデルを使う場合はオプションで指定。
* [Julius 4.3.1](http://julius.sourceforge.jp/)
 * Julius 本体と「Julius記述文法音声認識実行キット」が必要です
* aplay
 * wav 再生用。Mac で実行するならオプションで afplay を指定でもよい

# インストール方法
```
go get github.com/nasu-tomoyuki/oda-toriko
cd $GOPATH/src/github.com/nasu-tomoyuki/oda-toriko
make
```

# 使い方
## Julius をモジュールモードで起動
```
julius -C $GOPATH/src/github.com/nasu-tomoyuki/oda-toriko/assets/jconf/oda-toriko.jconf -C ~/julius-kits/grammar-kit-v4.1/hmm_ptm.jconf -module
```

## 小田鳥子を起動
```
./build/oda-toriko
```

これで、localhost:10500 で動いている Julius に telnet で接続します。Julius がリモートで動いている場合は、オプションでアドレスを指定してください。
オプションは -h で確認してください。

Julius はフォアグラウンドで動作するので実行中はターミナルに帰ってきません。「小田鳥子」と同じ PC で動作させる場合は以下のように行末に & を付けてバックグラウンドで立ち上げるなり、ターミナルを複数開くなりしてください。

```
cd $GOPATH/src/github.com/nasu-tomoyuki/oda-toriko
julius -C ./assets/jconf/oda-toriko.jconf -C ~/julius-kits/grammar-kit-v4.1/hmm_ptm.jconf -module &

./build/oda-toriko
```


## 注文の流れ
1. 「注文いいですか？」「オーダーお願いします」などで店員を呼びます。
1. メニューの中からアイテムを選びます。
 * 「ハンバーグを１つ」のような感じで注文します。「３つ」のように数だけを言うと最後に注文した数量を変更します。「キャンセル」で最後の注文を取り消すことが出来ます。「ハンバーグをキャンセル」のようにアイテム名を言うことで、特定のアイテムを取り消すことが出来ます。
1. 「以上でお願いします」「それで」のような感じで注文を終えると、いままでの注文を読み上げます。よければ「はい」「それでお願いします」という感じで注文を確定します。

### メニュー
メニューには次のアイテムがあります。

* ハンバーグ
* 和風ハンバーグ
* チーズハンバーグ
* サイコロステーキ
* ヒレステーキ
* チキンステーキ
* とんかつ定食
* ショートケーキ
* チーズケーキ
* オレンジジュース
* コーラ
* ウーロン茶
* コーヒー
* アイスコーヒー
* ホットコーヒー

# 必要ソフトウェアのインストール例
インストールされていない場合や動作しない場合は以下を参考にしてください。
基本的には yum や apt-get では古いバージョンがインストールされるので、ソースからコンパイルします。

## Go 言語
git がインストールされていない場合は apt-get などであらかじめインストールします。
```
git clone https://go.googlesource.com/go go1.4
cd go1.4/src
./all.bash
mkdir -p ~/go/src
export GOPATH=$HOME/go
export GOROOT=$HOME/go1.4
export PATH=$PATH:$GOROOT/bin
```

## Julius
「Julius記述文法音声認識実行キット」もインストールします。

### インストール例

#### Julius 本体
```
wget --trust-server-names 'http://sourceforge.jp/frs/redir.php?m=iij&f=%2Fjulius%2F60273%2Fjulius-4.3.1.tar.gz'
tar zxvf julius-4.3.1.tar.gz
cd julius-4.3.1
./configure
make
sudo make install
```

#### グラマーキット
```
wget http://sourceforge.jp/projects/julius/downloads/51159/grammar-kit-v4.1.tar.gz
tar zxvf grammar-kit-v4.1
mkdir -p ~/julius-kits
mv grammar-kit-v4.1 ~/julius-kits
```

#### Raspberry Pi に USB マイクを付けている場合
```
export ALSADEV="plughw:1,0"
```
PulseAudio をインストールすればこれは必要ないけれど、自分の環境では PulseAudio をインストールしていると aplay がときどき止まるのでアンインストールしています。

## Open JTalk
Open JTalk は UTF-8 対応にコンパイルします。コンパイルには hts_engine が必要です。

### インストール例

#### hts_engine
```
wget http://downloads.sourceforge.net/hts-engine/hts_engine_API-1.09.tar.gz
tar zxvf hts_engine_API-1.09.tar.gz
cd hts_engine_API-1.09
./configure
make
sudo make install
```

#### Open JTalk 本体
```
wget http://downloads.sourceforge.net/open-jtalk/open_jtalk-1.08.tar.gz
tar zxvf open_jtalk-1.08.tar.gz
cd open_jtalk-1.08
./configure --with-charset=UTF-8 
make
sudo make install
```

#### 辞書
```
wget http://downloads.sourceforge.net/open-jtalk/open_jtalk_dic_utf_8-1.08.tar.gz
tar zxvf open_jtalk_dic_utf_8-1.08
sudo mkdir -p /var/lib/mecab/dic/open-jtalk/
sudo mv open_jtalk_dic_utf_8-1.08/ /var/lib/mecab/dic/open-jtalk/naist-jdic
```

#### 音声モデル
```
wget --trust-server-names http://sourceforge.net/projects/mmdagent/files/MMDAgent_Example/MMDAgent_Example-1.4/MMDAgent_Example-1.4.zip/download
sudo mkdir -p /usr/share/hts-voice/mei
sudo mv MMDAgent_Example-1.4/Voice/mei /usr/share/hts-voice/mei
```

# 技術解説
Julius による音声認識と、Open JTalk による音声合成が１０割です。プログラミングはほとんど必要とせずに会話を出来るのがすばらしいですね。

Julius には自由文を認識出来るディクテーションモードと、決められたルールで認識するグラマーモードとがあります。前者は Julius が認識した文章を自前で解析する必要があるのと、認識の処理負荷が高いので今回は利用していません。
文法を複数用意して、各状態で必要な文法のみ有効にしています。状態は HSM で管理していて、 https://github.com/hhkbp2/go-hsm の実装を利用しています。


