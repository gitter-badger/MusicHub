auto:
	go install server/server/serveraux
	go install server/server/mp3metap
	go install server/server/returnMD5
	go install
runserver:
	server
