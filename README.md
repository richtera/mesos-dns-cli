# mesos-dns-cli
Sample client for mesos-dns written in go.

# Simple Usage
Usage of ./mesos-dns-cli:
  -dns="127.0.0.1:53": dns to query
  -domain="mesos": Domain
  -framework="marathon": Framework
  -hname="": hostname prefix if needed
  -index="": Index name to add to primary name (hack)
  -pname=".": port prefix if needed
  -protocol="tcp": Protocol tcp/udp and so on.
  -version=false: output the version

# What are Hostname Prefix and Port Prefix?
If you specify a --hname "-eh " and --pname " -ep " then the system will return your address as -eh +ip+ -ep +port+.

# Relative searches
As long as MARATHON_APP_ID is in the environment you can use relative service names to search for things.
Like

./mesos-dns-cli ../data/mysql
