default:
	./pkgtemplates.sh
	GOOS=darwin go build -v -o jiffy

install-mac:
	./pkgtemplates.sh
	GOOS=darwin go build -v -o jiffy
	mv jiffy /usr/local/bin

install-linux:
	./pkgtemplates.sh
	GOOS=linux go build -v -o jiffy
	mv jiffy /usr/local/bin

docker:
	./pkgtemplates.sh
	GOOS=linux go build -v -o jiffy
	# docker build -t pkger:jiffy .
	# docker run -p 3000:3000 pkger:jiffy

all:
	./pkgtemplates.sh

	# macos
	GOOS=darwin GOARCH=amd64 go build -v -o jiffy 
	mv jiffy ./binary_distributions/darwin_amd64

	# linux
	GOOS=linux GOARCH=amd64 go build -v -o jiffy 
	mv jiffy ./binary_distributions/linux_amd64

	GOOS=linux GOARCH=arm go build -v -o jiffy 
	mv jiffy ./binary_distributions/linux_arm

	GOOS=linux GOARCH=arm64 go build -v -o jiffy 
	mv jiffy ./binary_distributions/linux_arm64

	GOOS=linux GOARCH=ppc64 go build -v -o jiffy 
	mv jiffy ./binary_distributions/linux_ppc64

	GOOS=linux GOARCH=mips64 go build -v -o jiffy 
	mv jiffy ./binary_distributions/linux_mips64
	
	GOOS=linux GOARCH=s390x go build -v -o jiffy 
	mv jiffy ./binary_distributions/linux_s390x

	# freebsd
	GOOS=freebsd GOARCH=amd64 go build -v -o jiffy 
	mv jiffy ./binary_distributions/freebsd_amd64
	
	GOOS=freebsd GOARCH=arm go build -v -o jiffy 
	mv jiffy ./binary_distributions/freebsd_arm

	# openbsd
	GOOS=openbsd GOARCH=amd64 go build -v -o jiffy 
	mv jiffy ./binary_distributions/openbsd_amd64
	
	GOOS=openbsd GOARCH=arm go build -v -o jiffy 
	mv jiffy ./binary_distributions/openbsd_arm




clean:
	rm /usr/local/bin/jiffy
