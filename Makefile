
run:
	go build -v -o mutualRead.exe ./cmd/
	./mutualRead.exe

store:
	go build -v --tags release -o ./mutualRead ./cmd/
	overseerctl restart mutualRead

kill:
	ps aux|grep mutualRead|grep -v grep |awk '{print $2}'|xargs kill -9
