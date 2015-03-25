pdirector
=========

This simple TCP port proxy helps exposing services bound to localhost or other unaccessible network interfaces.

Usage
-----

	pdirector <local-port> <proxy-port> [<proxy-address>]
		local-port    - An opened port usually bound to the localhost
		proxy-port    - A proxy port for a localhost connection, which is remotely available
		proxy-address - A specific ip or named address where proxy-port should be opened.
	                	Default - 0.0.0.0
