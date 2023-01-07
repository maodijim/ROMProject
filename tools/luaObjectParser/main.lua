
--require("mobdebug").start() -- same as start("localhost", 8172)
--
--require("oop")
--
--require("Import")
--
--require("FunctionLogin")
--
--require("NetProtocol")
--
--require("pbMgr")
--
--print("Start")
--local foo = 0
--for i = 1, 3 do
--    local function bar()
--        print("In bar")
--    end
--    foo = i
--    print("Loop")
--    bar()
--end
--print("End")

local json = require("dkjson")

require("Table_Buffer")
tableIn = Table_Buffer

-- require("Table_Item")
-- tableIn = Table_Item

-- require("Table_Monster")
-- tableIn = Table_Monster

-- require("Table_Exchange")
-- tableIn = Table_Exchange

----require("Table_ItemType")
----tableIn = Table_ItemType

tableIn = json.encode(tableIn, { indent = true })

io.open("tableOut.json", "w"):write(tableIn):close()

-- function tableToJson(inTable)
--     print("{")
--     local curIndex = 1
--     local count = 0
--     for _ in pairs(inTable) do count = count + 1 end
--     for index, data in pairs(inTable) do
--         local addComma = ""
--         print(string.format('"%s": {', index))
--
--         local subCount = 0
--         local subIndex = 1
--         for _ in pairs(data) do subCount = subCount + 1 end
--         for key, value in pairs(data) do
--             addSubComma = ""
--             if type(value) == 'table' then
--                 tableToJson(value)
--             else
--                 if subCount > subIndex then addSubComma = "," end
--                 print(string.format('"%s": "%s"%s', key, value, addSubComma))
--                 subIndex = subIndex + 1
--             end
--         end
--
--
--         if count > curIndex then addComma = "," end
--         print(string.format("}%s", addComma))
--         curIndex = curIndex + 1
--     end
--
--     print("}")
-- end

-- tableToJson(tableIn)


--local f = loadfile("Table_Boss.luac")
--print(f)

--
--function dump(o)
--    if type(o) == 'table' then
--        local s = '{ '
--        for k,v in pairs(o) do
--            if type(k) ~= 'number' then k = '"'..k..'"' end
--            s = s .. '['..k..'] = ' .. dump(v) .. ','
--        end
--        return s .. '} '
--    else
--        return tostring(o)
--    end
--end
--
--require("Table_Rune_21")
--require("Table_Rune_221")
--print(dump(Table_Rune_21))
--print(dump(Table_Rune_221))
