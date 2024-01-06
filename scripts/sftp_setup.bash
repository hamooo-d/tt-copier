#!/usr/bin/bash

# NOTE: This script is used to set up the SFTP users and directories in development environment ONLY.

LOGFILE="/var/log/sftp_setup.log"

log() {
	echo "$(date '+%Y-%m-%d %H:%M:%S') - $1" | tee -a "$LOGFILE"
}

set -e
trap 'log "An error occurred. Exiting..."; exit 1' ERR

log "Starting SFTP user and directory setup."

sudo groupadd sftpusers
log "Group sftpusers created."

USERS=("atib" "nab" "sb" "tt" "med")

sudo mkdir -p /home/sftp/files
sudo chown root:root /home/sftp/files
sudo chmod 755 /home/sftp/files
log "Parent directory /home/sftp/files created and permissions set."

for user in "${USERS[@]}"; do
	user_dir="/home/sftp/files/${user^^}"
	sudo mkdir -p "$user_dir"
	sudo chown root:root "$user_dir"
	sudo chmod 755 "$user_dir"
	log "Home directory $user_dir created and permissions set."

	sudo useradd "$user" -g sftpusers -s /bin/false -d "$user_dir"
	log "User $user added to group sftpusers with no shell access."

	read -sp "Enter password for $user: " PASSWORD
	echo "$user:$PASSWORD" | sudo chpasswd
	unset PASSWORD
	log "Password set for $user."

	for subdir in "UAT" "Prod"; do
		for subsubdir in "from_tadawul" "to_tadawul"; do
			full_path="${user_dir}/${subdir}/${subsubdir}"
			sudo mkdir -p "$full_path"
			sudo chown "${user}:sftpusers" "$full_path"
			sudo chmod 700 "$full_path"
		done

		sudo chown "${user}:sftpusers" "${user_dir}/${subdir}"
	done
done

sudo mkdir -p /home/sftp/files/TTP
sudo chown root:root /home/sftp/files/TTP
sudo chmod 755 /home/sftp/files/TT
log "TTP directory created and permissions set."

log "SFTP user and Directory setup completed."
