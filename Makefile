run_gsheet:
	go run ./gsheet_test/main.go

run_tg:
	go run ./telegram_test/main.go

run:
	go run ./main.go


GIT_SSH_COMMAND='ssh -i PATH/TO/KEY/FILE -o IdentitiesOnly=yes' git clone git@github.com:OWNER/REPOSITORY
