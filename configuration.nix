{ modulesPath, config, lib, pkgs, ... }: {
  imports = [
    (modulesPath + "/installer/scan/not-detected.nix")
    (modulesPath + "/profiles/qemu-guest.nix")
    ./disk-config.nix
  ];
  boot.loader.grub = {
    # no need to set devices, disko will add all devices that have a EF02 partition to the list already
    # devices = [ ];
    efiSupport = true;
    efiInstallAsRemovable = true;
  };
  services.openssh.enable = true;
  services.tailscale.enable = true;
  networking.hostName = "nanode";

  environment.systemPackages = map lib.lowPrio [
    pkgs.curl
    pkgs.gitMinimal
    pkgs.tailscale
  ];

  users.users.root.openssh.authorizedKeys.keys = [
    # change this to your ssh key
    "CHANGE"
  ];

  # Define a user account. Don't forget to set a password with ‘passwd’.
  users.users.azcoemu = {
    isNormalUser = true;
    description = "azurecosmos emu";
    extraGroups = [ "wheel" ];
    packages = with pkgs; [
    #  thunderbird
    ];
  };
  virtualisation.oci-containers = {
    backend = "docker";
    containers = {
      cosmos = {
        autoStart = true;
        image = "mcr.microsoft.com/cosmosdb/linux/azure-cosmos-emulator:mongodb";
        ports = [
          "8081:8081"
          "10250-10255:10250-10255"
        ];
        environment = {
          AZURE_COSMOS_EMULATOR_ENABLE_MONGODB_ENDPOINT = "4.0";
        };
      };
    };
  };

  system.stateVersion = "23.11";
}
