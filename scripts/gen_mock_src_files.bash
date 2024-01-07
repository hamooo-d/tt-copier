#!/bin/bash

# NOTE: This script is used to set up the SFTP users and directories in development environment ONLY.

base_dir="/online/mxpprod/selectsystem_files"
declare -a sub_dirs=("cardholder/out" "transaction/out" "merchant/out" "evoucher/out")

declare -a prefixes=("CL." "KYCFile_" "reload_" "Rev_Reload_" "redemp_" "POS_RevAuthFile_" "APPLICATION." "KYC_ATM_" "SETT_TOPUP." "CORP_TOPUP." "EV_MERC" "PersoFile" "CL_TT" "FSO." "PAYOUT.")

declare -a bank_ids=("000001" "000002" "000003" "000004" "000005" "000006")

current_date=$(date +"%d%m%Y")

for dir in "${sub_dirs[@]}"; do
	mkdir -p "$base_dir/$dir"

	for prefix in "${prefixes[@]}"; do
		for id in "${bank_ids[@]}"; do
			touch "$base_dir/$dir/${prefix}${current_date}.${id}"
		done
	done
done
