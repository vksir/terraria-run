#!/usr/bin/env bash


mkdir -p ~/tml
cd ~/tml
curl -OL https://github.com/tModLoader/tModLoader/releases/download/v0.11.8.5/tModLoader.Linux.v0.11.8.5.tar.gz -x socks5://127.0.0.1:10808
tar -xzvf tModLoader.Linux.v0.11.8.5.tar.gz
rm tModLoader.Linux.v0.11.8.5.tar.gz


mkdir -p /home/steam/.local/share/Terraria/ModLoader/Mods
cd /home/steam/.local/share/Terraria/ModLoader/Mods
curl -OL https://mirror7.sgkoi.dev/tModLoader/download.php?Down=file/CalamityMod.tmod
curl -OL https://mirror7.sgkoi.dev/tModLoader/download.php?Down=file/CalamityModMusic.tmod
curl -OL https://mirror7.sgkoi.dev/tModLoader/download.php?Down=file/AutoReroll.tmod
curl -OL https://mirror7.sgkoi.dev/tModLoader/download.php?Down=file/AutoTrash.tmod
curl -OL https://mirror7.sgkoi.dev/tModLoader/download.php?Down=file/BossChecklist.tmod
curl -OL https://mirror7.sgkoi.dev/tModLoader/download.php?Down=file/LansUnlimitedBuffSlots.tmod
curl -OL https://mirror7.sgkoi.dev/tModLoader/download.php?Down=file/Luiafk.tmod
curl -OL https://mirror7.sgkoi.dev/tModLoader/download.php?Down=file/MaxStackPlus.tmod
curl -OL https://mirror7.sgkoi.dev/tModLoader/download.php?Down=file/RecipeBrowser.tmod
curl -OL https://mirror7.sgkoi.dev/tModLoader/download.php?Down=file/VeinMiner.tmod
curl -OL https://mirror7.sgkoi.dev/tModLoader/download.php?Down=file/OmniSwing.tmod
curl -OL https://mirror7.sgkoi.dev/tModLoader/download.php?Down=file/MagicStorageExtra.tmod
