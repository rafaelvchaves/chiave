TYPE=$1
echo "latency,throughput" >> results_${TYPE}.csv
for ((d = 0; d < 3; d++)) do
	for ((k = 1; k <= 9; k++)) do
		NOPS=$((10**d * k * 1000))
		if [[ $NOPS -gt 50000 ]]; then
			break
		fi
		go run benchmark.go -mode a -wp 1 -nops $NOPS >> results_${TYPE}.csv
		sleep 2
	done
done
