# asdfによるgoのinstall

```
$ asdf plugin-add golang https://github.com/kennyp/asdf-golang.git
$ asdf install
```

# 起動

```
# mysql起動
$ docker-compose up -d
# schemaの適用（パスワードはroot）
$ mysql -uroot -h 127.0.0.1 -p < schema.sql
$ go run main.go
```

# apiのschema
```
# login
POST /api/login

# 請求書一覧
GET /api/invoices
{
  "start_date": "2024-01-01",
  "end_date": "2024-01-31"
}

# 請求書作成
GET /api/invoices
{
  "partner_id": 1001, # 取引先ID 現在静的に管理しているため1000 or 1001
  "issue_date": "2024-01-01", # 発行日
  "payment_amount": 10000, # 支払金額
  "payment_due_date": "2024-01-31" # 支払期日
}
```
