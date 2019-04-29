# natureRemo-post-mackerel

- やること
  - NatureRemoで取得したデータ(気温・湿度・照度)をMackerelにPostさせる。

- 開発環境
  - go: 1.12.4
  - goland
  - Mackerel(Freeプラン)

- 実行環境
  - AWS Lambda(Go1.x)
    - Cloudwatch-Event(定期実行)

- デプロイ方法
```
$ make setup cross-build
$ zip -j deployment.zip ./build/pkg/main_linux_amd64/main
$ aws lambda update-function-code --function-name ${LAMBDA_FUNCTION_NAME} --zip-file fileb://deployment.zip --region ap-northeast-1
```

- Lambdaの設定
  - 事前にLambda関数を作成する
  - ランタイムはGo 1.x
  - メモリ: 128MB
  - タイムアウト: 5s
  - ネットワーク: 非VPC
  - ハンドラ: main
  - 環境変数
    - TZ: Asia/Tokyo
    - MKRKEY: mackerel_api_keyをセット
    - REMOTOKEN: natureRemoトークン
