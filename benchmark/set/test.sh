TYPE=$1
NOPS=$2
echo "latency,throughput" >> results_${TYPE}.csv
go run benchmark.go -mode a -nsec 5 -nops $NOPS >> results_${TYPE}.csv
