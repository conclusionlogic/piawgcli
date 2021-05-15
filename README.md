# piawgcli

A tool to quickly and easily create WireGuard configuration files for PIA.
Use this tool on systems where there is no official PIA client support or
you need a command line based tool on, for example, a headless Linux server.

## Supported Systems

Builds are made available for the following platforms:

  * Linux/amd64
  * FreeBSD/amd64
  * Windows/amd64

If there are other platforms/architectures you would like to see builds for
then please open up a ticket.

## Usage

To generate a WireGuard configuration file suitable for processing by `wg-quick`:

```
piawgcli create-config --pia-id <id> --pia-password <pwd> --pia-region-id <regionid>
```

Use your PIA credentials and a valid PIA region id.  When sucessful, a valid
WireGuard config file will be output to stdout.  You may also use the `--outfile`
option to write the generated config to a file instead of stdout.

Add `--help` to the command line for a full listing of all available config options.

### Valid PIA Region IDs

So how do you find a valid PIA region id?  Use the `show-regions` command:

```
piawgcli show-regions [--search SEARCH] [--ping]
```

With no options, this command will display all available regions in alphabetical order.
Use the `--search` options to filter for specific region names or ids.  Use the
`--ping` option to ping each region and sort the results by ping time instead of
alphabetically.  The value in the `ID` column of the output of this command is the
region id you need to feed into the `create-config` command to generate your WireGuard
config file.

## Shortlived Sessions

Though the generated configs will work, they will not work forever.  If traffic stops
flowing on a configured PIA wg interface for a short time then PIA will drop your public
key from their server and you will then have to generate a new configuration, the
existing one will no longer work.  PIA has also documented that they routinely reboot
their servers for maintenance and other purposes and when that happens you would then
have to create a new configuration.

Because of this, a tool more tightly integrated with your router/gateway that can not
only generate valid configurations but also monitor and automatically regenerate 
connections as needed would be a more reliable solution.  I have other projects that
have tigher integrations with [VyOS]() and [pfSense]() that would be more useful for
those platforms.  However, this tool does provide a solid building block for those
platforms that do not have a more tightly integrated solution (vanilla Linux distros,
FreeBSD, etc.).  There is also a build of this for Windows, but is mostly just created
because it's "free" for me to buid it for Windows and my development machine is Windows
based.  Most Windows users of PIA should be using the official PIA client for Windows.