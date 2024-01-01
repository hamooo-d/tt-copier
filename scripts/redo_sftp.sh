#!/usr/bin/bash

# Define the log file
LOGFILE="/var/log/sftp_undo.log"

# Function to log messages
log() {
	echo "$(date '+%Y-%m-%d %H:%M:%S') - $1" | tee -a "$LOGFILE"
}

# Exit script on error
set -e
trap 'log "An error occurred. Exiting..."; exit 1' ERR

log "Starting SFTP user and directory undo process."

# Users to remove
USERS=("atib" "nab" "sb" "tt" "med")

# Remove users
for user in "${USERS[@]}"; do
	sudo userdel "$user"
	log "User $user removed."

	# Optionally, remove home directories as well
	if [ -d "/home/$user" ]; then
		sudo rm -rf "/home/$user"
		log "Home directory for $user removed."
	fi
done

# Remove group
sudo groupdel sftpusers
log "Group sftpusers removed."

# Directories to remove
DIRS=("/home/sftp/files/TTP" "/home/sftp/files")

# Remove directories
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

# End of script
