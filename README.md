## Prerequisites

### Download files from Android Emulator
```
.\adb.exe -s emulator-5558 pull /sdcard/Android/data/com.xd.ro.roapk/files/Android/resources/pbbytes 'D:\Downloads\ROAPK\'
.\adb.exe -s emulator-5558 pull /sdcard/Android/data/com.xd.ro.roapk/files/Android/resources/script2 'D:\Downloads\ROAPK\'
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
https://github.com/protocolbuffers/protobuf/releases

# or

go install google.golang.org/protobuf/cmd/protoc-gen-go@latest
```

### 6. Generate proto file for go
```
protoc -I <path to the proto file folder> --go_out=proto_go .\*.proto
# Once Done copy all .pb.go file to ROMProject/Cmds
```

### 7. Clone this project
```
https://github.com/maodijim/ROMEncryption
```

### 8. Export Table_exchange, Table_Items, Table_Skill_ClsBranch_* from script2 folder using AssetSutdio


### 9. Use ROMEncryption to decrypt table_exchange.bytes into table_exchange.lua

### 10. Run ./tools/main.lua against table_exchange.lua to get the json formatted exchangeItems.json


### 11. Change the variable in_files to table_items.bytes path and run ./tools/key_val_to_json.py to get json formatted items.json

### 12. Change the variable in_files to parent folder path of Table_Skill_ClsBranch_* with only Table_Skill_ClsBranch_* files inside the folder then run ./tools/key_val_to_json.py to get skills.json

### 13. Get access token from emulator
```
.\adb.exe -s emulator-5558 shell cat /data/data/com.xd.ro.roapk/shared_prefs/XDUserToken.xml
```
