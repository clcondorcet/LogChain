cmd=$1
docker run --rm -it -v $PWD:/data/ golang:1.18 sh -c "cd /data/ ; $cmd"