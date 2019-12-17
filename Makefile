default:
	pkgtemplates.sh
	GOOS=darwin go build -v -o jiffy

install:
	pkgtemplates.sh
	GOOS=darwin go build -v -o jiffy
	mv jiffy /usr/local/bin

docker:
	pkgtemplates.sh
	GOOS=linux go build -v -o jiffy
	# docker build -t pkger:jiffy .
	# docker run -p 3000:3000 pkger:jiffy

clean:
	rm /usr/local/bin/jiffy