シークレットマネージャのOpenAPI定義は以下のページで公開されています。

https://manual.sakura.ad.jp/api/cloud/security-encryption/#tag/secretmanager-vault

secretmanager-api-goではここで公開されている定義からCloudHSM / KSM向けの定義を削除したのを利用しています。

現在APIの実装とOpenAPIの定義に食い違いがありライブラリ側で修正した定義を利用しています。
サービス側で修正があり次第本来の定義を利用するように変更します。