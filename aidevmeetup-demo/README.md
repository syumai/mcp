# aidevmeetup-demo

# 参考

以下のCloudflareのガイドをベースに作業します

-   https://developers.cloudflare.com/agents/guides/remote-mcp-server/

# 手順

## プロジェクトの初期化

-   認証不要のリモートMCPサーバーを作成

```
pnpm create cloudflare@latest \
calculator --template=cloudflare/ai/demos/remote-mcp-authless
```

## ローカルでのMCPサーバーの起動

```
cd calculator
pnpm start
```

## デバッグ

-   MCPのインスペクタを起動する
-   `http://localhost:8787/mcp` を登録する
    -   `/mcp` はStreamable HTTPのエンドポイント

```
npx @modelcontextprotocol/inspector@latest
```

## デプロイ

```
pnpm run deploy
```

## デプロイ後の動作確認

-   CloudflareのAI Playgroundを使う
    -   https://playground.ai.cloudflare.com/
