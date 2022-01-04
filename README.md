## autoEncode

---

## 特徴

 [nekopanda/Amatsukaze](https://github.com/nekopanda/Amatsukaze)に対象ファイルを自動追加するための実行ファイル。

autoEncodeは所定のディレクトリを監視。TSファイルが存在した場合、AmatsukazeAddTask.exeに対象ファイルを渡してエンコード対象に自動追加するプログラムである。なお、重複追加しないような仕組みを実装している。



## 対応環境

Windows 10以降(64bit)



## 使い方

add.exe, start.exe, finish.exeを任意のフォルダに配置(Amatsukaze\batが良いかも)する。

下記のようなバッチファイルを作成、実行することで追加されたTSファイルのみをAddすることが可能になる。

```cmd
@echo off
add.exe -i <エンコード対象ディレクトリ> \ -o <エンコードファイル出力フォルダ> -e <AmatsukazeAddTask.exeの格納ディレクトリ> -p <プロファイル名>
```



また、実行後_サンプル.bat、実行前_サンプル.batとして以下のバッチファイルを作成、登録することでエンコード進捗状況も管理できる。

* 実行前_サンプル.bat

```cmd
@echo off

echo "エンコード終了 for" %IN_PATH%
"start.exe -f" %IN_PATH%
```



* 実行後_サンプル.bat

```cmd
@echo off

echo "エンコード終了 for" %IN_PATH%
"finish.exe -f" %IN_PATH%
```







## 著者

Yosuke Moriyama