touch results.csv
echo "TP,L99\n" > results.csv
for i in {0..50}
do
	go run tpvsl.go -nops $(( 1000*i ))
done
