outtest:
	go test -coverprofile=c.out -coverpkg ebook-cloud,ebook-cloud/api,ebook-cloud/api/apiv1,ebook-cloud/config,ebook-cloud/models,ebook-cloud/server,ebook-cloud/client,ebook-cloud/render,ebook-cloud/search

showcover:
	go tool cover -html=c.out

cover2html:
	go tool cover -html=c.out -o coverage.html 
	
rundebug:
	air -c .air.conf

build:
	sh build.sh & tar -cvf ebook-cloud.tar dist/* static/