#!/usr/bin/bash

LOGFILE="/var/log/sftp_undo.log"

log() {
	echo "$(date '+%Y-%m-%d %H:%M:%S') - $1" | tee -a "$LOGFILE"
}

set -e
trap 'log "An error occurred. Exiting..."; exit 1' ERR

log "Starting SFTP user and directory undo process."

USERS=("atib" "nab" "sb" "tt" "med" "ncb")

for user in "${USERS[@]}"; do
	sudo userdel "$user"
	log "User $user removed."
	if [ -d "/home/$user" ]; then
		sudo rm -rf "/home/$user"
		log "Home directory for $user removed."
	fi
done

sudo groupdel sftpusers
log "Group sftpusers removed."

DIRS=("/home/sftp/files/TTP" "/home/sftp/files")

for dir in "${DIRS[@]}"; do
	if [ -d "$dir" ]; then
		rm -rf "$dir"/{ATIB,SB,TT,MED,NAB}
		log "User directories removed under $dir."
	else
		log "Directory not found, so not removed: $dir"
	fi
done

log "SFTP user and directory undo process completed."
