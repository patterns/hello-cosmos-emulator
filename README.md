# hello-cosmos-emulator
Learn Azure Cosmos DB by using the emulator (container)


### Requirements
- NixOS development environment
- remote target VM (Linode / DigitalOcean)

### Quickstart
Follow the steps from [nixos-anywhere](https://github.com/nix-community/nixos-anywhere).

Roughly,
1. vi configuration.nix (add configuration e.g. oci-containers)
2. nix flake lock
3. nix run github:nix-community/nixos-anywhere -- --flake /home/mydir/test#hetzner-cloud root@37.27.18.135

### Notes
With the Linode "Boot" configuration post-installer, the kernel dropdown needs to be "Direct Disk". Tried "Grub 2" by mistake.

Resorted to the Debian 11 disk image before reboot succeeded.

## Mahalo
[nixos-anywhere](https://github.com/nix-community/nixos-anywhere/blob/main/docs/quickstart.md)

[systemd container](https://nixos.wiki/wiki/Docker)

[Azure Cosmos Emulator](https://techcommunity.microsoft.com/t5/educator-developer-blog/local-development-using-azure-cosmos-db-emulator-at-no-cost/ba-p/4153822)

