function encodeCSVRecord(record)
    if type(record) ~= "table" then
        return nil
    end

    -- RFC 4180
    local out = {}
    for i, s in ipairs(record) do
        local o = encodeCSVField(s)
        if o == nil then
            return nil
        end
        out[i] = o
    end
    return table.concat(out, ",") .. "\r\n"
end

function encodeCSVField(s)
    if type(s) ~= "string" then
        return nil
    end

    -- RFC 4180
    if string.match(s, "\r\n") ~= nil or string.match(s, '[",]') ~= nil then
        s = string.gsub(s, '"', '""')
        s = '"' .. s .. '"'
    end
    return s
end

function escapeQuery(s)
    if type(s) ~= "string" then
        return nil
    end

    local out = {}
    for i = 1, #s do
        local c = string.byte(s, i)
        local o
        if (
            c == 0x2D or                    -- -
            c == 0x2E or                    -- .
            (0x30 <= c and c <= 0x39) or    -- 0..9
            (0x41 <= c and c <= 0x5A) or    -- A..Z
            c == 0x5F or                    -- _
            (0x61 <= c and c <= 0x7A) or    -- a..z
            c == 0x7E                       -- ~
        ) then
            o = string.char(c)
        elseif c == 0x20 then   -- space
            o = "+"
        else
            o = string.format("%%%02X", c)
        end
        table.insert(out, o)
    end
    return table.concat(out)
end
