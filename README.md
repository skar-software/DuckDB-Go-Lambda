# Go-DuckDB-Lambda
this repository contains examples to successfully compile a golang binary to run on a AWS lambda

- main.go has a simple duckdb program where duckdb is run in memory mode and reads data from a csv file in a Cloudflare R2 Bucket (similar to S3, specific bucket name redacted)
- since duckdb uses C++, the go-duckdb lib requires CGO (via which we can call C Code from Golang)
- hence we need to have CSGO_ENABLED=true (Envirnoment Variable)
- but when we build from our host machine (in my case ubuntu 22.04) the binary will reference a perticular version of GLIBC (OS Level dep) where as when we try to run that binary in AWS Lambda where AWS Linux 2023 is used it has a different version of GLIBC hence we will get an error stating that this perticular version of GLIBC is not found
- hence we will have to compile the binary from that very OS Environment
- we build it using either a `EC2 Instance` or `Docker File`
## Build using EC2 Instance
- spin a ec2 instance with the same OS (AWS Linux 2023 in my case) as your lamba
- install golang and all required dependencies and compile the program and the resulting binary should work in AWS Lambda functions
## Build using Docker File
- this is the easy way, just utilize the Dockerfile given in this repository
- build the docker file
- `docker build -t my-golang-builder .`
- run the docker image built so we build the program
- `docker run --name my-golang-container my-golang-builder`
- copy the file from container to host fs 
- `docker cp my-golang-container:/built_file /home/user`

## HOME env not set
- in few environments such as AWS Lambda the $HOME env var is not set which is required in order to install extensions for duckdb in our example we are using httpfs to fetch from R2, you will get the below error
- `An error occurred while trying to automatically install the required extension 'httpfs':
	Can't find the home directory at ''
	Specify a home directory using the SET home_directory='/path/to/dir' option.`
- as this https://github.com/duckdb/duckdb/issues/3855 suggests, in AWS Lambda environment variables set variable `HOME` to `/tmp`