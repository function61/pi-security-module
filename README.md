[![Build Status](https://img.shields.io/travis/function61/pi-security-module.svg?style=for-the-badge)](https://travis-ci.org/function61/pi-security-module)
[![Download](https://img.shields.io/bintray/v/function61/pi-security-module/main.svg?style=for-the-badge&label=Download)](https://bintray.com/function61/pi-security-module/main/_latestVersion#files)

What is this?
-------------

Software for a separate trusted hardware device ("hardware security module") which
essentially acts just like [Keepass](http://keepass.info/) and only serves the
function of storing secrets.

If you use Keepass on your PC and your PC gets compromised by a virus or a hacker,
it's game over. But if you use a separate device for storing secrets, your PC compromise
does not expose your secrets. This software only exposes your secret when you physically
press a button on the device - and only exposes one secret per push acknowledge.


Links
-----

- [Architecture summary](https://docs.function61.com/pi-security-module/architecture/)
- [Comparison to alternatives](https://docs.function61.com/pi-security-module/comparison-to-alternatives/)
- [All documentation](https://docs.function61.com/pi-security-module) - everything you
  seek is probably here. The above links were just some of the most important links to
  this documentation site.


Features
--------

- No cloud
- Physical acknowledgement to expose a password by pressing a button on a U2F key (YubiKey for example)
- Supported secrets:
	* Passwords
	* OTP tokens (Google Authenticator)
	* SSH keys (via SSH agent protocol)
	* Keylists (["printed OTP list"](https://en.wikipedia.org/wiki/One-time_password#Hardcopy))
	* Freetext (any text content is treated as secret data)
- Create, view and list secrets in a folder hierarchy.
- Export database to Keepass format (for viewing in mobile devices when traveling etc.)
- Import data from Keepass format


Recommended hardware
--------------------

![](docs/pi-zero-in-wood-case.png)

I'm using [Raspberry Zero W](https://www.raspberrypi.org/products/pi-zero-w/)
with [wooden case](https://thepihut.com/products/zebra-zero-for-raspberry-pi-zero-wood)
and a [capacitive pushbutton](http://www.ebay.com/sch/?_nkw=ttp223).

It doesn't matter much which hardware you use, as long as you don't run anything else on
that system - to minimize the attack surface. For such a light use Raspberry Pi is
economical, although this project runs across processor architectures and operating systems
because Golang is so awesome. :)


Download & running
------------------

Click the "Download" badge at top of this readme and locate the binary for your OS/arch combo:

- For Raspberry Pi, download `pism_linux-arm`
- For Linux PC, download `pism_linux-amd64`

Note: don't worry about `public.tar.gz` - it's downloaded automatically if it doesn't exist.

Rename the downloaded binary to `pism`.

Pro-tip: you can download this directly to your Pi from command line:

```
$ mkdir pi-security-module/
$ cd pi-security-module
$ curl --fail --location -o pism <url to pism_linux-arm from Bintray>

# mark the binary as executable
$ chmod +x pism
```

Installation & running:

```
$ ./pism server install
Wrote unit file to /etc/systemd/system/pi-security-module.service
Run to enable on boot & to start now:
        $ systemctl enable pi-security-module
        $ systemctl start pi-security-module
        $ systemctl status pi-security-module
```

Looks good. You should now be able to access the web interface at `http://<ip of your pi>`.


How to build & develop
----------------------

[How to build & develop](https://github.com/function61/turbobob/blob/master/docs/external-how-to-build-and-dev.md)
(with Turbo Bob, our build tool). It's easy and simple!


TODO
----

- Tags to .JS command definitions
- Enter to confirm command dialog
- Data types for command fields (password)
