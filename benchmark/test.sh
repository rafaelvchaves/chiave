FILE=results.csv
if [ ! -f "$FILE" ]; then
    echo "TP,L95" > results.csv
fi
for i in {0..50}
do
	go run tpvsl.go -nops $(( 1000*i ))
done
python3 graph.py