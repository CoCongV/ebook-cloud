outtest:
	go test -coverprofile=c.out -coverpkg ebook-cloud,ebook-cloud/api,ebook-cloud/api/apiv1,ebook-cloud/config,ebook-cloud/models,ebook-cloud/server,ebook-cloud/client,ebook-cloud/render

showcover:
	go tool cover -html=c.out
	
rundebug:
	air -c .air.conf