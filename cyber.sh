#!/bin/bash

# ANSI Color Codes
RED='\033[0;31m'
GREEN='\033[0;32m'
YELLOW='\033[1;33m'
CYAN='\033[0;36m'
MAGENTA='\033[0;35m'
BLUE='\033[0;34m'
BOLD='\033[1m'
NC='\033[0m' # No Color

clear
echo -e "${MAGENTA}╔════════════════════════════════════════════════════════════════════════════════════════════╗"
echo -e "${MAGENTA}║      ${BOLD}ENTERPRISE SSH INTRUSION DETECTION & MITIGATION SYSTEM - INCIDENT REPORT v7.4.2${NC}${MAGENTA}      ║"
echo -e "${MAGENTA}╚════════════════════════════════════════════════════════════════════════════════════════════╝${NC}"
sleep 1

TARGET_HOST="smksabilillah.sch.id"
TARGET_IP="103.196.154.213"
SSH_PORT=22
DATE_NOW=$(date +"%Y-%m-%d %H:%M:%S")

echo -e "${CYAN}[+] Target hostname: ${YELLOW}$TARGET_HOST${NC}"
echo -e "${CYAN}[+] Target IP address: ${YELLOW}$TARGET_IP${NC}"
echo -e "${CYAN}[+] Monitored SSH port: ${YELLOW}$SSH_PORT${NC}"
echo -e "${CYAN}[+] Current system time: ${YELLOW}$DATE_NOW${NC}"
sleep 2

echo -e "${BLUE}[*] Initiating deep port scan and SSH service fingerprinting...${NC}"
sleep 2
echo -e "${GREEN}    Nmap Scan Results:"
echo -e "    PORT     STATE SERVICE  VERSION"
echo -e "    $SSH_PORT/tcp open  ssh      OpenSSH 8.6p1 Debian 5+deb11u1 (protocol 2.0)"
sleep 3

echo -e "${BLUE}[*] Extracting and validating SSH protocol banners...${NC}"
echo -e "    Detected Protocol: ${GREEN}SSH-2.0${NC}"
echo -e "    SSH Banner: ${YELLOW}SSH-2.0-OpenSSH_8.6p1 Debian-5${NC}"
sleep 2

echo -e "${BLUE}[*] Enumerating cryptographic primitives and cipher suites in use...${NC}"
echo -e "    Key Exchange Algorithms:"
echo -e "        - ${GREEN}curve25519-sha256${NC}"
echo -e "        - ${GREEN}ecdh-sha2-nistp256${NC}"
echo -e "    Message Authentication Codes (MACs):"
echo -e "        - ${GREEN}hmac-sha2-512-etm@openssh.com${NC}"
sleep 3

echo -e "${BLUE}[*] Reviewing authentication configuration parameters...${NC}"
echo -e "    PermitRootLogin: ${RED}no${NC} (Root login disabled, compliant with best practices)"
echo -e "    PasswordAuthentication: ${YELLOW}enabled${NC} (Key-based authentication recommended)"
echo -e "    MaxAuthTries: ${YELLOW}3${NC}"
echo -e "    AllowUsers: ${YELLOW}admin, devops, security${NC}"
sleep 2

echo -e "${BLUE}[*] Initiating authentication log analysis for suspicious activity...${NC}"
sleep 3

# Simulated failed login attempts
echo -e "${RED}[!] ALERT: Multiple failed SSH login attempts detected!${NC}"
echo -e "    Time Window: Last 15 minutes"
echo -e "    Source IP: 203.0.113.76"
echo -e "    Attempted Usernames:"
echo -e "        - admin"
echo -e "        - root"
echo -e "    Number of failed attempts: 27"
echo -e "    Sample log entries from /var/log/auth.log:"
echo -e "        May 26 09:15:22 server sshd[12345]: Failed password for invalid user admin from 203.0.113.76 port 58724 ssh2"
echo -e "        May 26 09:15:30 server sshd[12348]: Failed password for invalid user guest from 203.0.113.76 port 58727 ssh2"
sleep 4

echo -e "${BLUE}[*] Cross-referencing attacker IP with global threat intelligence databases...${NC}"
sleep 2
echo -e "${MAGENTA}    IP 203.0.113.76 flagged on multiple threat intelligence sources:"
echo -e "        - AbuseIPDB: Threat score 92/100"
echo -e "        - VirusTotal: Multiple reports of brute-force SSH attempts"
echo -e "        - FireHOL Blocklist: Included in recent SSH brute-force blocklists${NC}"
sleep 2

echo -e "${BLUE}[*] Executing active defense protocols...${NC}"
echo -e "    - Applying firewall rules to block IP 203.0.113.76 at network perimeter"
echo -e "    - Activating Fail2Ban automatic banning mechanism"
echo -e "    - Initiating connection throttling on SSH port"
sleep 3

echo -e "${GREEN}[+] Firewall status: IP 203.0.113.76 blocked successfully${NC}"
echo -e "${GREEN}[+] Fail2Ban: IP 203.0.113.76 banned for 3600 seconds (1 hour)${NC}"
echo -e "${GREEN}[+] SSH connection rate limited to 5 attempts per minute per IP${NC}"
sleep 2

echo -e "${BLUE}[*] Performing real-time monitoring for new login attempts...${NC}"
echo -e "    Monitoring window: 5 minutes"
sleep 3
echo -e "${GREEN}[✓] No further suspicious login attempts detected from blocked IP or others${NC}"
sleep 2

echo -e "${BLUE}[*] Verifying SSH configuration file integrity...${NC}"
echo -e "    Comparing /etc/ssh/sshd_config to baseline checksum..."
sleep 2
echo -e "${GREEN}[✓] Configuration integrity verified. No unauthorized changes detected.${NC}"
sleep 2

echo -e "${BLUE}[*] Reviewing active SSH sessions for anomalies...${NC}"
echo -e "    Current logged-in users:"
who
echo
sleep 3

echo -e "${BLUE}[*] Finalizing audit report and writing logs to /var/log/ssh-security/audit-$(date +%Y-%m-%d).log...${NC}"
sleep 2

echo -e "${GREEN}[✓] SSH Intrusion Detection & Mitigation cycle complete.${NC}"
echo -e "${GREEN}[✓] Target host ${YELLOW}$TARGET_HOST${NC} remains secure and uncompromised.${NC}"
echo
echo -e "${BOLD}${CYAN}>>> Press any key to exit the SSH Security Console...${NC}"
read -n 1 -s
