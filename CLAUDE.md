# CLAUDE.md

This file provides guidance to Claude Code (claude.ai/code) when working with code in this repository.

## プロジェクト概要

Go による HTML テンプレートレンダリング Web アプリ。

## ビルド・実行コマンド

```bash
go mod init github.com/ito/go-simple-webapp  # 初回のみ
go run ./...                                  # アプリ起動
go build -o bin/app ./...                     # バイナリビルド
go test ./...                                 # テスト実行
go test -v -run TestName ./...               # 特定テストのみ
```

## コードスタイル

- フォーマットは `gofmt`（編集後に自動適用）
- エラーハンドリングは明示的に行い、`_` での無視は原則禁止
- HTML テンプレートは `templates/` ディレクトリに配置
- ハンドラは `handlers/` パッケージに分離

## アーキテクチャ方針

- 標準ライブラリ優先（`net/http`）。外部フレームワークは必要になったときに検討
- テンプレートエンジンは `html/template`（XSS 対策のため `text/template` は使わない）

## 実装前のプランニング

大きな変更（新機能、リファクタリング、設計変更）を行う前に必ず実装計画を提示し、承認を得てから実装すること。
