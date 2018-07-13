
str='XHQ_Serializdfadsfasfddsa'
echo $1
str=$1
echo $str
#if (("$str"==""))
#then
#	echo "no args" 
#fi


sed -i "s/$str/XHQ_$str/g" `grep $str -rl ./`
