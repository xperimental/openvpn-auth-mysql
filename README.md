openvpn-auth-mysql
==================

A simple script written in Go to authenticate OpenVPN users against a simple MySQL table. The passwords are salted and hashed with SHA-256, a script to hash new passwords is also provided in `hash-password.go`.

[![Build Status](https://travis-ci.org/xperimental/openvpn-auth-mysql.png)](https://travis-ci.org/xperimental/openvpn-auth-mysql)

Usage
-----

To use this script in your OpenVPN installation add the following to your OpenVPN server configuration file:

    auth-user-pass-verify /path/to/script/openvpn-auth-mysql via-file

You can also use `via-env` instead of `via-file` if you want the credentials to be passed using environment variables instead of a temporary file.

MySQL table
-----------

The script expects a database with a table `openvpn_users` containing at least two columns `name` and `password`. The name field is directly matched against the username provided by the server. The password field's contents are formatted as follows:

    salt|algorithm|hash

- `salt` is the string used to salt the password hash
- `algorithm` is the used hashing algorithm (can only be `sha256` currently)
- `hash` is the hash of `salt + password` using the specified algorithm expressed as a lower-case hex-string

For example this would be the result of using the salt `f5cd8947` and the password `test`:

    f5cd8947|sha256|d9fb7d25153f6ec46c1e4cfd7f7eac02cbaccb0968692d4f3973eb9febad8402

Configuration
-------------

Before compiling the script change the configuration in `serverconfig.go` to match your MySQl server instance.
