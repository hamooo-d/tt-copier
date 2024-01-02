#!/bin/bash

directories=(
	"/online/mxpprod/selectsystem_files/cardholder/in"
	"/online/mxpprod/selectsystem_files/transaction/out"
	"/online/mxpprod/selectsystem_files/merchant/in"
	"/online/mxpprod/selectsystem_files/evoucher/out"
	"/online/mxpprod/selectsystem_files/merchant/out"
)

bank_files=("CL." "KYCFile_" "reload_" "Rev_Reload_" "redemp_" "POS_RevAuthFile_" "APPLICATION." "KYC_ATM_" "SETT_TOPUP." "CORP_TOPUP." "EV_MERC")

tt_files=("PersoFile" "CL_TT" "FSO." "PAYOUT.")

declare -A bank_ids
bank_ids["000003"]="SB"
bank_ids["000002"]="ATIB"
bank_ids["000004"]="NAB"
bank_ids["000005"]="MED"
bank_ids["000001"]="TT"
bank_ids["000006"]="NCB"

current_date=$(date +%Y%m%d)

for dir in "${directories[@]}"; do
	mkdir -p "$dir"
	echo "Created directory: $dir"

	for id in "${!bank_ids[@]}"; do
		# Create dummy bank files
		for prefix in "${bank_files[@]}"; do
			filename="${dir}/${id}_${prefix}${current_date}.txt"
			touch "$filename"
			echo "Created dummy bank file: $filename"
		done

		# Create dummy TT files
		for prefix in "${tt_files[@]}"; do
			filename="${dir}/${id}_${prefix}${current_date}.txt"
			touch "$filename"
			echo "Created dummy TT file: $filename"
		done
	done
done

echo "Folder structure and dummy files created successfully."
