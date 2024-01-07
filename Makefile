build:
	go build -o fccbot ./app/*.go

run:
	go run ./app/*.go

deploy:
	pm2 start "./fccbot -p" --name fccbot

redeploy:
	pm2 start "./fccbot -p -c" --name fccbot

stop:
	pm2 stop fccbot && pm2 delete fccbot
