moderniejapanizer
=================

これは[modern.IE](https://www.modern.ie/ja-jp)のWindowsで日本語環境セットアップを自動で行うためのコマンドラインツールです。

手動で日本語環境のセットアップを行うには以下の手順が必要ですが、これを自動化するものです。

* [VirtualBox - modern.IEのWindows 7で日本語の表示と入力をできるようにする - Qiita](http://qiita.com/hnakamur/items/5f2f9e817dd0de60abb2)
* [Windows8.xのmodern.IEで日本語を入力、表示できるようにする。 - Qiita](http://qiita.com/hnakamur/items/cd37c9c8826afe4b4dda)

このツールは「管理者：コマンドプロンプト」上で実行してください。

## 前提条件

このツールはWindowsの表示言語が英語の状態で使うことを想定しています。表示言語が英語以外の場合は、まず英語に切り替えてから実行してください。

## 使い方

1. modern.IEの仮想マシンを起動してください。
2. インターネットエクスプローラのURL欄に https://github.com/hnakamur/moderniejapanizer/raw/master/dist/win_32bit/moderniejapanizer.exe と入力してダウンロードします。
    1. インターネットエクスプローラの画面下部に[Save]ボタンが表示されたら、その右の下三角ボタン▼を押して[Save as]メニューを選びます。
    2. [Save as]ダイアログが開いたら、左のツリーで[Faviroites]/[Downloads]を選んで[Save]ボタンを押します。
3. 以下の手順で「管理者：コマンドプロンプト」を開きます。
    * Windows VistaかWindows 7の場合:
        1. スタートメニューを開きます。
        2. "Command prompt"メニューで右クリックして"Run as administrator"メニューを選択します。
    * Windows 8か8.1の場合:
        1. Windowsキーを押しながらdを押してデスクトップウィンドウを開きます。
        1. Windowsキーを押しながらxを押し、その後aを押して"Command Prompt (Admin)"メニューを選択します。
4. 「管理者：コマンドプロンプト」開いたら以下のコマンドを実行し、仮想マシンが自動的に再起動されるまで待ちます。
    * `cd \Users\IEUser\Downloads`
    * `moderniejapanizer.exe`
5. modern.IEの日本語環境をエンジョイしてください！

## TODO

Windows 10 support.

## License

MIT
