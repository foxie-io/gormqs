module example/pagination

go 1.25.2

replace github.com/foxie-io/gormqs => ./../..

require (
	gorm.io/driver/sqlite v1.6.0
	gorm.io/gorm v1.31.1
)

require github.com/mattn/go-sqlite3 v1.14.32 // indirect

require (
	github.com/foxie-io/gormqs v0.1.9
	github.com/jinzhu/inflection v1.0.0 // indirect
	github.com/jinzhu/now v1.1.5 // indirect
	golang.org/x/text v0.31.0 // indirect
)
