default:
	GOOS=darwin go build -v -o jiffy

install-mac:
	GOOS=darwin go build -v -o jiffy
	mv jiffy /usr/local/bin

install-linux:
	GOOS=linux go build -v -o jiffy
	mv jiffy /usr/local/bin

docker:
	GOOS=linux go build -v -o jiffy
	# docker build -t jiffy .
	# docker run -p 3000:3000 jiffy

all:

	# macos
	GOOS=darwin GOARCH=amd64 go build -v -o jiffy 
	mv jiffy ./binary_distributions/darwin_amd64

	# linux
	GOOS=linux GOARCH=amd64 go build -v -o jiffy 
	mv jiffy ./binary_distributions/jiffy_linux_amd64

	GOOS=linux GOARCH=arm go build -v -o jiffy 
	mv jiffy ./binary_distributions/jiffy_linux_arm

	GOOS=linux GOARCH=arm64 go build -v -o jiffy 
	mv jiffy ./binary_distributions/jiffy_linux_arm64

	GOOS=linux GOARCH=ppc64 go build -v -o jiffy 
	mv jiffy ./binary_distributions/jiffy_linux_ppc64

	GOOS=linux GOARCH=mips64 go build -v -o jiffy 
	mv jiffy ./binary_distributions/jiffy_linux_mips64
	
	GOOS=linux GOARCH=s390x go build -v -o jiffy 
	mv jiffy ./binary_distributions/jiffy_linux_s390x

	# freebsd
	GOOS=freebsd GOARCH=amd64 go build -v -o jiffy 
	mv jiffy ./binary_distributions/jiffy_freebsd_amd64
	
	GOOS=freebsd GOARCH=arm go build -v -o jiffy 
	mv jiffy ./binary_distributions/jiffy_freebsd_arm

	# openbsd
	GOOS=openbsd GOARCH=amd64 go build -v -o jiffy 
	mv jiffy ./binary_distributions/jiffy_openbsd_amd64
	
	GOOS=openbsd GOARCH=arm go build -v -o jiffy 
	mv jiffy ./binary_distributions/jiffy_openbsd_arm




clean:
	rm /usr/local/bin/jiffy
