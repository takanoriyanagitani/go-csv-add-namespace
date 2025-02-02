#!/bin/sh

namespace=cafef00d-dead-beaf-face-864299792458
csvname=./sample.d/${namespace}.csv

gencsv(){
	echo typ,val > "${csvname}"

	echo apple,634 >> "${csvname}"
	echo berry,333 >> "${csvname}"
	echo grape,127 >> "${csvname}"
	echo honey,333 >> "${csvname}"
	echo lemon,599 >> "${csvname}"
	echo mango,255 >> "${csvname}"
	echo melon,100 >> "${csvname}"
	echo peach,128 >> "${csvname}"
}

#gencsv

export ENV_INPUT_CSVNAME="${csvname}"

./csv2named > ./sample.d/output.csv

duckdb -c "SELECT * FROM './sample.d/output.csv'"
