migrateup:
	migrate -path ./db/migrations -database postgresql://postgres.jjhltxxtnfwiyypatwlw:e8Mtf9slTeXrk0gi@aws-1-ap-southeast-1.pooler.supabase.com:6543/postgres up

migratedown:
	migrate -path ./db/migrations -database postgresql://postgres.jjhltxxtnfwiyypatwlw:e8Mtf9slTeXrk0gi@aws-1-ap-southeast-1.pooler.supabase.com:6543/postgres down

migrateforce:
	migrate -path ./db/migrations -database postgresql://postgres.jjhltxxtnfwiyypatwlw:e8Mtf9slTeXrk0gi@aws-1-ap-southeast-1.pooler.supabase.com:6543/postgres force 1

sqlc:
	sqlc generate

server: 
	go run main.go

.PHONY: migrateup migratedown migratefoce slqc server