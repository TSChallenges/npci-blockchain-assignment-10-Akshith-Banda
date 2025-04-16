export PEER0_MEDLIFE_CA=/home/akshith/Practice/GoCode/bootcamp/npci-blockchain-assignment-10-Akshith-Banda/network/organizations/peerOrganizations/medlife.pharma.com/tlsca/tlsca.medlife.pharma.com-cert.pem

export CORE_PEER_TLS_ENABLED=true
export CORE_PEER_LOCALMSPID=MedlifeMSP
export CORE_PEER_TLS_ROOTCERT_FILE=$PEER0_MEDLIFE_CA
export CORE_PEER_MSPCONFIGPATH=/home/akshith/Practice/GoCode/bootcamp/npci-blockchain-assignment-10-Akshith-Banda/network/organizations/peerOrganizations/medlife.pharma.com/users/Admin@medlife.pharma.com/msp
export CORE_PEER_ADDRESS=localhost:8051