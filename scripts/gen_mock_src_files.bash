#!/bin/bash

# Define the base directory and the subdirectories
base_dir="/online/mxpprod/selectsystem_files"
declare -a sub_dirs=("cardholder/out" "transaction/out" "merchant/out" "evoucher/out")

# Bank prefixes and IDs
declare -a bank_prefixes=("CL." "KYCFile_" "reload_" "Rev_Reload_" "redemp_" "POS_RevAuthFile_" "APPLICATION." "KYC_ATM_" "SETT_TOPUP." "CORP_TOPUP." "EV_MERC")
declare -a bank_ids=("000001" "000002" "000003" "000004" "000005" "000006")

# Current date in the format DDMMYYYY
current_date=$(date +"%d%m%Y")

# Create directories and files
for dir in "${sub_dirs[@]}"; do
	# Create directory
	mkdir -p "$base_dir/$dir"

	# Create dummy files
	for prefix in "${bank_prefixes[@]}"; do
		for id in "${bank_ids[@]}"; do
			touch "$base_dir/$dir/${prefix}${current_date}.${id}"
		done
	done
done
