package config

type Mounts struct {
	// QFS mounts Quantos File System
	QFS string
	// QNS mounts Quantos Name server
	QNS string
	// QCNT mount Quantos content network
	QCNT string
	// QWAL mounts local wallets interface
	QWAL string
	// QLive is the live Quantos Network
	QLive string
	// QTest is the test net
	QTest string
	// QLocal is reserved for local / private tests
	QLocal string
}
