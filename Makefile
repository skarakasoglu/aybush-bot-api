build:
	echo "Building ${applicationName}..."
	go build -o ${executablePath}/${executableName} main.go

run:
	echo "Running ${applicationName} in ${ENVIRONMENT} mode..."
	${executablePath}/${executableName} --db-ip-address ${dbIpAddress} --db-port ${dbPort} --db-username ${dbUsername} --db-password ${dbPassword} --db-name ${dbName} \
		--cert-file ${certFile} --key-file ${keyFile}

all: build run