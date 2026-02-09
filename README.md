## Prerequisites

### Download files from Android Emulator
```
.\adb.exe -s emulator-5558 pull /sdcard/Android/data/com.xd.ro.roapk/files/Android/resources/pbbytes 'D:\Downloads\ROAPK\'
.\adb.exe -s emulator-5558 pull /sdcard/Android/data/com.xd.ro.roapk/files/Android/resources/script2 'D:\Downloads\ROAPK\'
```

(Optional) ### Decrypting unity3d and other files
If the unity3d files are encrypted following the steps below to reverse engineering the source code and find the secret to decrypt them, otherwise you can skip to step 3
1. Get the source code
```
find / -name "il2cpp.so"
# the global-metadata.dat file is also needed but it is not found in the latest version of the game to decompile the source code
find / -name "global-metadata.dat"
.\adb.exe -s emulator-5558 pull <path to il2cpp.so> 'D:\Downloads\ROAPK\'
.\adb.exe -s emulator-5558 pull <path to global-metadata.dat> 'D:\Downloads\ROAPK\'

# Note game guardian can be used to get the il2cpp.so and global-metadata.dat directly from App memory if you have trouble finding them with adb
```

2. Use Il2CppDumper to decompile the source code and find the secret to decrypt unity3d files
```
# Get Il2CppDumper from Github
https://github.com/Perfare/Il2CppDumper

# Decompile the source code
Il2CppDumper.exe <path to il2cpp.so> <path to global-metadata.dat>

# You will get a folder named Dummydll and file script.json, dump.cs and few other files in the same folder
```

3. Get the latest Ghidra from their official website and open the il2cpp.so file with it, then use the decompiled source code from Il2CppDumper to find the secret to decrypt unity3d files
one helpful link on how to do so is https://www.andnixsh.com/2023/05/how-to-use-il2cpph-scriptjson.html

3.1 Find the secret location by looking at the class and variable it is using to do that use tool like ILSpy or DnSpy 
to find the memory address offset which will be used in Ghidra to find the secret, for example in this case the secret is 
stored in a variable named "key" in a class named "EncryptHelper", then we can search for the class and variable in the decompiled source code from Il2CppDumper to find the memory address offset of the variable "key"
```
# Example code in C# from decompiled source code from ILSpy
[Token(Token = "0x600554D")]
[Address(RVA = "0x132A0B4", Offset = "0x132A0B4", VA = "0x132A0B4")]
public static byte[] DecryptBytes(AssetEncryptMode mode, byte[] datas)
```

3.2 Ghidra will help to get few useful information
1. Get the Field name of the secret variable which contains the secret field name for example
`Field$<>.4691E4C45AE35CECE39BC1214A59F5292AFD7` 
2. Use this key name to look up the memory offset location in dump.cs
```
# Example code in dump.cs
internal static readonly long 4691E4C45AE35CECE39BC1214A59F5292AFD79EB1C274F26E3F26793F2E33669 = 3328709681123367746; // 0x62A0
```
3. Take the offset value and open the Dummydll/Assembly-CSharp.dll in a Hex editor and search for the offset value
It should be in a format like in 
```
hex format `4f 66 66 73 65 74 <byte of offset value length e.g 06> <offset value> <byte of length of the data e.g 08> <the secret data in hex format for x bytes>`
plain text format `Offset.<offset value>.<none human readable string>`
```

### 1. Download AssetStudio from Github
```
https://github.com/Perfare/AssetStudio/releases/tag/v0.15.0
```

### 2. Use AssetStudio to extract all the .3d asset in pbtyes and scripts folder

### 3. Clone the following reop from Github
```
https://github.com/maodijim/pbtk
```

### 4. Convert all binary files to proto file
```
python <pbtk repo path>/extractors/from_binary.py <path to folder pbbytes asset extracted by AssetSutdio> proto
```

### 5. Install protobuf compiler
```
https://github.com/protocolbuffers/protobuf/releases/tag/v3.4.0

# or

go install google.golang.org/protobuf/cmd/protoc-gen-go@v1.25.0 
```

### 6. Generate proto file for go
```
cd <pbtk repo path>/extractors/proto
protoc -I . --go_out=proto_go .\*.proto
# Once Done copy all .pb.go file to ROMProject/Cmds
```

### 7. Clone this project
```
https://github.com/maodijim/ROMEncryption
```

### 8. Export Table_exchange, Table_Item, Table_Skill_ClsBranch_* from script2 folder using AssetSutdio


### 9. Use ROMEncryption to decrypt table_exchange.bytes into table_exchange.lua
```
#Visual studio is needed to compile the .cs file https://visualstudio.microsoft.com/
# Java https://www.java.com/en/download/

Build ROMEncryption\UtinyRipper first then
Build ROMEncryption
Create folder rawlua at the ROMEncryption.exe folder
Run ROMEncryption.exe
```

### 10. Run ./tools/luaObjectParser/[main.lua](tools%2FluaObjectParser%2Fmain.lua) against table_exchange.lua to get the json formatted exchangeItems.json
```
# Lua binaries
https://sourceforge.net/projects/luabinaries/files/5.4.2/Tools%20Executables
```

### 11. Change the variable in_files to table_item.bytes path and run ./tools/key_val_to_json.py to get json formatted items.json

### 12. Change the variable in_files to parent folder path of Table_Skill_ClsBranch_* with only Table_Skill_ClsBranch_* files inside the folder then run ./tools/key_val_to_json.py to get skills.json

### 13. Get access token from emulator
```
.\adb.exe -s emulator-5558 shell cat /data/data/com.xd.ro.roapk/shared_prefs/XDUserToken.xml
```

### 14. Build
```
$Env:GOOS = "linux"
$Env:GOARCH = "amd64"
go build -trimpath -ldflags "-w -s" .
```
