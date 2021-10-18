module github.com/martenwallewein/quicsample

go 1.16

require (
	github.com/anacrolix/tagflag v1.3.0
	github.com/inconshreveable/log15 v0.0.0-20180818164646-67afb5ed74ec
	github.com/johannwagner/scion-optimized-connection v0.3.0
	github.com/lucas-clemente/quic-go v0.21.1
	github.com/netsec-ethz/scion-apps v0.3.0
	github.com/scionproto/scion v0.6.0
)

replace github.com/johannwagner/scion-optimized-connection => /home/marten/go/src/github.com/martenwallewein/scion-optimized-connection
