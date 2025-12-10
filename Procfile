ratelimiter: cd ratelimiter && go run . --api 8080 --rate 1
parser: cd parser && go run . --api 8081 --ratelimiter 8080 --num-pages 4
querier: cd querier && go run . --api 8082 --parser 8081
webui: cd webui && go run . --browser 3000 --querier 8082
