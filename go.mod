module github.com/itoken417/go-simple-webapp

go 1.26.2

replace github.com/itoken417/goutils/logger => ../goutils/logger

replace github.com/itoken417/goutils/mailsender => ../goutils/mailsender

require (
	github.com/itoken417/goutils/logger v0.0.0-00010101000000-000000000000
	github.com/itoken417/goutils/mailsender v0.0.0-00010101000000-000000000000
	github.com/joho/godotenv v1.5.1
)

require github.com/mitchellh/panicwrap v1.0.0 // indirect
