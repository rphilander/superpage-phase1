ratelimiter: cd ratelimiter && go run . --api 8080 --rate 1
parser: cd parser && go run . --api 8081 --ratelimiter 8080 --num-pages 4
snapshotdb: cd snapshotdb && go run . --api 8082 --db snapshots.db --parser 8081 --freq 1800
windowviewer: cd windowviewer && go run . --api 8083 --snapshotdb 8082
windowui: cd windowui && go run . --browser 3000 --windowviewer 8083
