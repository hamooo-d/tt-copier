#!/usr/bin/bash

LOGFILE="/var/log/sftp_setup.log"

log() {
	echo "$(date '+%Y-%m-%d %H:%M:%S') - $1" | tee -a "$LOGFILE"
}

set -e
trap 'log "An error occurred. Exiting..."; exit 1' ERR

log "Starting SFTP user and directory setup."

sudo groupadd sftpusers
log "Group sftpusers created."

# Users to add
USERS=("atib" "nab" "sb" "tt" "med")

# Create parent directory for SFTP users and set proper permissions
sudo mkdir -p /home/sftp/files
sudo chown root:root /home/sftp/files
sudo chmod 755 /home/sftp/files
log "Parent directory /home/sftp/files created and permissions set."

# Add users and create their home directories with correct permissions
for user in "${USERS[@]}"; do
	user_dir="/home/sftp/files/${user^^}" # Convert username to uppercase
	sudo mkdir -p "$user_dir"
	sudo chown root:root "$user_dir"
	sudo chmod 755 "$user_dir"
	log "Home directory $user_dir created and permissions set."

	# Create user with specified home directory
	sudo useradd "$user" -g sftpusers -s /bin/false -d "$user_dir"
	log "User $user added to group sftpusers with no shell access."

	# Set password for user
	# It's better to prompt for a password or set it securely instead of hardcoding
	read -sp "Enter password for $user: " PASSWORD
	echo "$user:$PASSWORD" | sudo chpasswd
	unset PASSWORD
	log "Password set for $user."

	# Create subdirectories and set ownership to the user
	for subdir in "UAT" "Prod"; do
		for subsubdir in "from_tadawul" "to_tadawul"; do
			full_path="${user_dir}/${subdir}/${subsubdir}"
			sudo mkdir -p "$full_path"
			sudo chown "${user}:sftpusers" "$full_path"
			sudo chmod 700 "$full_path"
		done
		# Change ownership of UAT and Prod directories to the user
		sudo chown "${user}:sftpusers" "${user_dir}/${subdir}"
	done
done

# Create TTP directory with correct permissions
sudo mkdir -p /home/sftp/files/TTP
sudo chown root:root /home/sftp/files/TTP
sudo chmod 755 /home/sftp/files/TTP
log "TTP directory created and permissions set."

log "SFTP user and Directory setup completed."
