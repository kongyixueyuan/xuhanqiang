export NODE_ID=5000
rm -rf main
echo "go build main.go:"
go build main.go
echo "---------------------------------------------------"
echo "./main addresslists:"
#./main addresslists
#address=$(./main addresslists)
address=`./main addresslists`

addr="addr"

for x in $address
do
	echo $x

	if [[ $x =~ '1' ]]
	then
		if [[ $addr =~ 'addr' ]]
		then
			addr=$x
		fi

	fi

done


echo "---------------------------------------------------"
echo "./main addblock -address \"13aEcxAFNxHiYWCLLXgCoAKFU1bkHwYe95\""
./main addblock -address "13aEcxAFNxHiYWCLLXgCoAKFU1bkHwYe95"
echo "---------------------------------------------------"


#echo "./main send -from '[\"13aEcxAFNxHiYWCLLXgCoAKFU1bkHwYe95\"]' -to '[\"18YeNiBDFFjSYubeJRpeCZzwyj9PKMwjaY\"]' -amount '[\"3\"]'"
#./main send -from '["13aEcxAFNxHiYWCLLXgCoAKFU1bkHwYe95"]' -to '["18YeNiBDFFjSYubeJRpeCZzwyj9PKMwjaY"]' -amount '["2"]'
#echo "---------------------------------------------------"


for x in $address
do 
	if [[ $x =~ '1' ]]
       	then
	echo "./main getbalance -address $x"
	./main getbalance -address $x
	echo "---------------------------------------------------"
fi
done


echo $addr

./main printchain

: '
address2=`echo $address`

IFS='  '
OLD_IFS="$IFS"
arr=($address2)
#echo ${arr[3]}

for s in ${arr[@]}
do
	echo "$s"
done

for i in 1 2 3
do 
echo "./main getbalance -address ${arr[$i]}"
./main getbalance -address ${arr[$i]}
echo "---------------------------------------------------"
done
'



